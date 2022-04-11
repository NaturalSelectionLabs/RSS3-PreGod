package handler

import (
	"strings"
	"sync/atomic"
	"time"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
)

func MigrateLinkList(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		splits := strings.Split(file.Path, "-")

		links := make([]model.Link, 0, len(file.Content.Links))
		length := len(file.Content.List)
		for i, targetIdentity := range file.Content.List {
			links = append(links, model.Link{
				Type:             constants.LinkTypeFollow.Int(),
				From:             strings.ToLower(splits[0]),
				FromInstanceType: int(constants.InstanceTypeAccount),
				FromPlatformID:   constants.PlatformIDEthereum.Int(),
				To:               strings.ToLower(targetIdentity),
				ToInstanceType:   int(constants.InstanceTypeAccount),
				ToPlatformID:     constants.PlatformIDEthereum.Int(),
				Source:           constants.ProfileSourceIDCrossbell.Int(),
				Table: common.Table{
					CreatedAt: file.Content.DateCreated.Add(time.Duration(length-i) * time.Millisecond),
					UpdatedAt: file.Content.DateUpdated.Add(time.Duration(length-i) * time.Millisecond),
				},
			})
		}

		if err := tx.CreateInBatches(links, 1024).Error; err != nil {
			return err
		}

		atomic.AddInt64(&stats.Link, int64(len(links)))

		return nil
	})
}
