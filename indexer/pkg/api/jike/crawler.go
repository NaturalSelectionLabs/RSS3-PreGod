package jike

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type jikeCrawler struct {
	rss3Items []*model.Item

	rss3Notes []*model.ItemId
}

func NewJikeCrawler() crawler.Crawler {
	return &jikeCrawler{
		rss3Items: []*model.Item{},
		rss3Notes: []*model.ItemId{},
	}
}

func (mc *jikeCrawler) Work(param crawler.WorkParam) error {
	timeline, err := GetUserTimeline(param.Identity)

	if err != nil {
		return err
	}

	for _, item := range timeline {
		ni := model.NewItem(
			param.NetworkID,
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
			NetworkID: param.NetworkID,
			Proof:     item.Link,
		})
	}

	return nil
}

func (mc *jikeCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Notes: mc.rss3Notes,
		Items: mc.rss3Items,
	}
}
