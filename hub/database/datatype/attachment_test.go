package datatype_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/datatype"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "pregod-test-datatype-attachment.db")), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}

func TestAttachment(t *testing.T) {
	t.Parallel()

	type UserWithAttachments struct {
		ID          uint
		Name        string
		Attachments datatype.Attachments `gorm:"type:json"`
	}

	db.Migrator().DropTable(&UserWithAttachments{})

	if err := db.Migrator().AutoMigrate(&UserWithAttachments{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	attachments := datatype.Attachments{
		{Type: "image", Address: "http://example.com/logo.png"},
		{Type: "text", Content: "This is ", SizeInBytes: 100},
	}
	user := UserWithAttachments{Name: "jason", Attachments: attachments}
	db.Create(&user)

	result := UserWithAttachments{}

	if err := db.First(&result, "name = ? AND attachments = ?", "jason", attachments).Error; err != nil {
		t.Fatalf("Failed to find record with attachments, got error: %v", err)
	}

	AssertEqual(t, result.Attachments, attachments)
}
