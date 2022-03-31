package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Account{}

type Account struct {
	ID              string `gorm:"column:id;index:index_account"`
	Platform        int    `gorm:"column:platform_;index:index_account"`
	ProfileID       string `gorm:"column:profile_id;index:index_account"`
	ProfilePlatform int    `gorm:"column:profile_platform;index:index_account"`
	Source          int    `gorm:"column:source;index:index_account"`

	common.Table
}

func (a Account) TableName() string {
	return "account"
}
