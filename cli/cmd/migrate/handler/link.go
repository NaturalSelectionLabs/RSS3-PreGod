package handler

import (
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MigrateLinkList(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		splits := strings.Split(file.Path, "-")

		links := make([]model.Link, 0, len(file.Content.Links))
		for _, targetIdentity := range file.Content.List {
			links = append(links, model.Link{
				Type:             constants.LinkTypeFollow.Int(),
				From:             splits[0],
				FromInstanceType: int(constants.InstanceTypeAccount),
				FromPlatformID:   constants.PlatformIDEthereum.Int(),
				To:               targetIdentity,
				ToInstanceType:   int(constants.InstanceTypeAccount),
				ToPlatformID:     constants.PlatformIDEthereum.Int(),
				Source:           constants.ProfileSourceIDCrossbell.Int(),
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			})
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "type"},
				{Name: "from"},
				{Name: "from_instance_type"},
				{Name: "from_platform_id"},
				{Name: "to"},
				{Name: "to_instance_type"},
				{Name: "to_platform_id"},
				{Name: "source"},
			},
			UpdateAll: true,
		}).CreateInBatches(links, 1024).Error; err != nil {
			return err
		}

		atomic.AddInt64(&stats.Link, int64(len(links)))

		return nil
	})
}
