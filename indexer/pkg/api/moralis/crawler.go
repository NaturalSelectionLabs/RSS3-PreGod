package moralis

import (
	"fmt"
	"strings"
	"time"

	utils "github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/nft_utils"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"golang.org/x/sync/errgroup"
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
		tsp, err := item.GetTsp()
		if err != nil {
			logger.Warnf("asset: %s fails at GetTsp(): %v", item.String(), err)

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
		proof := asset.TokenAddress + "-" + asset.TokenId

		m, err := utils.ParseNFTMetadata(asset.MetaData)
		if err != nil {
			logger.Warnf("%v", err)
		}

		// find the note that has the same proof to get the tsp
		var tsp time.Time

		for _, note := range c.Notes {
			if note.MetadataProof == proof {
				tsp = note.DateCreated

				break
			}
		}

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
	newAssetProofs := lo.Map(c.Assets, func(asset model.Asset, _ int) string {
		return asset.MetadataProof
	})

	oldAssets, err := database.QueryAllAssets(database.DB, []string{owner})
	if err != nil {
		logger.Warnf("fail to query old assets: %v", err)
	} else {
		lop.ForEach(oldAssets, func(oldAsset model.Asset, _ int) {
			if !lo.Contains(newAssetProofs, oldAsset.MetadataProof) {
				// remove this old asset from database
				if _, err := database.DeleteAsset(database.DB, &oldAsset); err != nil {
					logger.Warnf("fail to delete old asset: %v", err)
				}
			}
		})
	}

	// complete attachments in parallel
	g := new(errgroup.Group)

	g.Go(func() error {
		lop.ForEach(c.Notes, func(note model.Note, i int) {
			if note.Attachments != nil {
				as, err := database.UnwrapJSON[datatype.Attachments](note.Attachments)
				if err != nil {
					return
				}
				utils.CompleteMimeTypes(as)
				c.Notes[i].Attachments = database.MustWrapJSON(as)
			}
		})

		return nil
	})

	g.Go(func() error {
		lop.ForEach(c.Assets, func(asset model.Asset, i int) {
			if asset.Attachments != nil {
				as, err := database.UnwrapJSON[datatype.Attachments](asset.Attachments)
				if err != nil {
					return
				}
				utils.CompleteMimeTypes(as)
				c.Assets[i].Attachments = database.MustWrapJSON(as)
			}
		})

		return nil
	})

	_ = g.Wait()

	return nil
}
