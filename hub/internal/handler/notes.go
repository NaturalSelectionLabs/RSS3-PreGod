package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type GetNoteListRequest struct {
	Limit          int        `form:"limit"`
	LastTime       *time.Time `form:"last_time" time_format:"2006-01-02T15:04:05.000Z"`
	Tags           []string   `form:"tags"`
	MimeTypes      []string   `form:"mime_types"`
	ItemSources    []string   `form:"item_sources"`
	LinkSources    []string   `form:"link_sources"`
	LinkType       string     `form:"link_type"`
	ProfileSources []string   `form:"profile_sources"`
}

// nolint:funlen // TODO
func GetNoteListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		_ = c.Error(err)

		return
	}

	request := GetNoteListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)

		return
	}

	// TODO Parse last time

	var noteModels []model.Note
	if len(request.LinkSources) != 0 || request.LinkType != "" {
		noteModels, err = getNoteListsByLink(instance, request)
	} else {
		noteModels, err = getNoteListByInstance(instance, request)
	}

	if err != nil {
		_ = c.Error(err)

		return
	}

	uri := rss3uri.New(instance)

	var dateUpdated *timex.Time

	noteList := make([]protocol.Item, len(noteModels))

	// nolint:dupl // TODO
	for i, noteModel := range noteModels {
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err = json.Unmarshal(noteModel.Attachments, &attachmentList); err != nil {
			_ = c.Error(err)

			return
		}

		internalTime := timex.Time(noteModel.DateUpdated)
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(noteModel.DateUpdated) {
			dateUpdated = &internalTime
		}

		noteList[i] = protocol.Item{
			Identifier:  noteModel.Identifier,
			DateCreated: timex.Time(noteModel.DateCreated),
			DateUpdated: timex.Time(noteModel.DateUpdated),
			RelatedURLs: noteModel.RelatedURLs,
			Links:       fmt.Sprintf("%s/links", noteModel.Owner),
			BackLinks:   fmt.Sprintf("%s/backlinks", noteModel.Owner),
			Tags:        noteModel.Tags,
			Authors:     noteModel.Authors,
			Title:       noteModel.Title,
			Summary:     noteModel.Summary,
			Attachments: attachmentList,
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: dateUpdated,
		// TODO
		Identifier:     fmt.Sprintf("%s/notes", uri.String()),
		IdentifierNext: fmt.Sprintf("%s/notes", uri.String()),
		Total:          len(noteList),
		List:           noteList,
	})
}

func getNoteListByInstance(instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, error) {
	// Get instance's profiles
	profiles, err := database.QueryProfiles(database.DB, instance.GetIdentity(), 1, []int{})
	if err != nil {
		return nil, err
	}

	profileIDs := make([]string, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}

	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		// Convert profile name to id
		var profileSources []int
		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	// Get instance's all accounts
	accounts := make([]model.Account, 0)
	if err := internalDB.
		Where("profile_id IN ?", profileIDs).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// TODO Refine it
	// Send get request to indexer
	go func() {
		if err := indexer.GetItems(accounts); err != nil {
			logger.Error(err)
		}
	}()

	// Get instance's notes
	internalDB = database.DB

	if request.LastTime != nil {
		internalDB = internalDB.Where("date_created <= ?", request.LastTime)
	}

	if request.Tags != nil && len(request.Tags) != 0 {
		internalDB = internalDB.Where("tags && ?", pq.StringArray(request.Tags))
	}

	if request.ItemSources != nil && len(request.ItemSources) != 0 {
		internalDB = internalDB.Where("source IN ?", request.ItemSources)
	}

	notes := make([]model.Note, 0)
	if err := internalDB.
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Limit(request.Limit).
		Order("date_created DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func getNoteListsByLink(instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, error) {
	links := make([]model.Link, 0)

	internalDB := database.DB

	if request.LinkType != "" {
		internalDB = internalDB.Where("type = ?", constants.LinkTypeName(request.LinkType).ID().Int())
	}

	if request.LinkSources != nil && len(request.LinkSources) != 0 {
		var sources []int
		for _, linkSource := range request.LinkSources {
			sources = append(sources, constants.LinkSourceName(linkSource).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", sources)
	}

	if err := internalDB.
		Where(&model.Link{
			From: instance.GetIdentity(),
		}).
		Find(&links).Error; err != nil {
		return nil, err
	}

	targets := make([]string, 0)

	for _, link := range links {
		targets = append(targets, link.To)
	}

	accounts := make([]model.Account, 0)
	if err := database.DB.
		Where("profile_id IN ?", targets).
		// TODO profile_platform
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// TODO Refine it
	// Send a request to indexer
	go func() {
		if err := indexer.GetItems(accounts); err != nil {
			logger.Error(err)
		}
	}()

	owners := make([]string, len(links))
	for _, link := range links {
		instance, err := rss3uri.NewInstance(
			constants.InstanceTypeID(link.ToInstanceType).String(),
			link.To,
			constants.PlatformID(link.ToPlatformID).Symbol().String(),
		)
		if err != nil {
			return nil, err
		}

		owners = append(owners, strings.ToLower(rss3uri.New(instance).String()))
	}

	internalDB = database.DB

	if request.LastTime != nil {
		internalDB = internalDB.Where("date_created >= ?", request.LastTime)
	}

	if request.Tags != nil && len(request.Tags) != 0 {
		internalDB = internalDB.Where("tags && ?", pq.StringArray(request.Tags))
	}

	if request.ItemSources != nil && len(request.ItemSources) != 0 {
		internalDB = internalDB.Where("source IN ?", request.ItemSources)
	}

	notes := make([]model.Note, 0)
	if err := internalDB.
		Where("owner IN ?", owners).
		Limit(request.Limit).
		Order("date_created DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}
