package jike

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	timeline, err := GetUserTimeline(param.Identity)

	if err != nil {
		return *result, err
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
		result.Items = append(result.Items, ni)

		result.Notes = append(result.Notes, &model.ItemId{
			NetworkID: param.NetworkID,
			Proof:     item.Link,
		})
	}

	return *result, nil
}
