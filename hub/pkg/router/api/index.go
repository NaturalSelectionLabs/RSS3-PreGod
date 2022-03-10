package api

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol/file"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
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

	account := model.Account{}
	if err := db.DB.Where(
		"account_id = ?",
		fmt.Sprintf("%s@%s", platformInstance.GetIdentity(), platformInstance.GetSuffix()),
	).First(&account).Error; err != nil {
		// TODO Account not found
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//}
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	// Query the max page index
	var followingPageIndex int
	if err := db.DB.Table("link").Select("max(page_index)").Where("rss3_id = ?", platformInstance.Identity).Row().Scan(&followingPageIndex); err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	// TODO Redesign database table
	var (
		notePageIndex  int
		assetPageIndex int
	)

	identifier := rss3uri.New(platformInstance).String()

	indexFile := file.Index{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:    protocol.Version,
				Identifier: identifier,
			},
		},
		Profile: file.IndexProfile{
			Name:    account.Name,
			Avatars: account.Avatars,
			Bio:     account.Bio,
			// TODO No data available
			// Attachments: nil,
		},
		Links: file.IndexLinks{
			Identifiers: []file.IndexLinkIdentifier{
				{
					Type:             "following",
					IdentifierCustom: fmt.Sprintf("%s/list/link/following/%d", identifier, followingPageIndex),
					Identifier:       fmt.Sprintf("%s/list/link/following", identifier),
				},
			},
			IdentifierBack: fmt.Sprintf("%s/list/backlink", identifier),
		},
		Items: file.IndexItems{
			Notes: file.IndexItemsNotes{
				IdentifierCustom: fmt.Sprintf("%s/list/note/%d", identifier, notePageIndex),
				Identifier:       fmt.Sprintf("%s/list/note", identifier),
			},
			Assets: file.IndexItemsAssets{
				IdentifierCustom: fmt.Sprintf("%s/list/asset/%d", identifier, assetPageIndex),
				Identifier:       fmt.Sprintf("%s/list/asset", identifier),
			},
		},
	}

	var accountPlatforms []model.AccountPlatform
	if err := db.DB.Where("account_id = ?", account.AccountID).Find(&accountPlatforms).Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	for _, accountPlatform := range accountPlatforms {
		indexFile.Profile.Accounts = append(indexFile.Profile.Accounts, file.IndexAccount{
			Identifier: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixNameAccount,
				Identity: accountPlatform.PlatformAccountID,
				Platform: accountPlatform.PlatformID.Symbol(),
			}).String(),
		})
	}

	c.JSON(http.StatusOK, &indexFile)
}
