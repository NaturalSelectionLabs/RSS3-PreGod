package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

// `link` model.
type Link struct {
	Base BaseModel `gorm:"embedded"`

	LinkID     string `gorm:"primaryKey;type:text;column:link_id"`
	LinkListID string `gorm:"type:text;column:link_list_id"`

	RSS3ID               string                   `gorm:"type:text;column:rss3_id"`
	PrefixID             constants.PrefixID       `gorm:"type:int;column:prefix_id"`
	PlatformNameID       constants.PlatformNameID `gorm:"type:int;column:platform_name_id"`
	TargetRSS3ID         string                   `gorm:"type:text;column:target_rss3_id"`
	TargetPrefixID       constants.PrefixID       `gorm:"type:int;column:target_prefix_id"`
	TargetPlatformNameID constants.PlatformNameID `gorm:"type:int;column:target_platform_name_id"`

	PageIndex int `gorm:"type:int"`
}
