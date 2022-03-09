package jike

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type jikeCrawler struct {
	rss3Items []*model.Item

	rss3Assets, rss3Notes []*model.ItemId
}

func NewJikeCrawler() crawler.Crawler {
	return &jikeCrawler{
		rss3Items: []*model.Item{},
		rss3Notes: []*model.ItemId{},
	}
}

func (mc *jikeCrawler) Work(userAddress string, networkId constants.NetworkID) error {
	timeline, err := GetUserTimeline(userAddress)

	if err != nil {
		return err
	}

	for _, item := range timeline {
		ni := model.NewItem(
			networkId,
			item.Link,
			model.Metadata{
				"network": constants.NetworkSymbolJike,
				"from":    item.Author,
			},
			constants.ItemTagsJikePost,
			[]string{item.Author},
			"",
			item.Summary,
			item.Attachments,
			item.Timestamp,
		)
		mc.rss3Items = append(mc.rss3Items, ni)

		mc.rss3Notes = append(mc.rss3Notes, &model.ItemId{
			NetworkId: networkId,
			Proof:     item.Link,
		})
	}

	return nil
}

func (mc *jikeCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: mc.rss3Assets,
		Notes:  mc.rss3Notes,
		Items:  mc.rss3Items,
	}
}
