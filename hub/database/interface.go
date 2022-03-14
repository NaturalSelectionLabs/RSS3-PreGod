package database

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	instance Database
	locker   sync.Mutex
)

type Database interface {
	DB(ctx context.Context) *gorm.DB
	Tx(ctx context.Context) *gorm.DB
}

func GetInstance() Database {
	// nolint:nestif // Must use nesting of if statements to ensure it is a single
	if instance == nil {
		locker.Lock()
		defer locker.Unlock()

		if instance == nil {
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
				log.Fatalln(err)
			}

			// Install uuid extension for postgres
			if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
				log.Fatalln(err)
			}

			if err := db.AutoMigrate(
				&model.Account{},
				&model.AccountPlatform{},
				&model.Instance{},
				&model.LinkList{},
				&model.Link{},
				&model.Signature{},
				&model.Asset{},
				&model.Note{},
			); err != nil {
				log.Fatalln(err)
			}

			instance = &database{
				db: db,
			}
		}
	}

	return instance
}
