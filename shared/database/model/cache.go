package model

import (
	"encoding/json"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
)

type Cache struct {
	Key      string          `gorm:"column:key;primaryKey"`
	Network  string          `gorm:"column:network;primaryKey"`
	Source   string          `gorm:"column:source;primaryKey"`
	BlockNum int64           `gorm:"column:block_num;index:index_block_num"`
	LogIndex int             `gorm:"column:log_index;primaryKey"`
	Data     json.RawMessage `gorm:"column:data;type:jsonb"`

	common.Table
}
