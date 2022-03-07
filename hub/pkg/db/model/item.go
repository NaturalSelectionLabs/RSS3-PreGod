package model

import (
	"time"
)

// `item` model.
type Item struct {
	ItemID string `gorm:"primaryKey;type:uuid;column:item_id"` // uuid

	ObjectID string `gorm:"type:text;column:object_id"`

	Proof     string `gorm:"type:text"`
	From      string `gorm:"type:text"`
	To        string `gorm:"type:text"`
	Auto      bool   `gorm:"type:bool"`
	PageIndex int    `gorm:"type:int"`

	PlatformCreatedAt time.Time `gorm:"index"` // create time on the platform

	Base BaseModel `gorm:"embedded"`
}
