package database

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		&model.Profile{},
		&model.Account{},
		&model.Link{},
		&model.Asset{},
		&model.Note{},
	); err != nil {
		return err
	}

	return nil
}

func QueryInstance(db *gorm.DB, id string, platform int) error {
	_, err := QueryProfiles(db, id, platform, []int{})

	return err
}

func QueryProfiles(db *gorm.DB, id string, platform int, profileSources []int) ([]model.Profile, error) {
	var profiles []model.Profile

	internalDB := db.Where(&model.Profile{
		ID:       id,
		Platform: platform,
	})

	if len(profileSources) > 0 {
		internalDB.Where("source IN ?", profileSources)
	}

	if err := internalDB.Find(&profiles).Error; err != nil {
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

func QueryLinks(db *gorm.DB, _type int, form string, linkSources []int, limit int) ([]model.Link, error) {
	var links []model.Link

	internalDB := db.Where(&model.Link{
		Type: _type,
		From: form,
	})

	if len(linkSources) > 0 {
		internalDB.Where("source IN ?", linkSources)
	}

	if limit > 0 {
		internalDB = internalDB.Limit(limit)
	}

	if err := internalDB.Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func QueryLinksByTo(db *gorm.DB, _type int, to string, linkSources []int, limit int) ([]model.Link, error) {
	var links []model.Link

	internalDB := db.Where(&model.Link{
		Type: _type,
		To:   to,
	})

	if len(linkSources) > 0 {
		internalDB.Where("source IN ?", linkSources)
	}

	if limit > 0 {
		internalDB = internalDB.Limit(limit)
	}

	if err := internalDB.Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func CreateNote(db *gorm.DB, note *model.Note) (*model.Note, error) {
	if err := db.Model(note).Clauses(clause.Returning{}).Create(&note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func CreateAsset(db *gorm.DB, asset *model.Asset) (*model.Asset, error) {
	if err := db.Clauses(clause.Returning{}).Create(&asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}
