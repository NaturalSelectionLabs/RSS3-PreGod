package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

// `link` model.
type Link struct {
	LinkID     string `gorm:"primaryKey;type:text;column:link_id"`
	LinkListID string `gorm:"type:text;column:link_list_id"`

	RSS3ID           string               `gorm:"type:text;column:rss3_id"`
	PrefixID         constants.PrefixID   `gorm:"type:int;column:prefix_id"`
	PlatformID       constants.PlatformID `gorm:"type:int;column:platform_id"`
	TargetRSS3ID     string               `gorm:"type:text;column:target_rss3_id"`
	TargetPrefixID   constants.PrefixID   `gorm:"type:int;column:target_prefix_id"`
	TargetPlatformID constants.PlatformID `gorm:"type:int;column:target_platform_id"`

	PageIndex int `gorm:"type:int"`

	BaseModel `gorm:"embedded"`
}
