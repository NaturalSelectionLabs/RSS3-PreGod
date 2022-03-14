package misskey

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type misskeyCrawler struct {
	rss3Items []*model.Item

	rss3Notes []*model.ItemId
}

func NewMisskeyCrawler() crawler.Crawler {
	return &misskeyCrawler{
		rss3Items: []*model.Item{},
		rss3Notes: []*model.ItemId{},
	}
}

func (mc *misskeyCrawler) Work(param crawler.WorkParam) error {
	noteList, err := GetUserNoteList(param.Identity, param.Limit, param.TimeStamp)

	if err != nil {
		logger.Errorf("%v : unable to retrieve misskey note list for %s", err, param.Identity)

		return err
	}

	for _, note := range noteList {
		ni := model.NewItem(
			param.NetworkID,
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
			NetworkID: param.NetworkID,
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
