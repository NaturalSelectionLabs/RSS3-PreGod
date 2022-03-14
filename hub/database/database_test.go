package database_test

import (
	"log"

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

	if err := database.Setup(); err != nil {
		log.Fatalln(err)
	}
}
