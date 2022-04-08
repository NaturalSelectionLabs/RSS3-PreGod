package handler

import (
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
)

func MigrateIndex(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Migrate crossbell account
		attachments := make(datatype.Attachments, 0)

		// https://stackoverflow.com/questions/38933898/error-parsing-regexp-invalid-or-unsupported-perl-syntax
		re := regexp.MustCompile(`<SITE#\S+>$`)
		result := re.FindString(file.Content.Profile.Bio)

		site := strings.TrimRight(strings.TrimLeft(result, "<SITE#"), ">")
		uris := fmt.Sprintf("https://%s", site)
		bio := strings.TrimRight(file.Content.Profile.Bio, fmt.Sprintf("<SITE#%s+>", site))

		if site != "" {
			attachments = append(attachments, datatype.Attachment{
				Type:     "website",
				Content:  uris,
				MimeType: "text/uri-list",
			})
		}

		tx.Create(&model.Profile{
			ID:          file.Path,
			Platform:    int(constants.PlatformIDEthereum),
			Name:        database.WrapNullString(file.Content.Profile.Name),
			Bio:         database.WrapNullString(bio),
			Avatars:     file.Content.Profile.Avatar,
			Attachments: attachments,
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

			for i, split := range splits[1:] {
				splits[i+1] = strings.Trim(split, "\\")
			}

			accountID := strings.Join(splits[1:], "-")

			if err := tx.Create(&model.Account{
				Identity:        strings.Trim(accountID, "@"),
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
