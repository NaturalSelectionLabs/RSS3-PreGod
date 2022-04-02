package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Account{}

type Asset struct {
	ID string `gorm:"column:id;index:index_account_id"`
	common.Table
}

func (_ Asset) TableName() string {
	return "account"
}
