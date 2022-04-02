package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Account{}

type Note struct {
	ID string `gorm:"column:id;index:index_account_id"`
	// Platform        int    `gorm:"column:platform;index:index_account_platform"`
	// ProfileID       string `gorm:"column:profile_id;index:index_account_profile_id"`
	// ProfilePlatform int    `gorm:"column:profile_platform;index:index_profile_platform"`
	// Source          int    `gorm:"column:source;index:index_account_source"`

	common.Table
}

func (_ Note) TableName() string {
	return "account"
}
