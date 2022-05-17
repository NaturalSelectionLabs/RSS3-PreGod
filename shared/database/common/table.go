package common

import (
	"time"

	"gorm.io/gorm"
)

type Table struct {
	CreatedAt time.Time      `gorm:"autoCreateTime;not null;default:now();index"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null;default:now();index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
