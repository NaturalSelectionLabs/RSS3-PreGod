package handler

import (
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"net/http"
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/gin-gonic/gin"
)

type GetNoteListRequest struct {
	Limit         int       `form:"limit"`
	LastTime      time.Time `form:"last_time"`
	Tags          []string  `form:"tags"`
	MimeTypes     []string  `form:"mime_types"`
	ItemSources   []string  `form:"item_sources"`
	LinkSource    string    `form:"link_source"`
	LinkType      string    `form:"link_type"`
	ProfileSource string    `form:"profile_source"`
}

func GetNoteListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		_ = c.Error(err)

		return
	}

	request := GetNoteListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)

		return
	}

	profiles, err := database.QueryProfiles(database.DB, instance.Identity, 1, []int{})
	if err != nil {
		_ = c.Error(err)

		return
	}

	uris := make([]string, 0)

	for _, profile := range profiles {
		accounts, err := database.QueryAccounts(database.DB, profile.ID, profile.Platform, 0)
		if err != nil {
			_ = c.Error(err)

			return
		}

		for _, account := range accounts {
			uris = append(uris, strings.ToLower(
				rss3uri.New(
					rss3uri.NewAccountInstance(account.ID, constants.PlatformID(account.Platform).Symbol()),
				).String(),
			))
		}
	}

	if err := indexer.GetItems(profiles); err != nil {
		_ = c.Error(err)

		return
	}

	// Query assets form database
	noteModels, err := database.QueryNotes(database.DB, uris)
	if err != nil {
		_ = c.Error(err)

		return
	}

	uri := rss3uri.New(instance)

	logger.Info(len(noteModels))
	noteList := make([]protocol.Item, len(noteModels))
	logger.Info(len(noteList))
	for i, noteModel := range noteModels {
		logger.Info(i)
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err := json.Unmarshal(noteModel.Attachments, &attachmentList); err != nil {
			_ = c.Error(err)

			return
		}

		noteList[i] = protocol.Item{
			Identifier:  noteModel.Identifier,
			DateCreated: noteModel.DateCreated,
			DateUpdated: noteModel.DateUpdated,
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
		DateUpdated: time.Now(),
		// TODO
		Identifier:     uri.String(),
		IdentifierNext: uri.String(),
		Total:          len(noteList),
		List:           noteList,
	})
}
