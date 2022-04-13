// nolint:dupl // TODO
package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

type GetLinkListRequest struct {
	Type           string     `form:"type"`
	Offset         int        `form:"offset"`
	Limit          int        `form:"limit"`
	LastTime       *time.Time `form:"last_time" time_format:"2006-01-02T15:04:05.000Z"`
	To             string     `form:"to"`
	LinkSources    []string   `form:"link_sources"`
	ProfileSources []string   `form:"profile_sources"`
}

func GetLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetInstance(c)
	if err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	request := GetLinkListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	linkModels, total, err := getLinkList(instance, request)
	if err != nil {
		api.SetError(c, api.ErrorIndexer, err)

		return
	}

	linkList := make([]protocol.Link, 0)

	for _, linkModel := range linkModels {
		linkList = append(linkList, protocol.Link{
			DateCreated: timex.Time(linkModel.CreatedAt),
			From:        rss3uri.New(rss3uri.NewAccountInstance(linkModel.From, constants.PlatformID(linkModel.FromPlatformID).Symbol())).String(),
			To:          rss3uri.New(rss3uri.NewAccountInstance(linkModel.To, constants.PlatformID(linkModel.ToPlatformID).Symbol())).String(),
			Type:        constants.LinkTypeID(linkModel.Type).Name().String(),
			Source:      constants.ProfileSourceID(linkModel.Source).Name().String(),
			Metadata: protocol.LinkMetadata{
				Network: constants.NetworkSymbolCrossbell.String(),
				Proof:   "TODO",
			},
		})
	}

	var dateUpdated *timex.Time

	for _, link := range linkList {
		internalTime := link.DateCreated
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(link.DateCreated.Time()) {
			dateUpdated = &internalTime
		}
	}

	uri := rss3uri.New(instance)

	var lastTime *time.Time

	for _, item := range linkList {
		assetDateCreated := item.DateCreated.Time()
		if lastTime == nil {
			lastTime = &assetDateCreated
		} else if lastTime.After(assetDateCreated) {
			lastTime = &assetDateCreated
		}
	}

	identifierNext := ""

	if len(linkList) == database.MaxLimit {
		nextQuery := c.Request.URL.Query()
		if lastTime != nil {
			nextQuery.Set("last_time", lastTime.Format(timex.ISO8601))
		}

		identifierNext = fmt.Sprintf("%s/links?%s", uri.String(), nextQuery.Encode())
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated:    dateUpdated,
		Identifier:     fmt.Sprintf("%s/links?%s", uri.String(), c.Request.URL.Query().Encode()),
		IdentifierNext: identifierNext,
		Total:          total,
		List:           linkList,
	})
}

func getLinkList(instance rss3uri.Instance, request GetLinkListRequest) ([]model.Link, int64, error) {
	internalDB := database.DB

	if request.Type != "" {
		internalDB = internalDB.Where("type = ?", constants.LinkTypeName(request.Type).ID().Int())
	}

	if request.LastTime != nil {
		internalDB = internalDB.Where("created_at <= ?", request.LastTime)
	}

	if request.To != "" {
		uri, err := rss3uri.Parse(strings.ToLower(request.To))
		if err != nil {
			return nil, 0, api.ErrorInvalidParams
		}

		internalDB = internalDB.Where(&model.Link{
			ToInstanceType: constants.StringToInstanceTypeID(uri.Instance.GetPrefix()).Int(),
			To:             strings.ToLower(uri.Instance.GetIdentity()),
			ToPlatformID:   constants.PlatformSymbol(uri.Instance.GetSuffix()).ID().Int(),
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
			FromPlatformID:   constants.StringToInstanceTypeID(instance.GetPrefix()).Int(),
			From:             strings.ToLower(instance.GetIdentity()),
			FromInstanceType: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
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
			FromPlatformID:   constants.StringToInstanceTypeID(instance.GetPrefix()).Int(),
			From:             strings.ToLower(instance.GetIdentity()),
			FromInstanceType: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
		}).
		Order("date_created DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	total := count - int64(request.Offset)

	if total < 0 {
		total = 0
	}

	return linkList, total, nil
}
