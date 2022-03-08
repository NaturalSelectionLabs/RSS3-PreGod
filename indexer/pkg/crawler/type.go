package crawler

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerResult struct {
	Assets []*model.ItemId
	Notes  []*model.ItemId
	Items  []*model.Item
}

type Crawler interface {
	Work(string, constants.NetworkID) error
	// GetResult return &{Assets, Notes, Items}
	GetResult() *CrawlerResult
}
