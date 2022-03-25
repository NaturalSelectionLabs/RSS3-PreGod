package autoupdater

import (
	"context"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/jike"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/twitter"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/vmihailenco/msgpack"
)

func AddToRecentVisitQueue(ctx context.Context, param *crawler.WorkParam) error {
	item, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}

	return RecentVisitQueue.Add(ctx, string(item))
}

func RunRecentVisitQueue(ctx context.Context) error {
	return RecentVisitQueue.Iter(ctx, func(s string) error {
		var c crawler.Crawler
		param := crawler.WorkParam{}
		if err := msgpack.Unmarshal([]byte(s), &param); err != nil {
			return err
		}

		// choose executor
		switch param.NetworkID {
		case constants.NetworkIDEthereumMainnet,
			constants.NetworkIDBNBChain,
			constants.NetworkIDAvalanche,
			constants.NetworkIDFantom,
			constants.NetworkIDPolygon:
			c = moralis.NewMoralisCrawler()
		case constants.NetworkIDMisskey:
			c = misskey.NewMisskeyCrawler()
		case constants.NetworkIDJike:
			c = jike.NewJikeCrawler()
		case constants.NetworkIDTwitter:
			c = twitter.NewTwitterCrawler()
		default:
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
