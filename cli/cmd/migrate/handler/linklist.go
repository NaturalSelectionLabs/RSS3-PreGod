package handler

import (
	"strconv"
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
)

func MigrateLinkList(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Migrate signature
		if file.Content.Signature != "" {
			if err := tx.Create(&model.Signature{
				FileURI:   strings.ReplaceAll(strings.ReplaceAll(file.Path, "-", "/"), ".", "/"),
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

		splits := strings.Split(file.Path, "-")
		pageIndex, err := strconv.Atoi(splits[len(splits)-1])
		if err != nil {
			return err
		}

		for _, targetIdentity := range file.Content.List {
			if err := tx.Create(&model.Link{
				Type:           1,         // Following
				Identity:       splits[0], // Ethereum wallet address
				PrefixID:       int(constants.PrefixIDAccount),
				SuffixID:       int(constants.PlatformIDEthereum),
				TargetIdentity: targetIdentity,
				TargetPrefixID: int(constants.PrefixIDAccount),
				TargetSuffixID: int(constants.PlatformIDEthereum),
				PageIndex:      pageIndex,
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			}).Error; err != nil {
				return err
			}

			atomic.AddInt64(&stats.LinkNumber, 1)
		}

		return nil
	})
}
