package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
	"time"
)

var _ schema.Tabler = &Note{}

type Note struct {
	Identifier      string         `gorm:"column:identifier;primaryKey"`
	Owner           string         `gorm:"colum:owner"`
	RelatedURLs     pq.StringArray `gorm:"column:related_urls;type:text[]"`
	Tags            pq.StringArray `gorm:"column:tags;type:text[]"`
	Authors         pq.StringArray `gorm:"column:authors;type:text[]"`
	Title           string         `gorm:"column:title"`
	Summary         string         `gorm:"column:summary"`
	Attachments     datatypes.JSON `gorm:"column:attachments"`
	Source          string         `gorm:"column:source"`
	MetadataNetwork string         `gorm:"column:metadata_network"`
	MetadataProof   string         `gorm:"column:metadata_proof"`
	Metadata        datatypes.JSON `gorm:"column:metadata"`
	DateCreated     time.Time      `gorm:"column:date_created"`
	DateUpdated     time.Time      `gorm:"column:date_updated"`

	common.Table
}

func (Note) TableName() string {
	return "note"
}
