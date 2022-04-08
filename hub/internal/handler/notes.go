package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
)

type GetNoteListRequest struct {
	Limit          int      `form:"limit"`
	LastTime       string   `form:"last_time"`
	Tags           []string `form:"tags"`
	MimeTypes      []string `form:"mime_types"`
	ItemSources    []string `form:"item_sources"`
	LinkSources    []string `form:"link_source"`
	LinkType       string   `form:"link_type"`
	ProfileSources []string `form:"profile_source"`
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

	//var lastTime *time.Time
	//if request.LastTime != "" {
	//	internalLastTime, err := timex.Parse(request.LastTime)
	//	if err != nil {
	//		_ = c.Error(api.ErrorInvalidParams)
	//
	//		return
	//	}
	//
	//	t := internalLastTime.Time()
	//
	//	lastTime = &t
	//}

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
			Links:       fmt.Sprintf("%s/links", uri.String()),
			BackLinks:   fmt.Sprintf("%s/backlinks", uri.String()),
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
	// Get instance's all profiles
	profiles, err := database.QueryProfiles(database.DB, instance.GetIdentity(), 1, []int{})
	if err != nil {
		return nil, err
	}

	profileIDs := make([]string, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}

	//db := database.DB

	//if request.ProfileSource != "" {
	//	db = db.Where("sour")
	//}

	// Get instance's all accounts
	accounts := make([]model.Account, 0)
	if err := database.DB.
		Where("profile_id IN ?", profileIDs).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// Send get request to indexer
	if err := indexer.GetItems(accounts); err != nil {
		return nil, err
	}

	// Get instance's all notes
	notes := make([]model.Note, 0)
	if err := database.DB.
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Order("date_created DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func getNoteListsByLink(instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, error) {
	links := make([]model.Link, 0)
	if err := database.DB.
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

	// Send a request to indexer
	if err := indexer.GetItems(accounts); err != nil {
		return nil, err
	}

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

	notes := make([]model.Note, 0)
	if err := database.DB.
		Where("owner IN ?", owners).
		Order("date_created DESC").
		Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}
