package autoupdater

import (
	"context"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/vmihailenco/msgpack"
)

func AddToRecentVisitQueue(ctx context.Context, param *crawler.WorkParam) error {
	item, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}

	return RecentVisitQueue.Add(ctx, item)
}

// TODO: Can be optimized and merged with item's Execute code
func RunRecentVisitQueue(ctx context.Context) error {
	return RecentVisitQueue.Iter(ctx, func(s string) error {
		var c crawler.Crawler
		param := crawler.WorkParam{}
		if err := msgpack.Unmarshal([]byte(s), &param); err != nil {
			return err
		}

		// choose executor
		c = crawler_handler.MakeCrawlers(param.NetworkID)
		if c == nil {
			return fmt.Errorf("unknown network id")
		}

		// Work and save result
		if err := c.Work(param); err != nil {
			return err
		}
		result := c.GetResult()
		if result.Items != nil {
			for _, item := range result.Items {
				if err := db.InsertItem(item).Err(); err != nil {
					return err
				}
			}
		}
		instance := rss3uri.NewAccountInstance(param.Identity, param.PlatformID.Symbol())
		if result.Assets != nil {
			db.SetAssets(instance, result.Assets, param.NetworkID)
		}
		if result.Notes != nil {
			db.AppendNotes(instance, result.Notes)
		}

		return nil
	})
}
