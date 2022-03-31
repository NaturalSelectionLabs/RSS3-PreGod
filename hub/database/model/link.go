package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Link{}

type Link struct {
	Type     int    `gorm:"column:type;index:index_link"`
	From     string `gorm:"column:from;index:index_link"`
	To       string `gorm:"column:to;index:index_link"`
	Source   int    `gorm:"column:source;index:index_link"`
	Metadata string `gorm:"column:metadata;index:index_link"`

	common.Table
}

func (l Link) TableName() string {
	return "link"
}
