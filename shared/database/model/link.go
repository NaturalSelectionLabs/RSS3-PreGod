package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Link{}

type Link struct {
	Type             int    `gorm:"column:type;index:index_link_type"`
	From             string `gorm:"column:from;index:index_link_from"`
	FromInstanceType int    `gorm:"column:from_instance_type"`
	FromPlatformID   int    `gorm:"column:from_platform_id"`
	To               string `gorm:"column:to;index:index_link_to"`
	ToInstanceType   int    `gorm:"column:to_instance_type"`
	ToPlatformID     int    `gorm:"column:to_platform_id"`
	Source           int    `gorm:"column:source"`
	Metadata         string `gorm:"column:metadata;default:{}"`

	common.Table
}

func (l Link) TableName() string {
	return "link"
}
