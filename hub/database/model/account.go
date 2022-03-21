package model

import (
	"database/sql"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/datatype"
	"github.com/lib/pq"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Account{}

type Account struct {
	ID          string               `gorm:"column:id;index:account_idx;primaryKey"`
	Platform    int                  `gorm:"column:platform;index:account_idx"`
	Name        sql.NullString       `gorm:"column:name"`
	Bio         sql.NullString       `gorm:"column:bio"`
	Avatars     pq.StringArray       `gorm:"column:avatars;type:text[]"`
	Attachments datatype.Attachments `gorm:"column:attachments;type:jsonb"`

	common.Table
}

func (a *Account) TableName() string {
	return "account"
}
