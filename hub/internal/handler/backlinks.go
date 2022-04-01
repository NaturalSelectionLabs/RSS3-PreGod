package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetBackLinkListRequest struct {
	Type          string   `form:"type"`
	Limit         int      `form:"limit"`
	LastInstance  string   `form:"last_instance"`
	Instance      string   `form:"instance"`
	LinkSources   []string `form:"link_sources"`
	ProfileSource string   `form:"profile_source"`
}

func GetBackLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	request := GetBackLinkListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	linkSources := make([]int, 0)
	for _, linkSource := range request.LinkSources {
		linkSources = append(linkSources, constants.LinkSourceName(linkSource).ID().Int())
	}

	linkModels, err := database.QueryLinksByTo(
		database.DB,
		constants.LinkTypeName(request.Type).ID().Int(),
		instance.Identity,
		linkSources,
		request.Limit,
	)
	if err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	links := make([]protocol.Link, 0)

	for _, linkModel := range linkModels {
		links = append(links, protocol.Link{
			DateCreated: linkModel.CreatedAt,
			From:        linkModel.From,
			To:          linkModel.To,
			Type:        constants.LinkTypeID(linkModel.Type).String(),
			Source:      constants.ProfileSourceID(linkModel.Source).Name().String(),
			Metadata: protocol.LinkMetadata{
				Network: cases.Title(language.English, cases.NoLower).String(constants.NetworkSymbolCrossbell.String()),
				Proof:   "TODO",
			},
		})
	}

	var dateUpdated sql.NullTime

	for _, link := range links {
		if !dateUpdated.Valid {
			dateUpdated.Valid = true
			dateUpdated.Time = link.DateCreated
		}

		if dateUpdated.Time.Before(link.DateCreated) {
			dateUpdated.Time = link.DateCreated
		}
	}

	if len(links) == 0 {
		_ = c.Error(api.ErrorNotFound)

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier:  fmt.Sprintf("%s/backlinks", rss3uri.New(instance)),
		DateUpdated: time.Now(),
		Total:       len(links),
		List:        links,
	})
}
