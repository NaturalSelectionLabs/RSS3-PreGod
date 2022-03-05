package db

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/dblogger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Setup() error {
	var err error

	// Use custom logger
	logger := dblogger.New()
	logger.SetAsDefault()

	// Establish a connection to the database
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: config.Config.Postgres.DSN,
	}), &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		NowFunc:                                  func() time.Time { return time.Now().UTC() },
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger,
	})

	if err != nil {
		return err
	}

	// Ping
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Set config
	sqlDB.SetMaxOpenConns(config.Config.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Config.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(config.Config.Postgres.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(config.Config.Postgres.ConnMaxLifetime)

	// defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.InstanceBase{},

		&model.Account{},
		&model.AccountPlatform{},

		&model.Object{},
		&model.Item{},
		&model.Asset{},
		&model.Note{},

		&model.Link{},
		&model.LinkList{},

		&model.ThirdPartyStorage{},

		&model.Signature{},
	)
}
