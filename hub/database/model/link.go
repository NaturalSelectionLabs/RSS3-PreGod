package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Link{}

type Link struct {
	Type     int    `gorm:"column:type"`
	From     string `gorm:"column:from"`
	To       string `gorm:"column:to"`
	Source   int    `gorm:"column:source"`
	Metadata string `gorm:"column:metadata"`

	common.Table
}

func (l Link) TableName() string {
	return "link"
}
