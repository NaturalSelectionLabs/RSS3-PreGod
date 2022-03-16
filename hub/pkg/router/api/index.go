package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetIndexRequest struct{}

//nolint:funlen // SQL logic will be wrapped up later
func GetIndexHandlerFunc(c *gin.Context) {
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

	// Query the max page index
	followingMaxPageIndex, err := database.Instance.QueryLinkWithMaxPageIndex(tx, 1, account.ID, account.Platform)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	// TODO No test data available
	var (
		notePageIndex  int
		assetPageIndex int
	)

	identifier := rss3uri.New(platformInstance).String()

	indexFile := protocol.Index{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:    protocol.Version,
				Identifier: identifier,
			},
		},
		Profile: protocol.IndexProfile{
			Name:    account.Name,
			Avatars: account.Avatars,
			Bio:     account.Bio,
			// TODO No data available
			// Attachments: nil,
		},
		Links: protocol.IndexLinks{
			Identifiers: []protocol.IndexLinkIdentifier{
				{
					Type:             "following",
					IdentifierCustom: fmt.Sprintf("%s/list/link/following/%d", identifier, followingMaxPageIndex),
					Identifier:       fmt.Sprintf("%s/list/link/following", identifier),
				},
			},
			IdentifierBack: fmt.Sprintf("%s/list/backlink", identifier),
		},
		Items: protocol.IndexItems{
			Notes: protocol.IndexItemsNotes{
				IdentifierCustom: fmt.Sprintf("%s/list/note/%d", identifier, notePageIndex),
				Identifier:       fmt.Sprintf("%s/list/note", identifier),
			},
			Assets: protocol.IndexItemsAssets{
				IdentifierCustom: fmt.Sprintf("%s/list/asset/%d", identifier, assetPageIndex),
				Identifier:       fmt.Sprintf("%s/list/asset", identifier),
			},
		},
	}

	accountPlatforms, err := database.Instance.QueryAccountPlatforms(tx, account.ID, account.Platform)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	for _, accountPlatform := range accountPlatforms {
		indexFile.Profile.Accounts = append(indexFile.Profile.Accounts, protocol.IndexAccount{
			Identifier: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixNameAccount,
				Identity: accountPlatform.PlatformAccountID,
				Platform: constants.PlatformID(accountPlatform.PlatformID).Symbol(),
			}).String(),
		})
	}

	signature, err := database.Instance.QuerySignature(
		tx,
		fmt.Sprintf("%s@%s", account.ID, constants.PlatformID(account.Platform).Symbol()),
	)

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

	indexFile.Signature = signature.Signature
	indexFile.Base.DateCreated = signature.CreatedAt.Format(time.RFC3339)
	indexFile.Base.DateUpdated = signature.UpdatedAt.Format(time.RFC3339)

	c.JSON(http.StatusOK, &indexFile)
}
