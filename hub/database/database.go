package database

import (
	"context"

	"gorm.io/gorm"
)

var (
	_ Database = &database{}
)

type database struct {
	db *gorm.DB
}

func (d *database) DB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *database) Tx(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx).Begin()
}
