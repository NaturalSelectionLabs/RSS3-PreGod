package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/common"
	"github.com/google/uuid"
	"gorm.io/gorm/schema"
)

var _ schema.Tabler = &Note{}

type Note struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ItemID uuid.UUID `gorm:"type:uuid"`

	common.Table
}

func (n *Note) TableName() string {
	return "note"
}
