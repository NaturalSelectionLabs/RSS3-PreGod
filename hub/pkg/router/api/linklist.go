package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
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
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidInstance, nil)

		return
	}

	if platformInstance.Prefix != constants.PrefixNameAccount || platformInstance.Platform != constants.PlatformSymbolEthereum {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidInstance, nil)

		return
	}

	// Begin a transaction
	tx := database.Instance.Tx(context.Background())
	defer tx.Rollback()

	account, err := database.Instance.QueryAccount(
		tx,
		platformInstance.GetIdentity(),
		int(constants.PlatformSymbol(platformInstance.GetSuffix()).ID()),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w := web.Gin{C: c}
			w.JSONResponse(http.StatusNotFound, status.CodeLinkListNotExist, nil)

			return
		}

		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	identifier := rss3uri.New(platformInstance).String()

	linkListFile := protocol.LinkList{
		ListSignedBase: protocol.ListSignedBase{
			SignedBase: protocol.SignedBase{
				Base: protocol.Base{
					Version:    protocol.Version,
					Identifier: fmt.Sprintf("%s/list/link/following/%d", identifier, request.PageIndex),
				},
			},
			IdentifierNext: func() string {
				if request.PageIndex == 0 {
					return ""
				}

				return fmt.Sprintf("%s/list/link/following/%d", identifier, request.PageIndex-1)
			}(),
		},
	}

	// TODO Define following type id
	links, err := database.Instance.QueryLinks(tx, 1, account.ID, account.Platform, request.PageIndex)
	if err != nil {
		w := web.Gin{C: c}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.JSONResponse(http.StatusNotFound, status.CodeLinkNotExist, nil)
		}

		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	for _, link := range links {
		linkListFile.List = append(linkListFile.List, protocol.LinkListItem{
			Type: constants.LinkTypeFollowing.String(),
			// TODO  Maybe it's an asset or a note
			IdentifierTarget: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixNameAccount,
				Identity: link.TargetIdentity,
				Platform: constants.PlatformSymbolEthereum,
			}).String(),
		})
	}

	signature, err := database.Instance.QuerySignature(
		tx,
		fmt.Sprintf(
			"%s@%s/list/link/following/%d",
			platformInstance.GetIdentity(),
			constants.PlatformSymbol(platformInstance.GetSuffix()),
			request.PageIndex,
		),
	)

	if err != nil {
		w := web.Gin{C: c}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.JSONResponse(http.StatusNotFound, status.CodeSignatureNotExist, nil)
		}

		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	if err := tx.Commit().Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	linkListFile.Signature = signature.Signature
	linkListFile.Base.DateCreated = signature.CreatedAt.Format(isotime.ISO8601)
	linkListFile.Base.DateUpdated = signature.UpdatedAt.Format(isotime.ISO8601)

	linkListFile.Total = len(linkListFile.List)

	c.JSON(http.StatusOK, &linkListFile)
}
