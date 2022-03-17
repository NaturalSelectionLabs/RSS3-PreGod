package model

import (
	"database/sql"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Instance{}

type Instance struct {
	ID         string         `gorm:"column:id;index:instance_idx"`
	PrefixID   int            `gorm:"column:prefix_id"`
	Controller sql.NullString `gorm:"column:controller"`

	common.Table
}

func (i *Instance) TableName() string {
	return "instance"
}
