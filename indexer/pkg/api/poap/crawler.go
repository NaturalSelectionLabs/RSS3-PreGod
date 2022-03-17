package poap

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	if param.NetworkID != constants.NetworkIDGnosisMainnet {
		return *result, fmt.Errorf("network is not gnosis")
	}

	networkSymbol := constants.NetworkSymbolGnosisMainnet

	networkId := networkSymbol.GetID()

	response, err := GetActions(param.Identity)
	if err != nil {
		logger.Error(err)

		return *result, err
	}

	author, err := rss3uri.NewInstance("account", param.Identity, string(constants.PlatformSymbolEthereum))
	if err != nil {
		logger.Error(err)

		return *result, err
	}

	//TODO: Since we are getting the full amount of interfaces,
	// I hope to get incremental interfaces in the future and use other methods to improve efficiency
	for _, poapResp := range response {
		tsp, err := poapResp.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		ni := model.NewItem(
			networkId,
			"",
			model.Metadata{
				"from": "0x0",
				"to":   poapResp.Owner,
			},
			constants.ItemTagsNFTPOAP,
			[]string{author.String()},
			"",
			"",
			[]model.Attachment{},
			tsp,
		)

		result.Items = append(result.Items, ni)
	}

	return *result, nil
}
