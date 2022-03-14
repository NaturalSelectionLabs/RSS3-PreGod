package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetIndexRequest struct {
	Instance string `uri:"instance" binding:"required"`
}

//nolint:funlen // SQL logic will be wrapped up later
func GetIndexHandlerFunc(c *gin.Context) {
	request := GetIndexRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	instance, err := rss3uri.ParseInstance(request.Instance)
	if err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	accountInstance, ok := instance.(*rss3uri.PlatformInstance)
	if !ok {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	if accountInstance.Prefix != constants.PrefixNameAccount || accountInstance.Platform != constants.PlatformSymbolEthereum {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	account := model.Account{}
	if err := database.Instance.DB(context.Background()).Where(
		"id = ?",
		fmt.Sprintf("%s@%s", instance.GetIdentity(), instance.GetSuffix()),
	).First(&account).Error; err != nil {
		// TODO Account not found
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//}
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.ERROR, nil)

		return
	}

	index := protocol.IndexFile{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version:     protocol.Version,
				Identifier:  rss3uri.New(instance).String(),
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
	if err := database.Instance.DB(context.Background()).Where("account_id = ?", account.ID).Find(&accountPlatforms).Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.ERROR, nil)

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
		IdentifierCustom: fmt.Sprintf("%s/list/link/following/1", rss3uri.New(accountInstance).String()),
		Identifier:       fmt.Sprintf("%s/list/link/following", rss3uri.New(accountInstance).String()),
	})

	index.Items.Notes = protocol.Notes{
		Identifier: fmt.Sprintf("%s/list/notes", rss3uri.New(accountInstance).String()),
	}

	index.Items.Assets = protocol.Assets{
		Identifier: fmt.Sprintf("%s/list/assets", rss3uri.New(accountInstance).String()),
	}

	c.JSON(http.StatusOK, index)
}
