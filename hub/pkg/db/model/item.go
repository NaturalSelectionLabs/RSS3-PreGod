package model

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

// `item` model.
type Item struct {
	Base BaseModel `gorm:"embedded"`

	ItemID string `gorm:"primaryKey;type:uuid;column:item_id"` // uuid

	ObjectID string `gorm:"type:text;column:object_id"`

	Proof     string               `gorm:"type:text"`
	Type      constants.ItemTypeID `gorm:"type:int"`
	From      string               `gorm:"type:text"`
	To        string               `gorm:"type:text"`
	Auto      bool                 `gorm:"type:bool"`
	PageIndex int                  `gorm:"type:int"`

	PlatformCreatedAt time.Time `gorm:"index"` // create time on the platform
}
