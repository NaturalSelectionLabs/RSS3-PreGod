package model

// `note` model.
type Note struct {
	Base BaseModel `gorm:"embedded"`

	NoteID string `gorm:"primaryKey;type:uuid;column:note_id"` // uuid
	ItemID string `gorm:"type:uuid;column:item_id"`            // uuid
}
