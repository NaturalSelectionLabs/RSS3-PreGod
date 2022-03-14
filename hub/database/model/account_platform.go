package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &AccountPlatform{}

type AccountPlatform struct {
	AccountID         string `gorm:"type:text;column:account_id;index:account_platform_idx"`
	AccountPlatformID int    `gorm:"column:account_platform_id;index:account_platform_idx"`
	PlatformAccountID string `gorm:"type:text;column:platform_account_id;index:account_platform_target_idx"`
	PlatformID        int    `gorm:"type:int;column:platform_id;index:account_platform_target_idx"`

	common.Table
}

func (a *AccountPlatform) TableName() string {
	return "account_platform"
}
