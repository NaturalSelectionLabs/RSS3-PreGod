package model

import (
	"database/sql"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/lib/pq"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Profile{}

type Profile struct {
	ID          string               `gorm:"column:id;index:index_profile_id"`
	Platform    int                  `gorm:"column:platform;index_profile_platform"`
	Source      int                  `gorm:"column:source;index_profile_source"`
	Name        sql.NullString       `gorm:"column:name"`
	Bio         sql.NullString       `gorm:"column:bio"`
	Avatars     pq.StringArray       `gorm:"column:avatars;type:text[]"`
	Attachments datatype.Attachments `gorm:"column:attachments;type:jsonb"`
	Metadata    string               `gorm:"metadata"`

	common.Table
}

func (p *Profile) TableName() string {
	return "profile"
}
