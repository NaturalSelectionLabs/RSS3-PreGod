package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
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

	fmt.Println(profiles)
	//if err := indexer.GetItems(profiles); err != nil {
	//	_ = c.Error(err)
	//
	//	return
	//}

	// Query assets form database
	assetModels, err := database.QueryAssets(database.DB, []string{rss3uri.New(instance).String()})
	if err != nil {
		_ = c.Error(err)

		return
	}

	uri := rss3uri.New(instance)

	var dateUpdated sql.NullTime
	assetList := make([]protocol.Item, 0, len(assetModels))
	for i, assetModel := range assetModels {
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err = json.Unmarshal(assetModel.Attachments, &attachmentList); err != nil {
			_ = c.Error(err)

			return
		}

		if !dateUpdated.Valid {
			dateUpdated.Valid = true
			dateUpdated.Time = assetModel.DateUpdated
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

	if len(assetList) == 0 {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: time.Now(),
		// TODO
		Identifier:     uri.String(),
		IdentifierNext: uri.String(),
		Total:          len(assetList),
		List:           assetList,
	})
}
