package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Signature{}

type Signature struct {
	FileURI   string         `gorm:"column:file_uri;type:text;primaryKey;"`
	Signature string         `gorm:"column:signature;type:text"`
	Agents    datatypes.JSON `gorm:"column:agents;type:jsonb"`

	common.Table
}

func (s *Signature) TableName() string {
	return "signature"
}
