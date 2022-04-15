package model

import (
	"database/sql"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Profile{}

type Profile struct {
	ID              string         `gorm:"column:id;index;primaryKey"`
	Platform        int            `gorm:"column:platform;primaryKey"`
	Source          int            `gorm:"column:source;primaryKey"`
	Name            sql.NullString `gorm:"column:name"`
	Bio             sql.NullString `gorm:"column:bio"`
	Avatars         pq.StringArray `gorm:"column:avatars;type:text[]"`
	Attachments     datatypes.JSON `gorm:"column:attachments;type:jsonb"`
	MetadataNetwork string         `gorm:"column:metadata_network"`
	MetadataProof   string         `gorm:"column:metadata_proof"`
	Metadata        string         `gorm:"metadata"`

	common.Table
}

func (p *Profile) TableName() string {
	return "profile"
}
