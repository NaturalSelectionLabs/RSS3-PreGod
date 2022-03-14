package database_test

import (
	"log"
	"testing"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func init() {
	if err := config.Setup(); err != nil {
		log.Fatalln(err)
	}

	if err := logger.Setup(); err != nil {
		log.Fatalln(err)
	}
}

func TestName(t *testing.T) {
	t.Parallel()

	db := database.GetInstance()
	t.Log(db)
}
