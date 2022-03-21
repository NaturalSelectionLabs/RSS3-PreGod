package handler

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	mongomodel "github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/cmd/migrate/stats"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func MigrateLinkList(db *gorm.DB, file mongomodel.File) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Field format update
		fileURI := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(file.Path, "-", "/"), ".", "/"), "links", "link")

		// Migrate signature
		if file.Content.Signature != "" {
			// Field format update
			identity := strings.Split(fileURI, "/")[0]
			fileURI = strings.ReplaceAll(fileURI, identity, fmt.Sprintf("%s@ethereum", identity))

			if err := tx.Create(&model.Signature{
				FileURI:   fileURI,
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

		linkList := model.LinkList{
			Type:         1, // Following
			Identity:     splits[0],
			PrefixID:     int(constants.PrefixIDAccount),
			SuffixID:     int(constants.PlatformIDEthereum),
			ItemCount:    len(file.Content.List),
			MaxPageIndex: 0, // TODO Page break
			Table: common.Table{
				CreatedAt: file.Content.DateCreated,
				UpdatedAt: file.Content.DateUpdated,
			},
		}

		// Returning all fields
		if err := tx.Clauses(clause.Returning{}).Create(&linkList).Error; err != nil {
			return err
		}

		links := make([]model.Link, len(file.Content.Links))
		for _, targetIdentity := range file.Content.List {
			links = append(links, model.Link{
				Type:           1,         // Following
				Identity:       splits[0], // Ethereum wallet address
				PrefixID:       int(constants.PrefixIDAccount),
				SuffixID:       int(constants.PlatformIDEthereum),
				TargetIdentity: targetIdentity,
				TargetPrefixID: int(constants.PrefixIDAccount),
				TargetSuffixID: int(constants.PlatformIDEthereum),
				LinkListID:     linkList.ID,
				PageIndex:      pageIndex,
				Table: common.Table{
					CreatedAt: file.Content.DateCreated,
					UpdatedAt: file.Content.DateUpdated,
				},
			})
		}

		if err := tx.CreateInBatches(links, 1024).Error; err != nil {
			return err
		}

		atomic.AddInt64(&stats.LinkNumber, int64(len(links)))

		return nil
	})
}
