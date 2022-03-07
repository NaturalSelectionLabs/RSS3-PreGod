package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"autoCreateTime;not null;default:now()"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
