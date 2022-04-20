package util

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

var keyOffset = make(map[string]int)

func GotKey(strategy string, indexer_id string, keys []string) string {
	if len(strategy) == 0 {
		strategy = "round-robin"
	}

	if len(indexer_id) == 0 {
		indexer_id = "."
	}

	var offset int

	var key string

	if strategy == "first-always" {
		key = "Bearer " + indexer_id
	} else {
		count, ok := keyOffset[indexer_id]

		if !ok {
			keyOffset[indexer_id] = 0
		}

		offset = count % len(keys)
		keyOffset[indexer_id] = count + 1
		key = keys[offset]
	}

	return key
}

func EllipsisContent(summary string, maxLength int) string {
	if maxSummaryLength := maxLength; len(summary) > maxSummaryLength { // TODO: define the max length specifically in protocol?
		summary = string([]rune(summary)[:maxSummaryLength]) + "..."
	}

	return summary
}

func GetCrawlerMetadata(identity string, platformID constants.PlatformID) (int64, error) {
	metadata, err := database.QueryCrawlerMetadata(database.DB, identity, platformID)
	if err != nil {
		return 0, fmt.Errorf("query crawler metadata error: %s", err)
	}

	if metadata == nil {
		return 0, fmt.Errorf("crawler metadata not found")
	}

	return metadata.LastBlock, nil
}

func SetCrawlerMetadata(
	instance string,
	fromHeight int64,
	platformID constants.PlatformID) error {

	if _, err := database.CreateCrawlerMetadata(database.DB, &model.CrawlerMetadata{
		AccountInstance: instance,
		PlatformID:      platformID,
		LastBlock:       fromHeight,
	}, true); err != nil {
		return fmt.Errorf("set last position error: %s", err)
	}

	return nil
}
