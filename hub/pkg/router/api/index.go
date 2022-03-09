package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
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

	index := protocol.IndexFile{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:     protocol.Version,
				Identifier:  rss3uri.New(platformInstance).String(),
				DateCreated: account.CreatedAt.Format(time.RFC3339),
				DateUpdated: account.UpdatedAt.Format(time.RFC3339),
			},
		},
		Agent: protocol.Agent{},
		Profile: protocol.Profile{
			Name:    account.Name,
			Avatars: account.Avatars,
			Bio:     account.Bio,
		},
		Links: protocol.Links{},
		Items: protocol.Items{},
	}

	var accountPlatforms []model.AccountPlatform
	if err := db.DB.Where("account_id = ?", account.AccountID).Find(&accountPlatforms).Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.CodeError, nil)

		return
	}

	for _, platform := range accountPlatforms {
		account := protocol.Account{
			Identifier: rss3uri.New(&rss3uri.PlatformInstance{
				Prefix:   constants.PrefixNameAccount,
				Identity: platform.PlatformAccountID,
				Platform: constants.PlatformSymbolEthereum,
			}).String(),
		}
		index.Profile.Accounts = append(index.Profile.Accounts, account)
	}

	index.Links.Identifiers = append(index.Links.Identifiers, protocol.LinkIdentifier{
		Type: constants.LinkTypeFollowing.String(),
		// TODO Refine rss3uri package
		// TODO Max page index
		IdentifierCustom: fmt.Sprintf("%s/list/link/following/1", rss3uri.New(platformInstance).String()),
		Identifier:       fmt.Sprintf("%s/list/link/following", rss3uri.New(platformInstance).String()),
	})

	index.Items.Notes = protocol.Notes{
		Identifier: fmt.Sprintf("%s/list/notes", rss3uri.New(platformInstance).String()),
	}

	index.Items.Assets = protocol.Assets{
		Identifier: fmt.Sprintf("%s/list/assets", rss3uri.New(platformInstance).String()),
	}

	c.JSON(http.StatusOK, index)
}
