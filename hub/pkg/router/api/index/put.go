package index

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/verify"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/isotime"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PutIndexRequest struct{}

func PutIndexHandlerFunc(c *gin.Context) {
	w := web.Gin{C: c}

	// Parse the request
	platformInstance, err := middleware.GetInstance(c)
	if err != nil {
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	// Get the new put index file
	var indexFile protocol.Index
	if err = c.ShouldBind(&indexFile); err != nil {
		w.JSONResponse(http.StatusBadRequest, status.CodeInvalidParams, nil)

		return
	}

	// Check the signature
	jsonBytes, err := json.Marshal(indexFile)
	if err != nil {
		w.JSONResponse(http.StatusInternalServerError, status.CodeInvalidParams, nil)

		return
	}

	if ok, err := verify.Signature(jsonBytes, platformInstance.GetIdentity(), indexFile.Identifier); !ok || err != nil {
		w.JSONResponse(http.StatusBadRequest, status.CodeSignatureInvalid, nil)

		return
	}

	// Get the original index file
	oIndexFile, httpStatus, bizStatus := getIndexFile(platformInstance)
	if oIndexFile == nil {
		w.JSONResponse(httpStatus, bizStatus, nil)

		return
	}

	// Start the transaction
	tx := database.Instance.DB(context.Background()).Begin()
	defer tx.Rollback()

	httpStatus, bizStatus = runUpdateTask(updateAccount, tx, &indexFile, oIndexFile)
	if bizStatus != status.CodeOK {
		logger.Error(bizStatus)
		w.JSONResponse(httpStatus, bizStatus, nil)

		return
	}

	httpStatus, bizStatus = runUpdateTask(updateAccountPlatform, tx, &indexFile, oIndexFile)
	if bizStatus != status.CodeOK {
		logger.Error(bizStatus)
		w.JSONResponse(httpStatus, bizStatus, nil)

		return
	}

	httpStatus, bizStatus = runUpdateTask(updateSignature, tx, &indexFile, oIndexFile)
	if bizStatus != status.CodeOK {
		logger.Error(bizStatus)
		w.JSONResponse(httpStatus, bizStatus, nil)

		return
	}

	if err := tx.Commit().Error; err != nil {
		w.JSONResponse(http.StatusInternalServerError, status.CodeDatabaseError, nil)

		return
	}

	w.JSONResponse(http.StatusOK, status.CodeOK, indexFile)

	return
}

type UpdateTaskFunc func(db *gorm.DB, old, newIndex *protocol.Index) (int, status.Code)

func runUpdateTask(f UpdateTaskFunc, db *gorm.DB, oldIndex, newIndex *protocol.Index) (int, status.Code) {
	httpStatus, bizStatus := f(db, oldIndex, newIndex)

	return httpStatus, bizStatus
}

func updateSignature(db *gorm.DB, indexFile *protocol.Index, oIndexFile *protocol.Index) (int, status.Code) {
	u, err := rss3uri.Parse(indexFile.Identifier)
	if err != nil {
		logger.Error(err)

		return http.StatusInternalServerError, status.CodeInvalidInstance
	}

	needUpdate := false

	toUpdate := model.Signature{}

	// Parse file date
	oldDateUpdated, err := time.Parse(isotime.ISO8601, oIndexFile.DateUpdated)
	if err != nil {
		logger.Error(err)

		return http.StatusBadRequest, status.CodeFileFieldError
	}

	newDateUpdated, err := time.Parse(isotime.ISO8601, indexFile.DateUpdated)
	if err != nil {
		logger.Error(err)

		return http.StatusBadRequest, status.CodeFileFieldError
	}

	// TODO
	logger.Error(oldDateUpdated)
	logger.Error(newDateUpdated)
	logger.Error(oldDateUpdated.After(newDateUpdated))
	logger.Error(newDateUpdated.Before(time.Now().Add(-time.Hour * 2)))

	// oldDateUpdated > newDateUpdated || newDateUpdated < now - 2 hours
	if oldDateUpdated.After(newDateUpdated) || newDateUpdated.Before(time.Now().Add(-time.Hour*2)) {
		logger.Error("invalid date")

		return http.StatusBadRequest, status.CodeFileFieldError
	} else {
		needUpdate = true
		toUpdate.UpdatedAt = newDateUpdated
	}

	if needUpdate {
		if err := db.Model(&model.Signature{
			FileURI: fmt.Sprintf("%s@%s", u.Instance.GetIdentity(), u.Instance.GetSuffix()),
		}).UpdateColumns(&toUpdate).Error; err != nil {
			logger.Error(err)

			return http.StatusInternalServerError, status.CodeDatabaseError
		}
	}

	return http.StatusOK, status.CodeOK
}

func updateAccount(db *gorm.DB, indexFile *protocol.Index, oIndexFile *protocol.Index) (int, status.Code) {
	u, err := rss3uri.Parse(indexFile.Identifier)
	if err != nil {
		logger.Error(err)

		return http.StatusInternalServerError, status.CodeInvalidInstance
	}

	needUpdate := false

	toUpdate := model.Account{}

	// Get columns diff
	if *indexFile.Profile.Name != *oIndexFile.Profile.Name {
		needUpdate = true
		toUpdate.Name = database.WrapNullString(*indexFile.Profile.Name)
	}

	if *indexFile.Profile.Bio != *oIndexFile.Profile.Bio {
		needUpdate = true
		toUpdate.Bio = database.WrapNullString(*indexFile.Profile.Bio)
	}

	if !isArrayEqual(indexFile.Profile.Avatars, oIndexFile.Profile.Avatars) {
		needUpdate = true
		toUpdate.Avatars = indexFile.Profile.Avatars
	}

	if !isAttachmentsEqual(indexFile.Profile.Attachments, oIndexFile.Profile.Attachments) {
		needUpdate = true
		toUpdate.Attachments = indexFile.Profile.Attachments.ToDBStruct()
	}

	if needUpdate {
		if err := db.Model(&model.Account{
			ID:       u.Instance.GetIdentity(),
			Platform: int(constants.PlatformSymbol(u.Instance.GetSuffix()).ID()),
		}).Updates(&toUpdate).Error; err != nil {
			logger.Error(err)

			return http.StatusInternalServerError, status.CodeDatabaseError
		}
	}

	return http.StatusOK, status.CodeOK
}

// nolint:funlen,gocognit // TODO
func updateAccountPlatform(db *gorm.DB, indexFile *protocol.Index, oIndexFile *protocol.Index) (int, status.Code) {
	var toUpdate, toCreate, toDelete []model.AccountPlatform

	oldAccountMap := make(map[string]string)
	newAccountMap := make(map[string]string)

	for _, account := range indexFile.Profile.Accounts {
		newAccountMap[account.Identifier] = account.Signature
	}

	for _, account := range oIndexFile.Profile.Accounts {
		oldAccountMap[account.Identifier] = account.Signature
	}

	accountURI, err := rss3uri.Parse(indexFile.Identifier)
	if err != nil {
		logger.Error(err)

		return http.StatusInternalServerError, status.CodeInvalidInstance
	}

	for _, account := range indexFile.Profile.Accounts {
		signature, exist := oldAccountMap[account.Identifier]
		if !exist {
			platformURI, err := rss3uri.Parse(account.Identifier)
			if err != nil {
				logger.Error(err)

				return http.StatusInternalServerError, status.CodeInvalidInstance
			}

			toCreate = append(toCreate, model.AccountPlatform{
				AccountID:         accountURI.Instance.GetIdentity(),
				AccountPlatformID: int(constants.PlatformSymbol(accountURI.Instance.GetSuffix()).ID()),
				PlatformAccountID: platformURI.Instance.GetIdentity(),
				PlatformID:        int(constants.PlatformSymbol(platformURI.Instance.GetSuffix()).ID()),
			})

			continue
		}

		platformURI, err := rss3uri.Parse(account.Identifier)
		if err != nil {
			logger.Error(err)

			return http.StatusInternalServerError, status.CodeInvalidInstance
		}

		if signature != account.Signature {
			toUpdate = append(toUpdate, model.AccountPlatform{
				AccountID:         accountURI.Instance.GetIdentity(),
				AccountPlatformID: int(constants.PlatformSymbol(accountURI.Instance.GetSuffix()).ID()),
				PlatformAccountID: platformURI.Instance.GetIdentity(),
				PlatformID:        int(constants.PlatformSymbol(platformURI.Instance.GetSuffix()).ID()),
			})
		}
	}

	for _, account := range oIndexFile.Profile.Accounts {
		_, exist := newAccountMap[account.Identifier]
		if !exist {
			platformURI, err := rss3uri.Parse(account.Identifier)
			if err != nil {
				logger.Error(err)

				return http.StatusInternalServerError, status.CodeInvalidInstance
			}

			toDelete = append(toDelete, model.AccountPlatform{
				AccountID:         accountURI.Instance.GetIdentity(),
				AccountPlatformID: int(constants.PlatformSymbol(accountURI.Instance.GetSuffix()).ID()),
				PlatformAccountID: platformURI.Instance.GetIdentity(),
				PlatformID:        int(constants.PlatformSymbol(platformURI.Instance.GetSuffix()).ID()),
			})
		}
	}

	if len(toDelete) > 0 {
		// TODO
		for i := 0; i < len(toDelete); i++ {
			if err := db.Where(&toDelete[i]).Delete(&model.AccountPlatform{}).Error; err != nil {
				logger.Error(err)

				return http.StatusInternalServerError, status.CodeDatabaseError
			}
		}
	}

	if len(toUpdate) > 0 {
		if err := db.Updates(&toUpdate).Error; err != nil {
			logger.Error(err)

			return http.StatusInternalServerError, status.CodeDatabaseError
		}
	}

	if len(toCreate) > 0 {
		if err := db.Create(&toCreate).Error; err != nil {
			logger.Error(err)

			return http.StatusInternalServerError, status.CodeDatabaseError
		}
	}

	return http.StatusOK, status.CodeOK
}

func isArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func isAttachmentsEqual(a, b []protocol.IndexProfileAttachment) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		aa := a[i]
		bb := b[i]

		if aa.Type != bb.Type ||
			aa.Content != bb.Content ||
			aa.Address != bb.Address ||
			aa.MimeType != bb.MimeType ||
			aa.SizeInBytes != bb.SizeInBytes {
			return false
		}
	}

	return true
}
