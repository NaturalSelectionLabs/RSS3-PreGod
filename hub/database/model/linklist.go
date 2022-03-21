package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &LinkList{}

type LinkList struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	Metadata datatypes.JSONMap `gorm:"type:jsonb"`

	Type int `gorm:"column:type;not null;index"`

	Identity string `gorm:"column:identity;not null;index"`
	PrefixID int    `gorm:"column:prefix_id;not null;index"`
	SuffixID int    `gorm:"column:suffix_id;not null;index"`

	ItemCount    int `gorm:"item_count"`
	MaxPageIndex int `gorm:"max_page_index"`

	common.Table
}

func (l *LinkList) TableName() string {
	return "linklist"
}
