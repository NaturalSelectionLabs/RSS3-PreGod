package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Note{}

type Note struct {
	Identifier  string         `gorm:"column:identifier"`
	RelatedURLs pq.StringArray `gorm:"column:related_urls;type:text[]"`
	Links       string         `gorm:"column:links"`
	BackLinks   string         `gorm:"column:backlinks"`
	Tags        pq.StringArray `gorm:"column:tags;type:text[]"`
	Authors     pq.StringArray `gorm:"column:authors;type:text[]"`
	Title       string         `gorm:"column:title"`
	Summary     string         `gorm:"column:summary"`
	Attachments datatypes.JSON `gorm:"column:attachments"`
	Source      string         `gorm:"column:source"`
	Metadata    datatypes.JSON `gorm:"column:metadata"`

	common.Table
}

func (Note) TableName() string {
	return "note"
}
