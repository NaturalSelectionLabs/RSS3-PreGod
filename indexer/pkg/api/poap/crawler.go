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

type poapCrawler struct {
	crawler.CrawlerResult
}

func NewPoapCrawler() crawler.Crawler {
	return &poapCrawler{
		crawler.CrawlerResult{
			Items:  []*model.Item{},
			Assets: []*model.ItemId{},
			Notes:  []*model.ItemId{},
		},
	}
}

func (pc *poapCrawler) Work(param crawler.WorkParam) error {
	if param.NetworkID != constants.NetworkIDGnosisMainnet {
		return fmt.Errorf("network is not gnosis")
	}

	networkSymbol := constants.NetworkSymbolGnosisMainnet

	networkId := networkSymbol.GetID()

	poapResps, err := GetActions(param.Identity)
	if err != nil {
		logger.Error(err)

		return err
	}

	author, err := rss3uri.NewInstance("account", param.Identity, string(constants.PlatformSymbolEthereum))
	if err != nil {
		logger.Error(err)

		return err
	}

	//TODO: Since we are getting the full amount of interfaces,
	// I hope to get incremental interfaces in the future and use other methods to improve efficiency
	for _, poapResp := range poapResps {
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

		pc.Items = append(pc.Items, ni)
	}

	return nil
}

func (pc *poapCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: pc.Assets,
		Notes:  pc.Notes,
		Items:  pc.Items,
	}
}
