package crawlers

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerResult struct {
	Assets  []*model.ItemId
	Notes   []*model.ItemId
	Items   []*model.Item
	Objects []*model.Object
}

type Crawler interface {
	Work(string, constants.NetworkName) error
	// GetResult return &{Assets, Notes, Items, Objects}
	GetResult() *CrawlerResult
}
