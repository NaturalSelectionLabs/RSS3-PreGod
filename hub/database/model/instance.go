package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Instance{}

type Instance struct {
	ID         string `gorm:"column:id;index:instance_idx;primaryKey"`
	Platform   int    `gorm:"column:platform"`
	Controller string `gorm:"column:controller"`

	common.Table
}

func (i *Instance) TableName() string {
	return "instance"
}
