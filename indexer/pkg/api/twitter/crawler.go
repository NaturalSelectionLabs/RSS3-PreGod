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

const DefaultTwitterCount = 200

func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	if param.NetworkID != constants.NetworkIDTwitter {
		return *result, fmt.Errorf("network is not twitter")
	}

	networkSymbol := constants.NetworkSymbolTwitter

	networkId := networkSymbol.GetID()

	contentInfos, err := GetTimeline(param.Identity, DefaultTwitterCount)
	if err != nil {
		return *result, err
	}

	author, err := rss3uri.NewInstance("account", param.Identity, string(constants.PlatformSymbolTwitter))
	if err != nil {
		return *result, err
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

		result.Items = append(result.Items, ni)
	}

	return *result, nil
}
