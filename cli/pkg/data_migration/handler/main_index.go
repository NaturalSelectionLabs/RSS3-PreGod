package handler

import (
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/pkg/data_migration/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

func getAccountPlatform(platform string) constants.PlatformNameID {
	switch strings.ToLower(platform) {
	case "evm+":
		return constants.PlatformNameID_Evm
	case "solana":
		return constants.PlatformNameID_Solana
	case "flow":
		return constants.PlatformNameID_Flow
	case "twitter":
		return constants.PlatformNameID_Twitter
	case "misskey":
		return constants.PlatformNameID_Misskey
	case "jike":
		return constants.PlatformNameID_Jike
	default:
		return constants.PlatformNameID_Unknown
	}
}

//nolint:funlen // no need to split
func MainIndex(content bson.D) error {
	// handle main index
	var mainIndex protocol.RSS3Index031
	// Unmarshal
	doc, err := bson.Marshal(content)
	if err != nil {
		return err
	}

	if err = bson.Unmarshal(doc, &mainIndex); err != nil {
		return err
	}

	// Split & save into db

	var instanceBase model.InstanceBase

	newID := mainIndex.ID + "@" + string(constants.PlatformName_Evm)

	CreatedAt, err := time.Parse("2006-01-02T15:04:05.000Z", mainIndex.DateCreated)
	if err != nil {
		return err
	}

	UpdatedAt, err := time.Parse("2006-01-02T15:04:05.000Z", mainIndex.DateUpdated)
	if err != nil {
		return err
	}

	// if err := db.DB.First(&instanceBase, "rss3_id = ?", newID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	// Already exists
	// 	return nil // skip
	// }

	// New instance
	instanceBase = model.InstanceBase{
		RSS3ID:         newID,
		PrefixID:       constants.PrefixID_Account,
		InstanceTypeID: constants.InstanceType_Account,
		Base: model.BaseModel{
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		},
	}

	// Accounts
	accounts := []model.AccountPlatform{
		{
			AccountID:         newID,
			PlatformNameID:    constants.PlatformNameID_Evm,
			PlatformAccountID: mainIndex.ID,
			Base: model.BaseModel{
				CreatedAt: CreatedAt,
				UpdatedAt: UpdatedAt,
			},
		},
	}

	for _, additionalAccount := range mainIndex.Profile.Accounts {
		splits := strings.Split(additionalAccount.ID, "-")

		accounts = append(accounts, model.AccountPlatform{
			AccountID:         splits[0] + "@" + splits[1],
			PlatformNameID:    getAccountPlatform(splits[0]),
			PlatformAccountID: splits[1],
			Base: model.BaseModel{
				CreatedAt: CreatedAt,
				UpdatedAt: UpdatedAt,
			},
		})
	}

	// Account
	account := &model.Account{
		AccountID: newID,
		Name:      mainIndex.Profile.Name,
		Bio:       mainIndex.Profile.Bio,
		Avatars:   mainIndex.Profile.Avatar,

		InstanceBase:    instanceBase,
		AccountPlatform: accounts,

		Base: model.BaseModel{
			CreatedAt: CreatedAt,
			UpdatedAt: UpdatedAt,
		},
	}

	// todo: migrate application data (signature, account tags, custom fields, etc)

	// Save
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&instanceBase).Error; err != nil {
			return err
		}
		if err := tx.Save(&accounts).Error; err != nil {
			return err
		}
		if err := tx.Create(&account).Error; err != nil {
			return err
		}

		return nil
	})
}
