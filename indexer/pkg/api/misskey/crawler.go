package misskey

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type misskeyCrawler struct {
	rss3Items []*model.Item

	rss3Notes []*model.ItemId
}

func NewmisskeyCrawler() crawler.Crawler {
	return &misskeyCrawler{
		rss3Items: []*model.Item{},
		rss3Notes: []*model.ItemId{},
	}
}

func (mc *misskeyCrawler) Work(userAddress string, networkId constants.NetworkID) error {
	count := 100
	tsp := time.Now()

	noteList, err := GetUserNoteList(userAddress, count, tsp)

	if err != nil {
		logger.Errorf("%v : unable to retrieve misskey note list for %s", err, userAddress)

		return err
	}

	for _, note := range noteList {
		ni := model.NewItem(
			networkId,
			note.Link,
			model.Metadata{
				"network": constants.NetworkSymbolMisskey,
				"from":    note.Author,
			},
			constants.ItemTagsMisskeyNote,
			[]string{note.Author},
			"",
			note.Summary,
			note.Attachments,
			note.CreatedAt,
		)
		mc.rss3Items = append(mc.rss3Items, ni)

		mc.rss3Notes = append(mc.rss3Notes, &model.ItemId{
			NetworkId: networkId,
			Proof:     note.Link,
		})
	}

	return nil
}

func (mc *misskeyCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Notes: mc.rss3Notes,
		Items: mc.rss3Items,
	}
}
