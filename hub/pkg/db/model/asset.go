package model

// `asset` model.
type Asset struct {
	AssetID string `gorm:"primaryKey;type:uuid;column:asset_id"` // uuid
	ItemID  string `gorm:"type:uuid;column:item_id"`             // uuid

	BaseModel `gorm:"embedded"`
}
