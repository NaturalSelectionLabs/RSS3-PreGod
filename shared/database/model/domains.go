package model

import "time"

type Domains struct {
	BlockTimestamp      time.Time `gorm:"block_timestamp"`
	TransactionHash     []byte    `gorm:"transaction_hash"`
	TransactionLogIndex int       `gorm:"transaction_log_index"`
	Type                string    `gorm:"type"`
	Name                string    `gorm:"name"`
	AddressOwner        []byte    `gorm:"address_owner"`
	AddressTarget       []byte    `gorm:"address_target"`
	ExpiredAt           time.Time `gorm:"expired_at"`
	Source              string    `gorm:"source"`
	CreatedAt           time.Time `gorm:"created_at"`
	UpdatedAt           time.Time `gorm:"updated_at"`
}

func (d Domains) TableName() string {
	return "domains"
}
