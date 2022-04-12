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

// nolint:dupl,funlen // TODO
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

	var noteModels []model.Note

	var total int64

	if len(request.LinkSources) != 0 || request.LinkType != "" {
		noteModels, total, err = getNoteListsByLink(instance, request)
	} else {
		noteModels, total, err = getNoteListByInstance(instance, request)
	}

	if err != nil {
		_ = c.Error(err)

		return
	}

	uri := rss3uri.New(instance)

	var dateUpdated *timex.Time

	noteList := make([]protocol.Item, len(noteModels))

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
			Links:       fmt.Sprintf("%s/links", noteModel.Identifier),
			BackLinks:   fmt.Sprintf("%s/backlinks", noteModel.Identifier),
			Tags:        noteModel.Tags,
			Authors:     noteModel.Authors,
			Title:       noteModel.Title,
			Summary:     noteModel.Summary,
			Attachments: attachmentList,
		}
	}

	var lastTime *time.Time

	for _, item := range noteList {
		assetDateCreated := item.DateCreated.Time()
		if lastTime == nil {
			lastTime = &assetDateCreated
		} else if lastTime.After(assetDateCreated) {
			lastTime = &assetDateCreated
		}
	}

	identifierNext := ""

	if len(noteList) == database.MaxLimit {
		nextQuery := c.Request.URL.Query()
		if lastTime != nil {
			nextQuery.Set("last_time", lastTime.Format(timex.ISO8601))
		}

		identifierNext = fmt.Sprintf("%s/notes?%s", uri.String(), nextQuery.Encode())
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated:    dateUpdated,
		Identifier:     fmt.Sprintf("%s/notes?%s", uri.String(), c.Request.URL.Query().Encode()),
		IdentifierNext: identifierNext,
		Total:          total,
		List:           noteList,
	})
}

// nolint:funlen // TODO
func getNoteListByInstance(instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, int64, error) {
	// Get instance's profiles
	var profiles []model.Profile

	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) > 0 {
		var profileSources []int
		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	if err := internalDB.Where(&model.Profile{
		ID:       strings.ToLower(instance.GetIdentity()),
		Platform: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
	}).Find(&profiles).Error; err != nil {
		return nil, 0, err
	}

	profileIDs := make([]string, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}

	// Get accounts
	internalDB = database.DB

	// Get instance's all accounts
	accounts := make([]model.Account, 0)
	if err := internalDB.
		Where("profile_id IN ?", profileIDs).
		Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	// TODO Refine it
	// Send get request to indexer
	go func() {
		if err := indexer.GetItems(instance, accounts); err != nil {
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

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		authors := []string{
			rss3uri.New(instance).String(),
		}

		for _, account := range accounts {
			accountInstance := rss3uri.NewAccountInstance(account.Identity, constants.PlatformID(account.Platform).Symbol())
			authors = append(authors, rss3uri.New(accountInstance).String())
		}

		internalDB = internalDB.Where("authors && ?", pq.StringArray(authors))
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
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Note{}).
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Order("date_created DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return notes, count, nil
}

// nolint:funlen,gocognit // TODO
func getNoteListsByLink(instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, int64, error) {
	links := make([]model.Link, 0)

	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) > 0 {
		var profileSources []int
		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

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
		return nil, 0, err
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
		return nil, 0, err
	}

	// TODO Refine it
	// Send a request to indexer
	go func() {
		if err := indexer.GetItems(instance, accounts); err != nil {
			logger.Error(err)
		}
	}()

	owners := make([]string, len(links))

	for _, link := range links {
		ownerInstance, err := rss3uri.NewInstance(
			constants.InstanceTypeID(link.ToInstanceType).String(),
			link.To,
			constants.PlatformID(link.ToPlatformID).Symbol().String(),
		)
		if err != nil {
			return nil, 0, err
		}

		owners = append(owners, strings.ToLower(rss3uri.New(ownerInstance).String()))
	}

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

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		authors := []string{
			rss3uri.New(instance).String(),
		}

		for _, account := range accounts {
			accountInstance := rss3uri.NewAccountInstance(account.Identity, constants.PlatformID(account.Platform).Symbol())
			authors = append(authors, rss3uri.New(accountInstance).String())
		}

		internalDB = internalDB.Where("authors && ?", pq.StringArray(authors))
	}

	notes := make([]model.Note, 0)
	if err := internalDB.
		Where("owner IN ?", owners).
		Limit(request.Limit).
		Order("date_created DESC").
		Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Note{}).
		Where("owner IN ?", owners).
		Order("date_created DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return notes, count, nil
}
