package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type GetAssetListRequest struct {
	Limit         int       `form:"limit"`
	LastTime      time.Time `json:"last_time"`
	Tags          []string  `json:"tags"`
	MimeTypes     []string  `json:"mime_types"`
	ItemSources   []string  `json:"item_sources"`
	LinkSource    string    `json:"link_source"`
	LinkType      string    `json:"link_type"`
	ProfileSource string    `json:"profile_source"`
}

func GetAssetListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetAssetListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	client := resty.New()
	// proof=kallydev&platform_id=6&network_id=12&item_type=note
	// TODO
	//profiles, err := database.QueryProfiles(database.DB, "", 1, []int{1})
	//if err != nil {
	//	_ = c.Error(api.ErrorNotFound)
	//
	//	return
	//}

	var indexerResponse indexer.Response

	_, err = client.NewRequest().
		EnableTrace().
		// TODO
		SetQueryParams(map[string]string{
			"proof":       instance.Identity,
			"platform_id": "1",
			"network_id":  "1",
			"item_type":   "asset",
		}).
		SetResult(&indexerResponse).
		Get(indexer.EndpointItem)
	if err != nil {
		_ = c.Error(err)

		return
	}

	if indexerResponse.Error.Code != 0 {
		_ = c.Error(err)

		return
	}

	assetList := make([]protocol.Item, 0)

	for _, asset := range indexerResponse.Data.Asset {
		attachments := make([]protocol.ItemAttachment, 0)
		for _, attachment := range asset.Attachments {
			attachments = append(attachments, protocol.ItemAttachment{
				Type:        "object",
				Address:     attachment.Address,
				MimeType:    attachment.MimeType,
				SizeInBytes: attachment.SizeInBytes,
			})
		}

		// TODO
		//dateCrated, err := time.Parse("", asset.DateCreated)
		//if err != nil {
		//	_ = c.Error(err)
		//
		//	return
		//}

		uri := rss3uri.New(instance)

		authors := make([]string, 0)
		for _, author := range asset.Authors {
			authors = append(authors, author)
		}

		assetInstance, err := rss3uri.NewInstance(
			string(constants.ItemTypeNote),
			asset.ItemId.Proof,
			constants.NetworkID(asset.ItemId.NetworkId).Symbol().String(),
		)
		if err != nil {
			_ = c.Error(err)

			return
		}

		assetList = append(assetList, protocol.Item{
			Identifier: rss3uri.New(assetInstance).String(),
			// TODO
			DateCreated: time.Now(),
			DateUpdated: time.Time{},
			RelatedURLs: nil,
			Links:       fmt.Sprintf("%s/links", uri),
			BackLinks:   fmt.Sprintf("%s/backlinks", uri),
			Tags:        asset.Tags,
			Authors:     authors,
			Title:       asset.Title,
			Summary:     asset.Summary,
			Attachments: attachments,
		})
	}

	if len(assetList) == 0 {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier: fmt.Sprintf("%s/assets", rss3uri.New(instance)),
		Total:      len(assetList),
		List:       assetList,
	})
}
