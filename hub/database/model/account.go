package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Account{}

type Account struct {
	ID              string `gorm:"column:id;"`
	Platform        int    `gorm:"column:platform"`
	ProfileID       string `gorm:"column:profile_id"`
	ProfilePlatform int    `gorm:"column:profile_platform"`
	Source          int    `gorm:"column:source"`

	common.Table
}

func (a Account) TableName() string {
	return "account"
}
