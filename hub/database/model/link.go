package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/google/uuid"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Link{}

type Link struct {
	Type int `gorm:"column:type;not null;index"`

	Identity string `gorm:"column:identity;not null;index"`
	PrefixID int    `gorm:"column:prefix_id;not null;index"`
	SuffixID int    `gorm:"column:suffix_id;not null;index"`

	TargetIdentity string `gorm:"column:target_identity;not null;index"`
	TargetPrefixID int    `gorm:"column:target_prefix_id;not null;index"`
	TargetSuffixID int    `gorm:"column:target_suffix_id;not null;index"`

	PageIndex int `gorm:"column:page_index;index"`

	LinkListID uuid.UUID `gorm:"column:link_list_id;type:uuid"`
	LinkList   LinkList  `gorm:"foreignKey:LinkListID"`

	common.Table
}

func (l *Link) TableName() string {
	return "link"
}
