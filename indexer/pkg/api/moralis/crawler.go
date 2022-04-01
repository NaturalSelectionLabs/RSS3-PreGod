package moralis

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type moralisCrawler struct {
	crawler.DefaultCrawler
}

func NewMoralisCrawler() crawler.Crawler {
	return &moralisCrawler{
		crawler.DefaultCrawler{
			Items:  []*model.Item{},
			Assets: []*model.ObjectId{},
			Notes:  []*model.ObjectId{},
		},
	}
}

//nolint:funlen // disable line length check
func (c *moralisCrawler) Work(param crawler.WorkParam) error {
	chainType := GetChainType(param.NetworkID)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()
	networkId := networkSymbol.GetID()
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, config.Config.Indexer.Moralis.ApiKey)

	if err != nil {
		return err
	}

	//TODO: tsp
	assets, err := GetNFTs(param.Identity, chainType, config.Config.Indexer.Moralis.ApiKey)
	if err != nil {
		return err
	}
	//parser
	for _, nftTransfer := range nftTransfers.Result {
		c.Notes = append(c.Notes, &model.ObjectId{
			NetworkID: networkId,
			Proof:     nftTransfer.TransactionHash,
		})
	}

	for _, asset := range assets.Result {
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true

				c.Assets = append(c.Assets, &model.ObjectId{
					NetworkID: networkId,
					Proof:     nftTransfer.TransactionHash,
				})
			}
		}

		if !hasProof {
			// TODO: error handle here
			logger.Errorf("Asset doesn't has proof.")
		}
	}

	// make the item list complete
	for _, nftTransfer := range nftTransfers.Result {
		// TODO: make attachments
		tsp, err := nftTransfer.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolEthereum)

		hasObject := false

		for _, asset := range assets.Result {
			if nftTransfer.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
			}
		}

		if !hasObject {
			// TODO: get object
			logger.Errorf("Asset doesn't has the metadata.")
		}

		ni := model.NewItem(
			networkId,
			nftTransfer.TransactionHash,
			model.Metadata{
				"from": nftTransfer.FromAddress,
				"to":   nftTransfer.ToAddress,
			},
			constants.ItemTagsNFT,
			[]string{author.String()},
			"",
			"",
			[]model.Attachment{},
			tsp,
		)
		c.Items = append(c.Items, ni)
	}

	return nil
}
