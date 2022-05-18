package model

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Note{}

type Note struct {
	Identifier      string         `gorm:"column:identifier;primaryKey"`
	ContractAddress string         `gorm:"column:contract_address;index"`
	LogIndex        int            `gorm:"column:log_index;index"`
	TokenID         string         `gorm:"column:token_id;index"`
	Owner           string         `gorm:"colum:owner;index:index_note_owner"`
	ProfileSourceID int            `gorm:"colum:profile_source_id;index:index_note_profile_source_id"`
	RelatedURLs     pq.StringArray `gorm:"column:related_urls;type:text[]"`
	Tags            pq.StringArray `gorm:"column:tags;type:text[]"`
	Authors         pq.StringArray `gorm:"column:authors;type:text[]"`
	Title           string         `gorm:"column:title"`
	Summary         string         `gorm:"column:summary"`
	Attachments     datatypes.JSON `gorm:"column:attachments;not null;default:'[]'"`
	Source          string         `gorm:"column:source;index:index_note_source"`
	MetadataNetwork string         `gorm:"column:metadata_network"`
	MetadataProof   string         `gorm:"column:metadata_proof"`
	Metadata        datatypes.JSON `gorm:"column:metadata;not null;default:'{}'"`
	DateCreated     time.Time      `gorm:"column:date_created;index:index_note_date_created"`
	DateUpdated     time.Time      `gorm:"column:date_updated;index:index_note_date_updated"`

	common.Table
}

func (Note) TableName() string {
	return "note"
}

func (n Note) String() string {
	return n.Identifier + " " + n.Metadata.String()
}
