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
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetAssetListRequest struct {
	Limit         int       `form:"limit"`
	LastTime      time.Time `form:"last_time"`
	Tags          []string  `form:"tags"`
	MimeTypes     []string  `form:"mime_types"`
	ItemSources   []string  `form:"item_sources"`
	LinkSource    string    `form:"link_source"`
	LinkType      string    `form:"link_type"`
	ProfileSource string    `form:"profile_source"`
}

func GetAssetListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		_ = c.Error(err)

		return
	}

	request := GetAssetListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)

		return
	}

	profiles, err := database.QueryProfiles(database.DB, instance.Identity, 1, []int{})
	if err != nil {
		_ = c.Error(err)

		return
	}

	uris := make([]string, 0)
	uris = append(uris, rss3uri.New(instance).String())

	accounts := make([]model.Account, 0)

	for _, profile := range profiles {
		internalAccounts, err := database.QueryAccounts(database.DB, profile.ID, profile.Platform, 0)
		if err != nil {
			_ = c.Error(err)

			return
		}

		accounts = append(accounts, internalAccounts...)

		for _, account := range internalAccounts {
			uris = append(uris, strings.ToLower(
				rss3uri.New(
					rss3uri.NewAccountInstance(account.ID, constants.PlatformID(account.Platform).Symbol()),
				).String(),
			))
		}
	}
	if err := indexer.GetItems(accounts); err != nil {
		_ = c.Error(err)

		return
	}

	// Query assets form database
	assetModels, err := database.QueryAssets(database.DB, uris)
	if err != nil {
		_ = c.Error(err)

		return
	}

	uri := rss3uri.New(instance)

	var dateUpdated *time.Time
	assetList := make([]protocol.Item, len(assetModels))
	for i, assetModel := range assetModels {
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err = json.Unmarshal(assetModel.Attachments, &attachmentList); err != nil {
			_ = c.Error(err)

			return
		}

		if dateUpdated == nil {
			dateUpdated = &assetModel.DateUpdated
		}

		assetList[i] = protocol.Item{
			Identifier:  assetModel.Identifier,
			DateCreated: assetModel.DateCreated,
			DateUpdated: assetModel.DateUpdated,
			RelatedURLs: assetModel.RelatedURLs,
			Links:       fmt.Sprintf("%s/links", uri.String()),
			BackLinks:   fmt.Sprintf("%s/backlinks", uri.String()),
			Tags:        assetModel.Tags,
			Authors:     assetModel.Authors,
			Title:       assetModel.Title,
			Summary:     assetModel.Summary,
			Attachments: attachmentList,
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: dateUpdated,
		// TODO
		Identifier:     fmt.Sprintf("%s/assets", uri.String()),
		IdentifierNext: fmt.Sprintf("%s/assets", uri.String()),
		Total:          len(assetList),
		List:           assetList,
	})
}
