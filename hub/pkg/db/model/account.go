package model

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

// `account` model.
type Account struct {
	AccountID string         `gorm:"primaryKey;type:text;column:account_id" json:"account_id"`
	Name      string         `gorm:"type:text" json:"name"`
	Bio       string         `gorm:"type:text" json:"bio"`
	Avatars   pq.StringArray `gorm:"type:text[]" json:"avatars"`

	Attachments datatypes.JSON `gorm:"type:jsonb" json:"attachments,omitempty"`
	// The following fields are stored in `attachments` field above:
	// Banners   pq.StringArray `gorm:"type:text[]"`
	// Websites  pq.StringArray `gorm:"type:text[]"`

	InstanceBase    *InstanceBase     `gorm:"foreignkey:AccountID" json:"instance_base,omitempty"`    // belongs to
	AccountPlatform []AccountPlatform `gorm:"foreignkey:AccountID" json:"account_platform,omitempty"` // has many

	BaseModel `gorm:"embedded"`
}
