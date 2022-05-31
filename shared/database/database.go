package database

import (
	"encoding/json"
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
	Trigger  = `
-- CREATE FUNCTION
CREATE OR REPLACE FUNCTION serial_transaction_log_index() RETURNS TRIGGER AS
$$
DECLARE
    _transaction_log_index int;
BEGIN
    _transaction_log_index := (SELECT COALESCE(MAX(transaction_log_index), -1)
                               FROM note4
                               WHERE transaction_hash = NEW.transaction_hash);
    IF NEW.transaction_log_index = -1 THEN
        NEW.transaction_log_index = _transaction_log_index + 1;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

-- CREATE TRIGGER
BEGIN;
DROP TRIGGER IF EXISTS trigger_transaction_log_index ON note4;

CREATE TRIGGER trigger_transaction_log_index
    BEFORE INSERT
    ON note4
    FOR EACH ROW
EXECUTE FUNCTION serial_transaction_log_index();
END;
`
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
		CreateBatchSize:                          1000,
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

	//if err := DB.AutoMigrate(
	//	&model.Profile{},
	//	&model.Account{},
	//	&model.Asset{},
	//	&model.Note{},
	//	&model.CrawlerMetadata{},
	//	&model.Cache{},
	//); err != nil {
	//	return err
	//}

	// if err := DB.Exec("CREATE INDEX IF NOT EXISTS index_note_owner_and_date_created ON note (owner, date_created);").Error; err != nil {
	// 	return err
	// }

	// if err := DB.Exec(Trigger).Error; err != nil {
	// 	return err
	// }

	return nil
}

func CreateNote(db *gorm.DB, note *model.Note, updateAll bool) (*model.Note, error) {
	note.Identifier = strings.ToLower(note.Identifier)
	note.Owner = strings.ToLower(note.Owner)

	if err := db.Model(note).Clauses(NewCreateClauses(updateAll, true, true)...).Create(note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

func CreateAsset(db *gorm.DB, asset *model.Asset, updateAll bool) (*model.Asset, error) {
	asset.Identifier = strings.ToLower(asset.Identifier)
	asset.Owner = strings.ToLower(asset.Owner)

	if err := db.Clauses(NewCreateClauses(updateAll, true, true)...).Create(asset).Error; err != nil {
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

	if err := db.Clauses(NewCreateClauses(updateAll, true, true)...).Create(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func CreateNotesDoNothing(db *gorm.DB, notes []model.Note) ([]model.Note, error) {
	for i := range notes {
		notes[i].Identifier = strings.ToLower(notes[i].Identifier)
		notes[i].Owner = strings.ToLower(notes[i].Owner)

		if notes[i].Metadata == nil {
			notes[i].Metadata = []byte("{}")
		}
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&notes).Error; err != nil {
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

	if err := db.Clauses(NewCreateClauses(updateAll, true, true)...).Create(&assets).Error; err != nil {
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

func CreateCrawlerMetadata(db *gorm.DB, crawler *model.CrawlerMetadata, updateAll bool) (*model.CrawlerMetadata, error) {
	if err := db.Clauses(NewCreateClauses(updateAll, false, false)...).Create(&crawler).Error; err != nil {
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

func QueryCache(db *gorm.DB, key, network, source string) (json.RawMessage, error) {
	cache := model.Cache{}

	if err := db.
		Model((*model.Cache)(nil)).
		Where(map[string]interface{}{
			"key":     key,
			"network": network,
			"source":  source,
		}).
		First(&cache).
		Error; err != nil {
		return nil, err
	}

	return cache.Data, nil
}

func CreateCache(db *gorm.DB, key, network, source string, data json.RawMessage) error {
	return db.
		Model((*model.Cache)(nil)).
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&model.Cache{
			Key:     key,
			Network: network,
			Source:  source,
			Data:    data,
		}).
		Error
}

func NewCreateClauses(updateAll bool, updateMetadata bool, updateAttachments bool) []clause.Expression {
	clauses := []clause.Expression{
		// clause.Returning{}
	}

	assignmentArrary := []string{}

	if updateAll {
		assignmentArrary = append(assignmentArrary, "updated_at")
	}

	if updateMetadata {
		assignmentArrary = append(assignmentArrary, "metadata")
	}

	if updateAttachments {
		assignmentArrary = append(assignmentArrary, "attachments")
	}

	if len(assignmentArrary) > 0 {
		clauses = append(clauses, clause.OnConflict{
			DoUpdates: clause.AssignmentColumns(assignmentArrary),
			UpdateAll: true,
		})
	} else {
		clauses = append(clauses, clause.OnConflict{DoNothing: true})
	}

	return clauses
}
