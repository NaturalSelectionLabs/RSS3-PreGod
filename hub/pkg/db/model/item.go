package model

import (
	"time"
)

// `item` model.
type Item struct {
	ItemID string `gorm:"primaryKey;type:uuid;column:item_id"` // uuid

	ObjectID string `gorm:"type:text;column:object_id"`

	Auto      bool `gorm:"type:bool"`
	PageIndex int  `gorm:"type:int"`

	OriginalCreatedAt time.Time `gorm:"index"` // create time on the platform
	OriginalUpdatedAt time.Time `gorm:"index"` // update time on the platform

	BaseModel `gorm:"embedded"`
}
