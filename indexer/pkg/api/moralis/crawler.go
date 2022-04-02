package moralis

import (
	"fmt"
	"strings"
	"time"

	utils "github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/nft_utils"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
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
			Assets: []model.Asset{},
			Notes:  []model.Note{},
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
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, getApiKey())

	if err != nil {
		return err
	}

	//TODO: tsp
	assets, err := GetNFTs(param.Identity, chainType, getApiKey())
	if err != nil {
		return err
	}

	// check if each asset has a proof (only for logging issues)
	for _, asset := range assets.Result {
		hasProof := false

		for _, nftTransfer := range nftTransfers.Result {
			if nftTransfer.EqualsToToken(asset) {
				hasProof = true
			}
		}

		if !hasProof {
			logger.Warnf("Asset: " + asset.String() + " doesn't has proof.")
		}
	}

	// make the item list complete
	for _, item := range nftTransfers.Result {
		tsp, err := item.GetTsp()
		if err != nil {
			logger.Warnf("asset: %s fails at GetTsp(): %v", item.String(), err)

			tsp = time.Now()
		}

		author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolEthereum)

		hasObject := false

		var theAsset NFTItem

		for _, asset := range assets.Result {
			if item.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
				theAsset = asset
			}
		}

		m, err := utils.ParseNFTMetadata(theAsset.MetaData)
		if err != nil {
			logger.Warnf("%v", err)
		}

		if !hasObject {
			theAsset, err = GetMetadataByToken(item.TokenAddress, item.TokenId, chainType, getApiKey())
			if err != nil {
				logger.Warnf("fail to get metadata of token: " + item.String())
			}
		}

		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(item.TransactionHash, networkSymbol).String(),
			Owner:           author.String(),
			RelatedURLs:     GetTxRelatedURLs(networkSymbol, item.TokenAddress, item.TokenId, &item.TransactionHash),
			Tags:            constants.ItemTagsNFT.ToPqStringArray(),
			Authors:         []string{author.String()},
			Title:           m.Name,
			Summary:         m.Description,
			Attachments:     database.MustWrapJSON(utils.Meta2NoteAtt(m)),
			Source:          constants.NoteSourceNameEthereumNFT.String(),
			MetadataNetwork: string(networkSymbol),
			MetadataProof:   item.TransactionHash,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"from":               item.FromAddress,
				"to":                 item.ToAddress,
				"token_standard":     item.ContractType,
				"token_id":           item.TokenId,
				"token_symbol":       theAsset.Symbol,
				"collection_address": item.TokenAddress,
				"collection_name":    theAsset.Name,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		assetProof := item.GetAssetProof()
		asset := model.Asset{
			Identifier:      rss3uri.NewAssetInstance(assetProof, networkSymbol).String(),
			Owner:           author.String(),
			RelatedURLs:     GetTxRelatedURLs(networkSymbol, item.TokenAddress, item.TokenId, nil),
			Tags:            constants.ItemTagsNFT.ToPqStringArray(),
			Authors:         []string{author.String()},
			Title:           m.Name,
			Summary:         m.Description,
			Attachments:     database.MustWrapJSON(utils.Meta2AssetAtt(m)),
			Source:          constants.AssetSourceNameEthereumNFT.String(),
			MetadataNetwork: string(networkSymbol),
			MetadataProof:   assetProof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"token_standard":     item.ContractType,
				"token_id":           item.TokenId,
				"token_symbol":       theAsset.Symbol,
				"collection_address": item.TokenAddress,
				"collection_name":    theAsset.Name,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		c.Notes = append(c.Notes, note)
		c.Assets = append(c.Assets, asset)
	}

	return nil
}
