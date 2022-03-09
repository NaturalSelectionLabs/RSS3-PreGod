package model

import (
	datatype "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model/datatype"
	"github.com/lib/pq"
)

// `object` model.
type Object struct {
	ObjectID string `gorm:"primaryKey;type:text;column:object_id"` // uuid
	AlterID  string `gorm:"type:text;column:alter_id"`             // alternative readable id for item

	Title       string               `gorm:"type:text"`
	Summary     string               `gorm:"type:text"`
	Authors     pq.StringArray       `gorm:"type:text[]"`
	Tags        pq.StringArray       `gorm:"type:text[]"`
	Attachments datatype.Attachments `gorm:"type:jsonb"`

	BaseModel `gorm:"embedded"`
}
