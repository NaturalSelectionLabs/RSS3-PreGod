package moralis

import (
	"fmt"
	"log"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
)

type moralisCrawler struct {
	rss3Items   []*model.Item
	rss3Objects []*model.Object

	rss3Assets, rss3Notes []*model.ItemId
}

func NewMoralisCrawler() crawler.Crawler {
	return &moralisCrawler{
		rss3Items:   []*model.Item{},
		rss3Objects: []*model.Object{},

		rss3Assets: []*model.ItemId{},
		rss3Notes:  []*model.ItemId{},
	}
}

//nolint:funlen // disable line length check
func (mc *moralisCrawler) Work(userAddress string, itemType constants.NetworkName) error {
	chainType := GetChainType(itemType)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", itemType)
	}

	itemTypeID := chainType.GetNFTItemTypeID()
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
		tsp, err := nftTransfer.GetTsp()
		if err != nil {
			// TODO: log error
			logger.Error(tsp, err)
			tsp = time.Now()
		}

		ni := model.NewItem(
			nftTransfer.GetUid(),
			itemTypeID,
			nftTransfer.FromAddress,
			nftTransfer.ToAddress,
			nftTransfer.TransactionHash,
			tsp,
		)
		mc.rss3Items = append(mc.rss3Items, ni)
		mc.rss3Notes = append(mc.rss3Notes, &model.ItemId{
			ItemTypeID: itemTypeID,
			Proof:      nftTransfer.TransactionHash,
		})
	}

	for _, asset := range assets.Result {
		// TODO: make attachments and authors
		no := model.NewObject(nil, asset.GetUid(), itemTypeID, "", "", nil, nil)
		mc.rss3Objects = append(mc.rss3Objects, no)
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true

				mc.rss3Assets = append(mc.rss3Assets, &model.ItemId{
					ItemTypeID: itemTypeID,
					Proof:      nftTransfer.TransactionHash,
				})
			}
		}

		if !hasProof {
			// TODO: error handle here
			logger.Errorf("Asset doesn't has proof.")
		}
	}
	// make the object list complete
	for _, nftTransfer := range nftTransfers.Result {
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
	}

	return nil
}

func (mc *moralisCrawler) GetResult() *crawler.CrawlerResult {
	return &crawler.CrawlerResult{
		Assets:  mc.rss3Assets,
		Notes:   mc.rss3Notes,
		Items:   mc.rss3Items,
		Objects: mc.rss3Objects,
	}
}
