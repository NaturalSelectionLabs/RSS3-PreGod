package model

// `asset` model.
type Asset struct {
	Base BaseModel `gorm:"embedded"`

	AssetID string `gorm:"primaryKey;type:uuid;column:asset_id"` // uuid
	ItemID  string `gorm:"type:uuid;column:item_id"`             // uuid
}
