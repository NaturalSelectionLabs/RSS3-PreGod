package handler

import (
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
)

func MigrateIndex(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Migrate signature
		if file.Content.Signature != "" {
			if err := tx.Create(&model.Signature{
				FileURI:   file.Path,
				Signature: file.Content.Signature,
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			}).Error; err != nil {
				return err
			}

			atomic.AddInt64(&stats.SignatureNumber, 1)
		}

		// Migrate ethereum account
		if err := tx.Create(&model.Account{
			ID:       file.Path,
			Platform: int(constants.PlatformIDEthereum),
			Name:     file.Content.Profile.Name,
			Bio:      file.Content.Profile.Bio,
			Avatars:  file.Content.Profile.Avatar,
			Table: common.Table{
				CreatedAt: file.Content.DateCreated,
				UpdatedAt: file.Content.DateUpdated,
			},
		}).Error; err != nil {
			return err
		}

		atomic.AddInt64(&stats.AccountNumber, 1)

		// Migrate platform account
		for _, account := range file.Content.Profile.Accounts {
			splits := strings.Split(account.ID, "-")
			platform := splits[0]
			platformID := int(constants.PlatformSymbol(strings.ToLower(platform)).ID())
			if platformID == 0 {
				platformID = int(constants.PlatformIDEthereum)
			}

			accountID := splits[1]
			if err := tx.Create(&model.AccountPlatform{
				AccountID:         file.Content.ID,
				AccountPlatformID: int(constants.PlatformIDEthereum),
				PlatformAccountID: strings.Trim(strings.Trim(accountID, "@"), "\\"),
				PlatformID:        platformID,
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			}).Error; err != nil {
				return err
			}

			atomic.AddInt64(&stats.AccountPlatformNumber, 1)
		}

		return nil
	})
}
