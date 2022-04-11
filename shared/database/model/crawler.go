package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerMetadata struct {
	AccountInstance string              `gorm:"column:id,primaryKey"`
	NetworkId       constants.NetworkID `gorm:"column:network_id.primaryKey"`
	LastBlock       int                 `gorm:"column:last_block"`

	common.Table
}

func (CrawlerMetadata) TableName() string {
	return "crawler_metadata"
}
