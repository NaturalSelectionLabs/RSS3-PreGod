package api

import (
	"context"
	"database/sql"
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

type GetBackLinkListRequest struct {
	Limit        int    `form:"limit"`
	Instance     string `form:"instance"`
	LastInstance string `form:"last_instance"`
}

//nolint:funlen // SQL logic will be wrapped up later
func GetBackLinkListHandlerFunc(c *gin.Context) {
	request := GetBackLinkListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
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

	backLinkListFile := protocol.BackLinkList{
		ListUnsignedBase: protocol.ListUnsignedBase{
			UnsignedBase: protocol.UnsignedBase{
				Base: protocol.Base{
					Version:    protocol.Version,
					Identifier: fmt.Sprintf("%s/list/backlink", identifier),
				},
			},
		},
	}

	// TODO Define following type id
	links, err := database.Instance.QueryLinksByTarget(tx, 1, account.ID, account.Platform, request.Limit, request.Instance, request.LastInstance)
	if err != nil {
		w := web.Gin{C: c}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.JSONResponse(http.StatusNotFound, status.CodeLinkNotExist, nil)
		}

		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	var (
		dateCreated sql.NullTime
		dateUpdated sql.NullTime
	)

	for _, link := range links {
		if !dateCreated.Valid || link.CreatedAt.After(dateCreated.Time) {
			dateCreated.Time = link.CreatedAt
		}

		if !dateUpdated.Valid || link.CreatedAt.After(dateCreated.Time) {
			dateUpdated.Time = link.UpdatedAt
		}

		backLinkListFile.List = append(backLinkListFile.List, protocol.LinkListItem{
			Type: constants.LinkTypeFollowing.String(),
			// TODO  Maybe it's an asset or a note
			IdentifierTarget: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixID(link.PrefixID).String(),
				Identity: link.Identity,
				Platform: constants.PlatformID(link.SuffixID).Symbol(),
			}).String(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	backLinkListFile.Total = len(backLinkListFile.List)
	backLinkListFile.DateCreated = dateCreated.Time.Format(isotime.ISO8601)
	backLinkListFile.DateUpdated = dateUpdated.Time.Format(isotime.ISO8601)

	c.JSON(http.StatusOK, &backLinkListFile)
}
