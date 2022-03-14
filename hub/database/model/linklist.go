package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &LinkList{}

type LinkList struct {
	ID       uuid.UUID         `gorm:"type:uuid;primaryKey"`
	Metadata datatypes.JSONMap `gorm:"type:jsonb"`

	common.Table
}

func (l *LinkList) TableName() string {
	return "linklist"
}
