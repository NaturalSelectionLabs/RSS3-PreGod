package model

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// `account` model.
type Account struct {
	AccountID string         `gorm:"primaryKey;type:text;column:account_id"`
	Name      string         `gorm:"type:text"`
	Bio       string         `gorm:"type:text"`
	Avatars   pq.StringArray `gorm:"type:text[]"`

	Attachments datatypes.JSON `gorm:"type:jsonb"`
	// The following fields are stored in `attachments` field above:
	// Banners   pq.StringArray `gorm:"type:text[]"`
	// Websites  pq.StringArray `gorm:"type:text[]"`

	InstanceBase    InstanceBase      `gorm:"foreignkey:AccountID"` // belongs to
	AccountPlatform []AccountPlatform `gorm:"foreignkey:AccountID"` // has many

	Base BaseModel `gorm:"embedded"`
}
