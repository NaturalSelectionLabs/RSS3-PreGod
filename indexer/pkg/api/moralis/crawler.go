package moralis

import (
	"fmt"
	"log"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type moralisCrawler struct {
	rss3Items []*model.Item

	rss3Assets, rss3Notes []*model.ItemId
}

func NewMoralisCrawler() crawler.Crawler {
	return &moralisCrawler{
		rss3Items: []*model.Item{},

		rss3Assets: []*model.ItemId{},
		rss3Notes:  []*model.ItemId{},
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

	log.Println(nftTransfers.Total)
	//TODO: tsp
	assets, err := GetNFTs(userAddress, chainType, GetApiKey())
	if err != nil {
		return err
	}
	//parser
	for _, nftTransfer := range nftTransfers.Result {
		mc.rss3Notes = append(mc.rss3Notes, &model.ItemId{
			NetworkId: networkId,
			Proof:     nftTransfer.TransactionHash,
		})
	}

	for _, asset := range assets.Result {
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true

				mc.rss3Assets = append(mc.rss3Assets, &model.ItemId{
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
	// make the item list complete
	for _, nftTransfer := range nftTransfers.Result {
		// TODO: make attachments
		tsp, err := nftTransfer.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		author, err := rss3uri.NewInstance("account", nftTransfer.FromAddress, string(networkSymbol))
		if err != nil {
			// TODO
			logger.Error(tsp, err)
		}

		hasObject := false
		attachments := []model.Attachment{}

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
			attachments,
			tsp,
		)
		mc.rss3Items = append(mc.rss3Items, ni)
	}

	return nil
}

func (mc *moralisCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets: mc.rss3Assets,
		Notes:  mc.rss3Notes,
		Items:  mc.rss3Items,
	}
}
