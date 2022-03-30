package index

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

type GetIndexRequest struct{}

func GetIndexHandlerFunc(c *gin.Context) {
	w := web.Gin{C: c}

	value, exists := c.Get(middleware.KeyInstance)
	if !exists {
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	platformInstance, ok := value.(*rss3uri.PlatformInstance)
	if !ok {
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidInstance, nil)

		return
	}

	if platformInstance.Prefix != constants.PrefixNameAccount || platformInstance.Platform != constants.PlatformSymbolEthereum {
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidInstance, nil)

		return
	}

	indexFile, httpStatus, bizStatus := getIndexFile(platformInstance)

	w.JSONResponse(httpStatus, bizStatus, *indexFile)
}

// Gets the index file for the given instance.
// Returns:
//  - protocol.Index: The index file.
//  - int: The HTTP status code.
//  - status.Code: The business status code.
func getIndexFile(instance *rss3uri.PlatformInstance) (*protocol.Index, int, status.Code) {
	// Setup index file
	indexFile := protocol.NewIndex(instance)

	// Start the transaction
	tx := database.Instance.Tx(context.Background())
	defer tx.Rollback()

	// Query the account
	account, err := database.Instance.QueryAccount(
		tx, instance.GetIdentity(), int(constants.PlatformSymbol(instance.GetSuffix()).ID()),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusConflict, status.CodeIndexNotExist
		}

		return nil, http.StatusInternalServerError, status.CodeDatabaseError
	}

	if account.Name.Valid {
		indexFile.Profile.Name = &account.Name.String
	}

	indexFile.Profile.Avatars = account.Avatars

	if account.Bio.Valid {
		indexFile.Profile.Bio = &account.Bio.String
	}

	// Query the linklists
	linkLists, err := database.Instance.QueryLinkListsByOwner(
		tx, account.ID, int(constants.PrefixIDAccount), account.Platform,
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, status.CodeDatabaseError
	}

	for _, linkList := range linkLists {
		indexFile.AddLinkIdentifier(constants.LinkTypeID(linkList.Type), linkList.MaxPageIndex)
	}

	// Query the account platforms
	accountPlatforms, err := database.Instance.QueryAccountPlatforms(tx, account.ID, account.Platform)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, status.CodeDatabaseError
	}

	for _, accountPlatform := range accountPlatforms {
		indexFile.AddProfileAccount(
			accountPlatform.PlatformAccountID, constants.PlatformID(accountPlatform.PlatformID), protocol.SignatureNone,
		)
	}

	// Query the signature
	signature, err := database.Instance.QuerySignature(
		tx, fmt.Sprintf("%s@%s", account.ID, constants.PlatformID(account.Platform).Symbol()),
	)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, http.StatusInternalServerError, status.CodeDatabaseError
	}

	indexFile.Signature = signature.Signature
	indexFile.Base.DateCreated = signature.CreatedAt.Format(isotime.ISO8601)
	indexFile.Base.DateUpdated = signature.UpdatedAt.Format(isotime.ISO8601)

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, http.StatusInternalServerError, status.CodeDatabaseError
	}

	return indexFile, http.StatusOK, status.CodeOK
}
