package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"autoCreateTime;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
