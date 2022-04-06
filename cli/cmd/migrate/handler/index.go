package handler

import (
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
)

func MigrateIndex(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Migrate crossbell account
		tx.Create(&model.Profile{
			ID:          file.Path,
			Platform:    int(constants.PlatformIDEthereum),
			Name:        database.WrapNullString(file.Content.Profile.Name),
			Bio:         database.WrapNullString(file.Content.Profile.Bio),
			Avatars:     file.Content.Profile.Avatar,
			Attachments: nil,
			Table: common.Table{
				CreatedAt: file.Content.DateCreated,
				UpdatedAt: file.Content.DateUpdated,
			},
		})

		atomic.AddInt64(&stats.Profile, 1)

		// Migrate connected accounts
		for _, account := range file.Content.Profile.Accounts {
			splits := strings.Split(account.ID, "-")
			platform := splits[0]
			platformID := int(constants.PlatformSymbol(strings.ToLower(platform)).ID())
			if platformID == 0 {
				platformID = int(constants.PlatformIDEthereum)
			}

			accountID := splits[1]
			if err := tx.Create(&model.Account{
				ID:              strings.Trim(strings.Trim(accountID, "@"), "\\"),
				Platform:        platformID,
				ProfileID:       file.Content.ID,
				ProfilePlatform: int(constants.PlatformIDEthereum),
				Source:          constants.ProfileSourceIDCrossbell.Int(),
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			}).Error; err != nil {
				return err
			}

			atomic.AddInt64(&stats.Account, 1)
		}

		return nil
	})
}
