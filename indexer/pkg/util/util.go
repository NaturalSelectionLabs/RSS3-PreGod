package util

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
