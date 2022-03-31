package database

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Setup() error {
	db, err := gorm.Open(postgres.New(postgres.Config{
		// TODO Refactor config package
		DSN: config.Config.Postgres.DSN,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.New(),
	})

	if err != nil {
		return err
	}

	DB = db

	// Install uuid extension for postgres
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}

	if err := DB.AutoMigrate(
	//&model.Account{},
	//&model.AccountPlatform{},
	//&model.Instance{},
	//&model.LinkList{},
	//&model.Link{},
	//&model.Signature{},
	//&model.Asset{},
	//&model.Note{},
	); err != nil {
		return err
	}

	return nil
}

func QueryProfiles(db *gorm.DB, id string, platform int) ([]model.Profile, error) {
	var profiles []model.Profile
	if err := db.Where(&model.Profile{
		ID:       id,
		Platform: platform,
	}).Find(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}

func QueryAccounts(db *gorm.DB, profileID string, profilePlatform int, source int) ([]model.Account, error) {
	var accounts []model.Account
	if err := db.Where(&model.Account{
		ProfileID:       profileID,
		ProfilePlatform: profilePlatform,
		Source:          source,
	}).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}

func QueryLinks(db *gorm.DB, _type int, form string, source int) ([]model.Link, error) {
	var links []model.Link
	if err := db.Where(&model.Link{
		Type:   _type,
		From:   form,
		Source: source,
	}).Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func QueryLinksByTo(db *gorm.DB, _type int, to string, source int) ([]model.Link, error) {
	var links []model.Link
	if err := db.Where(&model.Link{
		Type:   _type,
		To:     to,
		Source: source,
	}).Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}
