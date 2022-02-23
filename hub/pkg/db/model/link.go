package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

// `link` model.
type Link struct {
	Base BaseModel `gorm:"embedded"`

	LinkID string `gorm:"primaryKey;type:text;column:asset_id"`
	ItemID string `gorm:"type:text;column:item_id"`

	RSS3ID       string           `gorm:"type:text;column:rss3_id"`
	Prefix       constants.Prefix `gorm:"type:text;column:prefix"`
	TargetRSS3ID string           `gorm:"type:text;column:target_rss3_id"`
	TargetPrefix constants.Prefix `gorm:"type:text;column:target_prefix"`

	PageIndex int `gorm:"type:int"`
}
