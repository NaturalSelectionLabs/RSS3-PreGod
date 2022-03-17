package misskey

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	noteList, err := GetUserNoteList(param.Identity, param.Limit, param.TimeStamp)

	if err != nil {
		logger.Errorf("%v : unable to retrieve misskey note list for %s", err, param.Identity)

		return *result, err
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
		result.Items = append(result.Items, ni)

		result.Notes = append(result.Notes, &model.ItemId{
			NetworkID: param.NetworkID,
			Proof:     note.Link,
		})
	}

	return *result, nil
}
