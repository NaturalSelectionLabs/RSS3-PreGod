package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol/file"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetBackLinkListRequest struct{}

//nolint:funlen // SQL logic will be wrapped up later
func GetBackLinkListHandlerFunc(c *gin.Context) {
	var (
		limit        = 0
		instance     = c.Query("instance")
		lastInstance = c.Query("lastInstance")
	)

	log.Println(limit)

	if c.Query("limit") != "" {
		var err error

		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			w := web.Gin{C: c}
			w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

			return
		}
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
			w.JSONResponse(http.StatusNotFound, status.CodeError, nil)

			return
		}

		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	identifier := rss3uri.New(platformInstance).String()

	backLinkListFile := file.BackLinkList{
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
	links, err := database.Instance.QueryLinksByTarget(tx, 1, account.ID, account.Platform, limit, instance, lastInstance)
	if err != nil {
		w := web.Gin{C: c}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.JSONResponse(http.StatusNotFound, status.CodeError, nil)
		}

		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

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

		backLinkListFile.List = append(backLinkListFile.List, file.LinkListItem{
			Type: constants.LinkTypeFollowing.String(),
			// TODO  Maybe it's an asset or a note
			IdentifierTarget: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixID(link.PrefixID).String(),
				Identity: link.Identity,
				Platform: constants.PlatformID(link.SuffixID).Symbol(),
			}).String(),
		})
	}

	if err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	if err := tx.Commit().Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	backLinkListFile.Total = len(backLinkListFile.List)
	backLinkListFile.DateCreated = dateCreated.Time.Format(time.RFC3339)
	backLinkListFile.DateUpdated = dateUpdated.Time.Format(time.RFC3339)

	c.JSON(http.StatusOK, &backLinkListFile)
}
