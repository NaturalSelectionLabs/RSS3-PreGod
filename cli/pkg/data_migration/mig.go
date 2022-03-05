package data_migration

import (
	"errors"
	"log"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/pkg/data_migration/handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/pkg/data_migration/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// what file:
// 1. main index file
// 2. link list (following)
// 3. link backlist (following)
// 4. auto asset list
func migrate(mongouri string) error {
	// set up mongodb
	err := mongo.Setup(mongouri)
	if err != nil {
		return err
	}

	var result []bson.D
	if err = mongo.GetAllData(&result); err != nil {
		return err
	}

	log.Println("got MongoDB data:", len(result))

	var warnings []string

	for i, item := range result {
		row := item.Map()

		// do something
		path, ok := row["path"].(string)
		if !ok {
			return errors.New("path is not string")
		}

		log.Println("item:", i, "path:", path)

		content, ok := row["content"].(bson.D)
		if !ok {
			return errors.New("content is not bson.D")
		}

		if err != nil {
			return err
		}

		var errInHandle error

		if !strings.Contains(path, "-") {
			// is main index
			errInHandle = handler.MainIndex(content)
		} else if strings.Contains(path, "-list-assets") {
			errInHandle = handler.AssetListAuto(content)
		} else if strings.Contains(path, "-list-links") {
			errInHandle = handler.LinkList(content)
		} else if strings.Contains(path, "-list-backlinks") {
			errInHandle = handler.LinkBackList(content)
		} else {
			warnings = append(warnings, "unknown path: "+path)
			// errInHandle = errors.New("Unknown file: " + path)
		}

		if errInHandle != nil {
			log.Println("Error: ", path)

			return errInHandle
		} else {
			log.Println("Success: ", path)
		}
	}

	if len(warnings) > 0 {
		log.Println("Warnings: ", warnings)
	}

	return nil
}
