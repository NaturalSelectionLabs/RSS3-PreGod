package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

// `account_platform` model.
type AccountPlatform struct {
	AccountID         string               `gorm:"primaryKey;type:text;column:account_id"`
	PlatformID        constants.PlatformID `gorm:"type:int;column:platform_id"`
	PlatformAccountID string               `gorm:"type:text;column:platform_account_id"` // account ID on the platform

	BaseModel `gorm:"embedded"`
}
