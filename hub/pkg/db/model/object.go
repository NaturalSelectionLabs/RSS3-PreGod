package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// `object` model.
type Object struct {
	BaseModel `gorm:"embedded"`

	ObjectID string `gorm:"primaryKey;type:text;column:object_id"` // uuid
	AlterID  string `gorm:"type:text;column:alter_id"`             // alternative readable id for item

	Title       string         `gorm:"type:text"`
	Summary     string         `gorm:"type:text"`
	Authors     pq.StringArray `gorm:"type:text[]"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	Attachments datatypes.JSON `gorm:"type:jsonb"`

	PlatformNameID   constants.PlatformNameID `gorm:"type:int;column:platform_name_id"`
	PlatformObjectID string                   `gorm:"type:text;column:platform_object_id"` // object ID on the platform
}
