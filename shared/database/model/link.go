package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Link{}

type Link struct {
	Type     int    `gorm:"column:type;index:index_link_type"`
	From     string `gorm:"column:from;index:index_link_from"`
	To       string `gorm:"column:to;index:index_link_to"`
	Source   int    `gorm:"column:source;index:index_link_source"`
	Metadata string `gorm:"column:metadata"`

	common.Table
}

func (l Link) TableName() string {
	return "link"
}
