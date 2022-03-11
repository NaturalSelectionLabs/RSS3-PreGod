package twitter

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type twitterCrawler struct {
	crawler.CrawlerResult
}

func NewTwitterCrawler() crawler.Crawler {
	return &twitterCrawler{
		crawler.CrawlerResult{
			Items: []*model.Item{},
			Notes: []*model.ItemId{},
		},
	}
}

const DefaultTwitterCount = 200

func (tc *twitterCrawler) Work(userAddress string, itemType constants.NetworkID) error {
	if itemType != constants.NetworkIDTwitter {
		return fmt.Errorf("network is not twitter")
	}

	networkSymbol := constants.NetworkSymbolTwitter

	networkId := networkSymbol.GetID()

	contentInfos, err := GetTimeline(userAddress, DefaultTwitterCount)
	if err != nil {
		return err
	}

	author, err := rss3uri.NewInstance("account", userAddress, string(constants.PlatformSymbolTwitter))
	if err != nil {
		return err
	}

	for _, contentInfo := range contentInfos {
		tsp, err := contentInfo.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		ni := model.NewItem(
			networkId,
			"",
			model.Metadata{},
			constants.ItemTagsTweet,
			[]string{author.String()},
			"",
			contentInfo.PreContent,
			[]model.Attachment{},
			tsp,
		)
		logger.Debugf("twitter item: %v", ni)

		tc.Items = append(tc.Items, ni)
	}

	return nil
}

func (pc *twitterCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: pc.Assets,
		Notes:  pc.Notes,
	}
}
