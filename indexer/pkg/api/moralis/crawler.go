package moralis

import (
	"fmt"
	"strings"
	"time"

	utils "github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/nft_utils"
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
			Assets:     []*model.ObjectId{},
			Notes:      []*model.ObjectId{},
			Items:      []*model.Item{},
			AssetItems: []*model.Item{},
		},
	}
}

func getApiKey() string {
	apiKey, err := jsoni.MarshalToString(config.Config.Indexer.Moralis.ApiKey)
	if err != nil {
		return ""
	}

	return strings.Trim(apiKey, "\"")
}

//nolint:funlen // disable line length check
func (c *moralisCrawler) Work(param crawler.WorkParam) error {
	chainType := GetChainType(param.NetworkID)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()
	networkId := networkSymbol.GetID()
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, getApiKey())

	if err != nil {
		return err
	}

	//TODO: tsp
	assets, err := GetNFTs(param.Identity, chainType, getApiKey())
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
			logger.Warnf("Asset: " + asset.String() + " doesn't has proof.")
		}
	}

	// make the item list complete
	for _, nftTransfer := range nftTransfers.Result {
		tsp, err := nftTransfer.GetTsp()
		if err != nil {
			logger.Warnf("asset: %s fails at GetTsp(): %v", nftTransfer.String(), err)

			tsp = time.Now()
		}

		author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolEthereum)

		hasObject := false
		var theAsset NFTItem

		for _, asset := range assets.Result {
			if nftTransfer.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
				theAsset = asset
			}
		}
		m, err := utils.ParseNFTMetadata(theAsset.MetaData)
		if err != nil {
			logger.Warnf("%v", err)
		}

		if !hasObject {
			theAsset, err = GetMetadataByToken(nftTransfer.TokenAddress, nftTransfer.TokenId, chainType, getApiKey())
			if err != nil {
				logger.Warnf("fail to get metadata of token: " + nftTransfer.String())
			}
		}

		noteItem := model.NewItem(
			networkId,
			nftTransfer.TransactionHash,
			model.Metadata{
				"from":           nftTransfer.FromAddress,
				"to":             nftTransfer.ToAddress,
				"token_standard": nftTransfer.ContractType,
				"token_id":       nftTransfer.TokenId,
				"token_symbol":   theAsset.Symbol,

				"collection_address": nftTransfer.TokenAddress,
				"collection_name":    theAsset.Name,
			},
			constants.ItemTagsNFT,
			[]string{author.String()},
			"", // title
			"", // summary
			utils.Meta2AssetAtt(m),
			tsp,
		)

		assetItem := model.NewItem(
			networkId,
			theAsset.GetAssetProof(),
			model.Metadata{
				"token_standard": nftTransfer.ContractType,
				"token_id":       nftTransfer.TokenId,
				"token_symbol":   theAsset.Symbol,

				"collection_address": nftTransfer.TokenAddress,
				"collection_name":    theAsset.Name,
			},
			constants.ItemTagsNFT,
			[]string{author.String()},
			m.Name,        // title
			m.Description, // summary
			utils.Meta2NoteAtt(m),
			tsp,
		)

		c.Items = append(c.Items, noteItem)

		c.AssetItems = append(c.AssetItems, assetItem)

	}

	return nil
}
