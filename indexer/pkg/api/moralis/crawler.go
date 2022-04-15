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
	"github.com/ethereum/go-ethereum/ethclient"
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

func getInfuraClient() {
	c, err := ethclient.Dial(infuraGateway + "/" + config.Config.Indexer.Infura.ApiKey)

	if err != nil {
		logger.Errorf("connect to Infura: %v", err)
	}

	client = c
}

//nolint:funlen,gocognit // disable line length check
func (c *moralisCrawler) Work(param crawler.WorkParam) error {
	chainType := GetChainType(param.NetworkID)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()

	// nftTransfers for notes
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, param.BlockHeight, getApiKey())
	if err != nil {
		return err
	}

	// get nft for assets
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

	// those two should be expected to be equal actually
	owner := rss3uri.NewAccountInstance(param.OwnerID, param.OwnerPlatformID.Symbol()).UriString()
	author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolEthereum).UriString()

	// complete the note list
	for _, item := range nftTransfers.Result {
		tsp, tspErr := item.GetTsp()
		if tspErr != nil {
			logger.Warnf("asset: %s fails at GetTsp(): %v", item.String(), tspErr)

			tsp = time.Now()
		}

		hasObject := false

		var theAsset NFTItem

		for _, asset := range assets.Result {
			if item.EqualsToToken(asset) && asset.MetaData != "" {
				hasObject = true
				theAsset = asset
			}
		}

		m, parseErr := utils.ParseNFTMetadata(theAsset.MetaData)
		if parseErr != nil {
			logger.Warnf("%v", parseErr)
		}

		if !hasObject {
			theAsset, err = GetMetadataByToken(item.TokenAddress, item.TokenId, chainType, getApiKey())
			if err != nil {
				logger.Warnf("fail to get metadata of token: " + item.String())
			}
		}

		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(item.TransactionHash, networkSymbol).UriString(),
			Owner:           owner,
			RelatedURLs:     GetTxRelatedURLs(networkSymbol, item.TokenAddress, item.TokenId, &item.TransactionHash),
			Tags:            constants.ItemTagsNFT.ToPqStringArray(),
			Authors:         []string{author},
			Title:           m.Name,
			Summary:         m.Description,
			Attachments:     database.MustWrapJSON(utils.Meta2NoteAtt(m)),
			Source:          constants.NoteSourceNameEthereumNFT.String(),
			MetadataNetwork: networkSymbol.String(),
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

		c.Notes = append(c.Notes, note)
	}

	// complete the asset list
	for _, asset := range assets.Result {
		m, parseErr := utils.ParseNFTMetadata(asset.MetaData)
		if parseErr != nil {
			logger.Warnf("%v", parseErr)
		}

		// find the note that has the same proof to get the tsp
		var tsp time.Time

		for _, note := range c.Notes {
			noteMetadata, unwrapErr := database.UnwrapJSON[map[string]interface{}](note.Metadata)
			if unwrapErr != nil {
				logger.Warnf("%v", unwrapErr) // should never be a problem

				continue
			}

			if noteMetadata["collection_address"] == asset.TokenAddress &&
				noteMetadata["token_id"] == asset.TokenId {
				tsp = note.DateCreated

				break
			}
		}

		proof := asset.TokenAddress + "-" + asset.TokenId
		asset := model.Asset{
			Identifier:      rss3uri.NewAssetInstance(proof, networkSymbol).UriString(),
			Owner:           owner,
			RelatedURLs:     GetTxRelatedURLs(networkSymbol, asset.TokenAddress, asset.TokenId, nil),
			Tags:            constants.ItemTagsNFT.ToPqStringArray(),
			Authors:         []string{author},
			Title:           m.Name,
			Summary:         m.Description,
			Attachments:     database.MustWrapJSON(utils.Meta2AssetAtt(m)),
			Source:          constants.AssetSourceNameEthereumNFT.String(),
			MetadataNetwork: string(networkSymbol),
			MetadataProof:   proof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"token_standard":     asset.ContractType,
				"token_id":           asset.TokenId,
				"token_symbol":       asset.Symbol,
				"collection_address": asset.TokenAddress,
				"collection_name":    m.Name,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		c.Assets = append(c.Assets, asset)
	}

	// find old data in the database
	// TODO: need to find a better way to do this
	//newAssetProofs := lo.Map(c.Assets, func(asset model.Asset, _ int) string {
	//	return asset.MetadataProof
	//})
	//oldAssets, err := database.QueryAllAssets(database.DB, []string{owner}, networkSymbol)
	//if err != nil {
	//	logger.Warnf("fail to query old assets: %v", err)
	//} else {
	//	lop.ForEach(oldAssets, func(oldAsset model.Asset, _ int) {
	//		if !lo.Contains(newAssetProofs, oldAsset.MetadataProof) {
	//			// remove this old asset from database
	//			if _, err := database.DeleteAsset(database.DB, &oldAsset); err != nil {
	//				logger.Warnf("fail to delete old asset: %v", err)
	//			}
	//		}
	//	})
	//}

	if err := utils.CompleteMimeTypesForItems(c.Notes, c.Assets); err != nil {
		logger.Error("moralis complete mime types error:", err)
	}

	ensList, err := GetENSList(param.Identity)

	if err != nil {
		return err
	}

	for _, ens := range ensList {
		metadata := make(map[string]interface{}, len(ens.Text))
		for k, v := range ens.Text {
			metadata[k] = v
		}

		profile := model.Profile{
			ID:          param.Identity,
			Platform:    constants.PlatformIDEthereum.Int(),
			Source:      constants.ProfileSourceIDENS.Int(),
			Name:        database.WrapNullString(ens.Domain),
			Bio:         database.WrapNullString(ens.Description),
			Avatars:     []string{ens.Text["avatar"]},
			Attachments: ens.Attachments,
		}

		c.Profiles = append(c.Profiles, profile)
	}

	return nil
}
