package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
)

func init() {
	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}
}

func main() {
	for {
		// get data from db
		// change db
		//save in db
	}

}
