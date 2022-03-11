package arbitrum

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type abCrawler struct {
	crawler.CrawlerResult
}

func NewArbitrumCrawler() crawler.Crawler {
	return &abCrawler{
		crawler.CrawlerResult{
			Assets: []*model.ItemId{},
			Notes:  []*model.ItemId{},
			Items:  []*model.Item{},
		},
	}
}

//nolint:funlen // disable line length check
func (ac *abCrawler) Work(userAddress string, network constants.NetworkID) error {
	nftTransfers, err := GetNFTTransfers(userAddress)
	if err != nil {
		return err
	}

	assets, err := GetNFTs(userAddress)
	if err != nil {
		return err
	}

	networkId := constants.NetworkSymbolArbitrum.GetID()

	// parse notes
	for _, v := range nftTransfers {
		ac.Notes = append(ac.Notes, &model.ItemId{
			NetworkId: networkId,
			Proof:     v.Hash,
		})
	}

	// parse assets
	for _, v := range assets {
		hasProof := false

		for _, nftTransfer := range nftTransfers {
			if nftTransfer.EqualsToToken(v) {
				hasProof = true

				ac.Assets = append(ac.Assets, &model.ItemId{
					NetworkId: networkId,
					Proof:     nftTransfer.Hash,
				})
			}
		}

		if !hasProof {
			logger.Error("Asset doesn't has proof")
		}
	}

	for _, v := range nftTransfers {
		tsp, err := time.Parse(time.RFC3339, v.TimeStamp)
		if err != nil {
			tsp = time.Now()
		}

		author, err := rss3uri.NewInstance("account", v.From, string(constants.NetworkSymbolArbitrum))
		if err != nil {
			// TODO
			logger.Error(err)
		}

		hasObject := false
		attachments := make([]model.Attachment, 0)

		for _, asset := range assets {
			if v.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
			}
		}

		if !hasObject {
			// TODO: get object
			logger.Errorf("Asset doesn't has the metadata.")
		}

		item := model.NewItem(
			networkId,
			v.Hash,
			model.Metadata{
				"network": constants.NetworkSymbolArbitrum,
				"from":    v.From,
				"to":      v.To,
			},
			constants.ItemTagsNFT,
			[]string{author.String()},
			"",
			"",
			attachments,
			tsp,
		)
		ac.Items = append(ac.Items, item)
	}

	return nil
}

func (ac *abCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: ac.Assets,
		Notes:  ac.Notes,
		Items:  ac.Items,
	}
}
