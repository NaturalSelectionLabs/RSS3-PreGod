package autoupdater

import (
	"context"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler_handler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
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
		tx := database.DB.Begin()
		defer tx.Rollback()

		if result.Assets != nil && len(result.Assets) > 0 {
			if _, err := database.CreateAssets(tx, result.Assets, true); err != nil {
				return err
			}
		}
		if result.Notes != nil && len(result.Notes) > 0 {
			if _, err := database.CreateNotes(tx, result.Notes, true); err != nil {
				return err
			}
		}

		return nil
	})
}
