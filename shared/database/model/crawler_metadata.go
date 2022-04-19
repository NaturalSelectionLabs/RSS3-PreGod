package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerMetadata struct {
	AccountInstance string               `gorm:"column:id;primaryKey"`
	PlatformID      constants.PlatformID `gorm:"column:platform_id;primaryKey"`
	LastBlock       int                  `gorm:"column:last_block"`

	common.Table
}

func (CrawlerMetadata) TableName() string {
	return "crawler_metadata"
}
