package moralis

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type moralisCrawler struct {
	crawler.CrawlerResult
}

func NewMoralisCrawler() crawler.Crawler {
	return &moralisCrawler{
		crawler.CrawlerResult{
			Items:  []*model.Item{},
			Assets: []*model.ItemId{},
			Notes:  []*model.ItemId{},
		},
	}
}

//nolint:funlen // disable line length check
func (mc *moralisCrawler) Work(userAddress string, network constants.NetworkID) error {
	chainType := GetChainType(network)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()
	networkId := networkSymbol.GetID()
	nftTransfers, err := GetNFTTransfers(userAddress, chainType, GetApiKey())

	if err != nil {
		return err
	}

	//TODO: tsp
	assets, err := GetNFTs(userAddress, chainType, GetApiKey())
	if err != nil {
		return err
	}
	//parser
	for _, nftTransfer := range nftTransfers.Result {
		mc.Notes = append(mc.Notes, &model.ItemId{
			NetworkId: networkId,
			Proof:     nftTransfer.TransactionHash,
		})
	}

	for _, asset := range assets.Result {
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true

				mc.Assets = append(mc.Assets, &model.ItemId{
					NetworkId: networkId,
					Proof:     nftTransfer.TransactionHash,
				})
			}
		}

		if !hasProof {
			// TODO: error handle here
			logger.Errorf("Asset doesn't has proof.")
		}
	}

	author, err := rss3uri.NewInstance("account", userAddress, string(constants.PlatformSymbolEthereum))

	if err != nil {
		// TODO
		logger.Error(err)
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
		mc.Items = append(mc.Items, ni)
	}

	return nil
}

func (mc *moralisCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: mc.Assets,
		Notes:  mc.Notes,
		Items:  mc.Items,
	}
}
