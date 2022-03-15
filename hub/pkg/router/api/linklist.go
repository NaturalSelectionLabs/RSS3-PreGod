package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol/file"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetLinkListRequest struct {
	LinkType  string `uri:"link_type" binding:"required"`
	PageIndex int    `uri:"page_index"`
}

//nolint:funlen // SQL logic will be wrapped up later
func GetLinkListHandlerFunc(c *gin.Context) {
	request := GetLinkListRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	// TODO Handle other types of requests
	if request.LinkType != "following" {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	value, exists := c.Get(middleware.KeyInstance)
	if !exists {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	platformInstance, ok := value.(*rss3uri.PlatformInstance)
	if !ok {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	if platformInstance.Prefix != constants.PrefixNameAccount || platformInstance.Platform != constants.PlatformSymbolEthereum {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	// TODO Check if the account exists

	identifier := rss3uri.New(platformInstance).String()

	linkListFile := file.LinkList{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:    protocol.Version,
				Identifier: fmt.Sprintf("%s/list/link/following/%d", identifier, request.PageIndex),
				// TODO IdentifierNext
				// TODO No test data available
				// DateCreated: "",
				// DateUpdated: "",
			},
		},
	}

	var links []model.Link
	if err := database.Instance.DB(context.Background()).Where(
		"rss3_id = ? and page_index = ?",
		platformInstance.GetIdentity(),
		request.PageIndex,
	).Find(&links).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// TODO Return 404 not found?
			return
		}

		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	for _, link := range links {
		linkListFile.List = append(linkListFile.List, file.LinkListItem{
			Type: constants.LinkTypeFollowing.String(),
			// TODO  Maybe it's an asset or a note
			IdentifierTarget: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixNameAccount,
				Identity: link.TargetRSS3ID,
				Platform: constants.PlatformSymbolEthereum,
			}).String(),
		})
	}

	linkListFile.Total = len(linkListFile.List)

	c.JSON(http.StatusOK, &linkListFile)
}
