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

//nolint:funlen // disable line length check
func Crawl(param *crawler.WorkParam, result *crawler.CrawlerResult) (crawler.CrawlerResult, error) {
	chainType := GetChainType(param.NetworkID)
	if chainType == Unknown {
		return *result, fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()
	networkId := networkSymbol.GetID()
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, GetApiKey())

	if err != nil {
		return *result, err
	}

	//TODO: tsp
	assets, err := GetNFTs(param.Identity, chainType, GetApiKey())
	if err != nil {
		return *result, err
	}
	//parser
	for _, nftTransfer := range nftTransfers.Result {
		result.Notes = append(result.Notes, &model.ItemId{
			NetworkID: networkId,
			Proof:     nftTransfer.TransactionHash,
		})
	}

	for _, asset := range assets.Result {
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true

				result.Assets = append(result.Assets, &model.ItemId{
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
		result.Items = append(result.Items, ni)
	}

	return *result, nil
}
