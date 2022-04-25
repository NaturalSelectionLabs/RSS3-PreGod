package database

import (
	"strings"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

const (
	MaxLimit = 100
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

	internalDB, err := DB.DB()
	if err != nil {
		return err
	}

	internalDB.SetMaxOpenConns(config.Config.Postgres.MaxOpenConns)
	internalDB.SetMaxIdleConns(config.Config.Postgres.MaxIdleConns)

	internalDB.SetConnMaxIdleTime(config.Config.Postgres.ConnMaxIdleTime)
	internalDB.SetConnMaxLifetime(config.Config.Postgres.ConnMaxLifetime)

	// Install uuid extension for postgres
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		return err
	}

	if err := DB.AutoMigrate(
		// &model.Profile{},
		// &model.Account{},
		// &model.Link{},
		// &model.Asset{},
		// &model.Note{},
		// &model.CrawlerMetadata{},
		&model.ProjectInfo{},
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

func QueryLinks(db *gorm.DB, _type *int, form string, linkSources []int, profileSources []int, lastTime *time.Time, limit int) ([]model.Link, error) {
	var links []model.Link

	internalDB := db.Where(&model.Link{
		From: form,
	})

	if _type != nil {
		internalDB = internalDB.Where("type = ?", *_type)
	}

	if len(linkSources) > 0 {
		internalDB = internalDB.Where("source IN ?", linkSources)
	}

	if lastTime != nil {
		internalDB = internalDB.Where("created_at < ?", *lastTime)
	}

	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
	} else {
		limit = MaxLimit
	}

	internalDB = internalDB.Limit(limit)

	if err := internalDB.Order("created_at DESC").Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func QueryLinksByTo(db *gorm.DB, _type *int, to string, linkSources []int, lastTime *time.Time, limit int) ([]model.Link, error) {
	var links []model.Link

	internalDB := db.Where(&model.Link{
		To: to,
	})

	if _type != nil {
		internalDB = internalDB.Where("type = ?", *_type)
	}

	if len(linkSources) > 0 {
		internalDB = internalDB.Where("source IN ?", linkSources)
	}

	if lastTime != nil {
		internalDB = internalDB.Where("created_at < ?", *lastTime)
	}

	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
	} else {
		limit = MaxLimit
	}

	internalDB = internalDB.Limit(limit)

	if err := internalDB.Order("created_at DESC").Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func CreateNote(db *gorm.DB, note *model.Note, updateAll bool) (*model.Note, error) {
	note.Identifier = strings.ToLower(note.Identifier)
	note.Owner = strings.ToLower(note.Owner)

	if err := db.Model(note).Clauses(NewCreateClauses(updateAll)...).Create(note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func CreateAsset(db *gorm.DB, asset *model.Asset, updateAll bool) (*model.Asset, error) {
	asset.Identifier = strings.ToLower(asset.Identifier)
	asset.Owner = strings.ToLower(asset.Owner)

	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

func CreateNotes(db *gorm.DB, notes []model.Note, updateAll bool) ([]model.Note, error) {
	for i := range notes {
		notes[i].Identifier = strings.ToLower(notes[i].Identifier)
		notes[i].Owner = strings.ToLower(notes[i].Owner)

		if notes[i].Metadata == nil {
			notes[i].Metadata = []byte("{}")
		}
	}

	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func CreateAssets(db *gorm.DB, assets []model.Asset, updateAll bool) ([]model.Asset, error) {
	for i := range assets {
		assets[i].Identifier = strings.ToLower(assets[i].Identifier)
		assets[i].Owner = strings.ToLower(assets[i].Owner)

		if assets[i].Metadata == nil {
			assets[i].Metadata = []byte("{}")
		}
	}

	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func DeleteAsset(db *gorm.DB, asset *model.Asset) (*model.Asset, error) {
	if err := db.Clauses(clause.Returning{}).Delete(&asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

func QueryAssets(db *gorm.DB, uris []string, lastTime *time.Time, limit int) ([]model.Asset, error) {
	var assets []model.Asset

	internalDB := db.
		Where("owner IN ?", uris).
		Order("date_created DESC")

	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
	} else {
		limit = MaxLimit
	}

	internalDB = internalDB.Limit(limit)

	if lastTime != nil {
		internalDB = internalDB.Where("created_at < ?", *lastTime)
	}

	if err := internalDB.Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func QueryAllAssets(db *gorm.DB, uris []string, network constants.NetworkSymbol) ([]model.Asset, error) {
	var assets []model.Asset

	internalDB := db.
		Where("owner IN ?", uris).
		Where("metadata_network = ?", network).
		Order("date_created DESC")

	if err := internalDB.Find(&assets).Error; err != nil {
		return nil, err
	}

	return assets, nil
}

func QueryNotes(db *gorm.DB, uris []string, lastTime *time.Time, limit int) ([]model.Note, error) {
	var notes []model.Note

	internalDB := db.
		Where("owner IN ?", uris).
		Order("date_created DESC")

	if limit > 0 {
		if limit > MaxLimit {
			limit = MaxLimit
		}
	} else {
		limit = MaxLimit
	}

	internalDB = internalDB.Limit(limit)

	if lastTime != nil {
		internalDB = internalDB.Where("created_at < ?", *lastTime)
	}

	if err := internalDB.Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func CreateProfile(db *gorm.DB, profile *model.Profile, updateAll bool) (*model.Profile, error) {
	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

func CreateProfiles(db *gorm.DB, profiles []model.Profile, updateAll bool) ([]model.Profile, error) {
	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}

func CreateCrawlerMetadata(db *gorm.DB, crawler *model.CrawlerMetadata, updateAll bool) (*model.CrawlerMetadata, error) {
	if err := db.Clauses(NewCreateClauses(updateAll)...).Create(&crawler).Error; err != nil {
		return nil, err
	}

	return crawler, nil
}

func QueryCrawlerMetadata(db *gorm.DB, identity string, platformId constants.PlatformID) (*model.CrawlerMetadata, error) {
	var crawler model.CrawlerMetadata
	r := db.Where(&model.CrawlerMetadata{
		AccountInstance: identity,
		PlatformID:      platformId,
	}).Find(&crawler)

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected == 0 {
		return nil, nil
	}

	return &crawler, nil
}

func NewCreateClauses(updateAll bool) []clause.Expression {
	clauses := []clause.Expression{clause.Returning{}}

	if updateAll {
		clauses = append(clauses, clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
			UpdateAll: true,
		})
	} else {
		clauses = append(clauses, clause.OnConflict{DoNothing: true})
	}

	return clauses
}
