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
	ID              string         `gorm:"column:id;primaryKey"`
	Platform        int            `gorm:"column:platform;primaryKey"`
	Source          int            `gorm:"column:source;primaryKey"`
	Name            sql.NullString `gorm:"column:name"`
	Bio             sql.NullString `gorm:"column:bio"`
	Avatars         pq.StringArray `gorm:"column:avatars;type:text[]"`
	Attachments     datatypes.JSON `gorm:"column:attachments;type:jsonb;default:'{}'"`
	MetadataNetwork string         `gorm:"column:metadata_network"`
	MetadataProof   string         `gorm:"column:metadata_proof"`
	Metadata        datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'"`

	common.Table
}

func (p *Profile) TableName() string {
	return "profile"
}
