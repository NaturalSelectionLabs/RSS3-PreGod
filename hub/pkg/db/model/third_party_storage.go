package model

import (
	"gorm.io/datatypes"
)

// `third_party_storage` model.
type ThirdPartyStorage struct {
	StorageID string `gorm:"primaryKey;type:text;column:storage_id"`
	RSS3ID    string `gorm:"type:text;column:rss3_id"`

	AppName     string         `gorm:"type:text"`
	Key         string         `gorm:"type:text"`
	Value       datatypes.JSON `gorm:"type:jsonb"`
	Description string         `gorm:"type:text"`

	BaseModel `gorm:"embedded"`
}
