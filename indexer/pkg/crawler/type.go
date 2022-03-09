package crawler

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type CrawlerResult struct {
	Assets []*model.ItemId
	Notes  []*model.ItemId
	Items  []*model.Item
}

type WorkParam struct {
	UserAddress string
	NetworkId   constants.NetworkID
	Limit       int // aka Count, limit the number of items to be crawled
	Timestamp   time.Time
}

type Crawler interface {
	Work(WorkParam) error
	// GetResult return &{Assets, Notes, Items}
	GetResult() *CrawlerResult
}
