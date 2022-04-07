package handler

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
)

type GetLinkListRequest struct {
	Type          string   `form:"type"`
	Limit         int      `form:"limit"`
	LastInstance  string   `form:"last_instance"`
	Instance      string   `form:"instance"`
	LinkSources   []string `form:"link_sources"`
	ProfileSource string   `form:"profile_source"`
}

func GetLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetLinkListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	linkSources := make([]int, 0)
	for _, linkSource := range request.LinkSources {
		linkSources = append(linkSources, constants.LinkSourceName(linkSource).ID().Int())
	}

	var linkType *int

	if request.Type != "" {
		internalLinkType := constants.LinkTypeName(request.Type).ID().Int()
		linkType = &internalLinkType
	}

	linkModels, err := database.QueryLinks(
		database.DB,
		linkType,
		instance.Identity,
		linkSources,
		request.Limit,
	)
	if err != nil {
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	links := make([]protocol.Link, 0)

	for _, linkModel := range linkModels {
		links = append(links, protocol.Link{
			DateCreated: timex.Time(linkModel.CreatedAt),
			From:        linkModel.From,
			To:          linkModel.To,
			Type:        constants.LinkTypeID(linkModel.Type).String(),
			Source:      constants.ProfileSourceID(linkModel.Source).Name().String(),
			Metadata: protocol.LinkMetadata{
				Network: constants.NetworkSymbolCrossbell.String(),
				Proof:   "TODO",
			},
		})
	}

	var dateUpdated *timex.Time

	for _, link := range links {
		internalTime := link.DateCreated
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(link.DateCreated.Time()) {
			dateUpdated = &internalTime
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier:     fmt.Sprintf("%s/links", rss3uri.New(instance)),
		IdentifierNext: fmt.Sprintf("%s/links", rss3uri.New(instance)),
		DateUpdated:    dateUpdated,
		Total:          len(links),
		List:           links,
	})
}
