// nolint:dupl // TODO
package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
)

type GetBackLinkListRequest struct {
	Type           string   `form:"type"`
	Offset         int      `form:"offset"`
	Limit          int      `form:"limit"`
	LastIdentifier string   `form:"last_identifier"`
	From           string   `form:"from"`
	LinkSources    []string `form:"link_sources"`
	ProfileSources []string `form:"profile_sources"`
}

func GetBackLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetInstance(c)
	if err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	request := GetBackLinkListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	backLinkModels, total, err := getBackLinkList(instance, request)
	if err != nil {
		api.SetError(c, api.ErrorIndexer, err)

		return
	}

	backLinkList := make([]protocol.Link, 0)

	for _, backLinkModel := range backLinkModels {
		backLinkList = append(backLinkList, protocol.Link{
			DateCreated: timex.Time(backLinkModel.CreatedAt),
			From:        rss3uri.New(rss3uri.NewAccountInstance(backLinkModel.From, constants.PlatformID(backLinkModel.FromPlatformID).Symbol())).String(),
			To:          rss3uri.New(rss3uri.NewAccountInstance(backLinkModel.To, constants.PlatformID(backLinkModel.ToPlatformID).Symbol())).String(),
			Type:        constants.LinkTypeID(backLinkModel.Type).Name().String(),
			Source:      constants.ProfileSourceID(backLinkModel.Source).Name().String(),
			Metadata: protocol.LinkMetadata{
				Network: constants.NetworkSymbolCrossbell.String(),
				Proof:   "TODO",
			},
		})
	}

	// Get date updated
	var dateUpdated *timex.Time

	for _, backLink := range backLinkList {
		internalTime := backLink.DateCreated
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(backLink.DateCreated.Time()) {
			dateUpdated = &internalTime
		}
	}

	uri := rss3uri.New(instance)

	var lastItem *protocol.Link

	for _, item := range backLinkList {
		internalItem := item

		if lastItem == nil || lastItem.DateCreated.Time().After(internalItem.DateCreated.Time()) {
			lastItem = &internalItem
		}
	}

	identifierNext := ""

	if len(backLinkList) == database.MaxLimit {
		nextQuery := c.Request.URL.Query()
		if lastItem != nil {
			nextQuery.Set("last_identifier", lastItem.From)
		}

		identifierNext = fmt.Sprintf("%s/backlinks?%s", uri.String(), nextQuery.Encode())
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated:    dateUpdated,
		Identifier:     fmt.Sprintf("%s/backlinks?%s", uri.String(), c.Request.URL.Query().Encode()),
		IdentifierNext: identifierNext,
		Total:          total,
		List:           backLinkList,
	})
}

func getBackLinkList(instance rss3uri.Instance, request GetBackLinkListRequest) ([]model.Link, int64, error) {
	internalDB := database.DB

	if request.Type != "" {
		internalDB = internalDB.Where("type = ?", constants.LinkTypeName(request.Type).ID().Int())
	}

	if request.LastIdentifier != "" {
		lastIdentifier, err := rss3uri.Parse(strings.ToLower(request.LastIdentifier))
		if err != nil {
			return nil, 0, err
		}

		var lastItem model.Link
		if err := database.DB.Where(&model.Link{
			From: lastIdentifier.Instance.GetIdentity(),
		}).First(&lastItem).Error; err != nil {
			return nil, 0, err
		}

		internalDB = internalDB.
			Where("created_at <= ?", lastItem.CreatedAt).
			Where("from != ?", lastItem.From)
	}

	if request.From != "" {
		uri, err := rss3uri.Parse(strings.ToLower(request.From))
		if err != nil {
			return nil, 0, api.ErrorInvalidParams
		}

		internalDB = internalDB.Where(&model.Link{
			FromInstanceType: constants.StringToInstanceTypeID(uri.Instance.GetPrefix()).Int(),
			From:             strings.ToLower(uri.Instance.GetIdentity()),
			FromPlatformID:   constants.PlatformSymbol(uri.Instance.GetSuffix()).ID().Int(),
		})
	}

	var linkSources []int

	if request.LinkSources != nil && len(request.LinkSources) > 0 {
		for _, source := range request.LinkSources {
			linkSources = append(linkSources, constants.LinkSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", linkSources)
	}

	linkList := make([]model.Link, 0)
	if err := internalDB.
		Where(&model.Link{
			ToInstanceType: constants.StringToInstanceTypeID(instance.GetPrefix()).Int(),
			To:             strings.ToLower(instance.GetIdentity()),
			ToPlatformID:   constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
		}).
		Offset(request.Offset).
		Limit(request.Limit).
		Order("created_at DESC").
		Find(&linkList).Error; err != nil {
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Link{}).
		Where(&model.Link{
			ToInstanceType: constants.StringToInstanceTypeID(instance.GetPrefix()).Int(),
			To:             strings.ToLower(instance.GetIdentity()),
			ToPlatformID:   constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
		}).
		Order("created_at DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return linkList, count, nil
}
