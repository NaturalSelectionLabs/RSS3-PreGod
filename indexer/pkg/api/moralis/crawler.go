package moralis

import (
	"fmt"
	"math/big"
	"strconv"
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
	mapset "github.com/deckarep/golang-set"
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

func getGatewayClient() {
	c, err := ethclient.Dial(config.Config.Indexer.Gateway.Endpoint)

	if err != nil {
		logger.Errorf("connect to json rpc endpoint error: %v", err)
	}

	client = c
}

//nolint:funlen,gocognit // disable line length check
func (c *moralisCrawler) setNFTTransfers(
	param crawler.WorkParam,
	owner string,
	author string,
	networkSymbol constants.NetworkSymbol,
	chainType ChainType) error {
	// nftTransfers for notes
	nftTransfers, err := GetNFTTransfers(param.Identity, chainType, param.BlockHeight, getApiKey())
	if err != nil {
		logger.Errorf("get nft transfers: %v", err)

		return err
	}

	// get nft for assets
	assets, err := GetNFTs(param.Identity, chainType, getApiKey())
	if err != nil {
		logger.Errorf("get nft: %v", err)

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

	// complete the note list
	for _, item := range nftTransfers.Result {
		tsp, tspErr := GetTsp(item.BlockTimestamp)
		if tspErr != nil {
			logger.Warnf("asset: %s fails at GetTsp(): %v", item.String(), tspErr)

			tsp = time.Now()
		}

		var theAsset NFTItem

		for _, asset := range assets.Result {
			if item.EqualsToToken(asset) && asset.MetaData != "" {
				theAsset = asset
			}
		}

		if theAsset.MetaData == "" {
			theAsset, err = GetMetadataByToken(item.TokenAddress, item.TokenId, chainType, getApiKey())
			if err != nil {
				logger.Warnf("fail to get metadata of token [" + item.String() + "] err[" + err.Error() + "]")
			}
		}

		m, parseErr := utils.ParseNFTMetadata(theAsset.MetaData)
		if parseErr != nil {
			logger.Warnf("%v", parseErr)
		}

		//convert to string
		proof := item.TransactionHash + "-" + strconv.FormatInt(item.LogIndex, 10) + "-" + item.TokenId
		note := model.Note{
			Identifier:      rss3uri.NewNoteInstance(proof, networkSymbol).UriString(),
			Owner:           owner,
			RelatedURLs:     GetTxRelatedURLs(networkSymbol, item.TokenAddress, item.TokenId, &item.TransactionHash),
			Tags:            constants.ItemTagsNFT.ToPqStringArray(),
			Authors:         []string{author},
			Title:           m.Name,
			Summary:         m.Description,
			Attachments:     database.MustWrapJSON(utils.Meta2NoteAtt(m)),
			Source:          constants.NoteSourceNameEthereumNFT.String(),
			ContractAddress: item.TokenAddress,
			LogIndex:        int(item.LogIndex),
			TokenID:         item.TokenId,
			MetadataNetwork: networkSymbol.String(),
			MetadataProof:   proof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"from":               strings.ToLower(item.FromAddress),
				"to":                 strings.ToLower(item.ToAddress),
				"token_standard":     item.ContractType,
				"token_id":           item.TokenId,
				"token_symbol":       theAsset.Symbol,
				"collection_address": strings.ToLower(item.TokenAddress),
				"collection_name":    theAsset.Name,
				"log_index":          item.LogIndex,
				"contract_type":      item.ContractType,
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
			ContractAddress: asset.TokenAddress,
			TokenID:         asset.TokenId,
			MetadataNetwork: string(networkSymbol),
			MetadataProof:   proof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"token_standard":     asset.ContractType,
				"token_id":           asset.TokenId,
				"token_symbol":       asset.Symbol,
				"collection_address": strings.ToLower(asset.TokenAddress),
				"collection_name":    m.Name,
				"contract_type":      asset.ContractType,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		c.Assets = append(c.Assets, asset)
	}

	return nil
}

// ERC20 used
type noteInstanceBuilder struct {
	countMap map[string]int
}

func getNewNoteInstanceBuilder() *noteInstanceBuilder {
	return &noteInstanceBuilder{
		countMap: map[string]int{},
	}
}

func setNoteInstance(
	niBuilder *noteInstanceBuilder,
	txHash string,
) (string, error) {
	if niBuilder == nil {
		return "", fmt.Errorf("note instance builder is nil")
	}

	if txHash == "" {
		return "", fmt.Errorf("tx hash is empty")
	}

	hashCount, ok := niBuilder.countMap[txHash]
	if !ok {
		niBuilder.countMap[txHash] = 0

		return txHash + "-0", nil
	}

	hashCount += 1

	niBuilder.countMap[txHash] = hashCount

	return txHash + "-" + strconv.Itoa(hashCount), nil
}

// nolint:funlen  // disable line length check
func (c *moralisCrawler) setERC20(
	param crawler.WorkParam,
	owner string,
	author string,
	networkSymbol constants.NetworkSymbol,
	chainType ChainType) error {
	result, err := GetErc20Transfers(param.Identity, chainType, getApiKey())
	if err != nil {
		logger.Errorf("chain type[%s], get erc20 transfers: %v", chainType.GetNetworkSymbol().String(), err)

		return err
	}

	// get the token address
	tokenAddressSet := mapset.NewSet()
	tokenAddresses := []string{}

	for _, item := range result {
		tokenAddressSet.Add(item.TokenAddress)
	}

	for _, tokenAddress := range tokenAddressSet.ToSlice() {
		addressStr, ok := tokenAddress.(string)
		if !ok {
			logger.Warnf("token address[%v] is not string", addressStr)

			continue
		}

		tokenAddresses = append(tokenAddresses, addressStr)
	}

	// get the token metadata
	erc20Tokens, err := GetErc20TokenMetaData(chainType, tokenAddresses, getApiKey())
	if err != nil {
		logger.Errorf("chain type[%s], get erc20 token metadata [%v]",
			chainType.GetNetworkSymbol().String(), err)

		return err
	}

	niBuilder := getNewNoteInstanceBuilder()

	// complete the note list
	for _, item := range result {
		tsp, tspErr := GetTsp(item.BlockTimestamp)
		if tspErr != nil {
			logger.Warnf("chain type[%s], item[%s], fails at GetTsp err[%v]",
				chainType.GetNetworkSymbol().String(),
				item.String(), tspErr)

			tsp = time.Now()
		}

		m := erc20Tokens[item.TokenAddress]

		proof, err := setNoteInstance(niBuilder, item.TransactionHash)
		if err != nil {
			logger.Warnf("chain type[%s], item[%s], get instance key err[%v]",
				chainType.GetNetworkSymbol().String(),
				item.TransactionHash, err)

			continue
		}

		decimals, err := strconv.Atoi(m.Decimals)
		if err != nil {
			logger.Warnf("chain type[%s], item[%s], get decimal err[%v]",
				chainType.GetNetworkSymbol().String(),
				item.TransactionHash, err)
		}

		note := model.Note{
			Identifier: rss3uri.NewNoteInstance(proof, networkSymbol).UriString(),
			Owner:      owner,
			RelatedURLs: []string{
				"https://etherscan.io/tx/" + item.TransactionHash,
			},
			Tags:            constants.ItemTagsToken.ToPqStringArray(),
			Authors:         []string{author},
			Source:          constants.NoteSourceNameEthereumERC20.String(),
			ContractAddress: item.TokenAddress,
			MetadataNetwork: networkSymbol.String(),
			MetadataProof:   proof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"network":          networkSymbol.String(),
				"from":             strings.ToLower(item.FromAddress),
				"to":               strings.ToLower(item.ToAddress),
				"amount":           item.Value,
				"decimal":          decimals,
				"token_standard":   "ERC20",
				"token_symbol":     m.Symbol,
				"token_address":    strings.ToLower(item.TokenAddress),
				"transaction_hash": item.TransactionHash,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		c.Notes = append(c.Notes, note)
	}

	return nil
}

func (c *moralisCrawler) setETH(
	param crawler.WorkParam,
	owner string,
	author string,
	networkSymbol constants.NetworkSymbol,
	chainType ChainType) error {
	result, err := GetEthTransfers(param.Identity, chainType, getApiKey())
	if err != nil {
		logger.Errorf("chain type[%s], get eth transfers: %v", chainType.GetNetworkSymbol().String(), err)

		return err
	}

	if len(result) <= 0 {
		return nil
	}

	niBuilder := getNewNoteInstanceBuilder()

	for _, item := range result {
		tsp, tspErr := GetTsp(item.BlockTimestamp)
		if tspErr != nil {
			logger.Warnf("chain type[%s], item[%s], fails at GetTsp err[%v]",
				chainType.GetNetworkSymbol().String(),
				item.String(), tspErr)

			tsp = time.Now()
		}

		// only value > 0 to pass
		formatedAmount := big.NewInt(0)
		formatedAmount.SetString(item.Value, 10)

		if formatedAmount.Cmp(big.NewInt(0)) != 1 {
			continue
		}

		proof, err := setNoteInstance(niBuilder, item.TransactionHash)
		if err != nil {
			logger.Warnf("chain type[%s], item[%s], get instance key err[%v]",
				chainType.GetNetworkSymbol().String(),
				item.TransactionHash, err)

			continue
		}

		note := model.Note{
			Identifier: rss3uri.NewNoteInstance(proof, networkSymbol).UriString(),
			Owner:      owner,
			RelatedURLs: []string{
				"https://etherscan.io/tx/" + item.TransactionHash,
			},
			Tags:            constants.ItemTagsETH.ToPqStringArray(), // will be change
			Authors:         []string{author},
			Source:          constants.NoteSourceNameEthereumETH.String(),
			ContractAddress: "0x0",
			MetadataNetwork: networkSymbol.String(),
			MetadataProof:   proof,
			Metadata: database.MustWrapJSON(map[string]interface{}{
				"network":          networkSymbol.String(),
				"from":             strings.ToLower(item.FromAddress),
				"to":               strings.ToLower(item.ToAddress),
				"amount":           item.Value,
				"token_standard":   "ERC20",
				"token_symbol":     "ETH",
				"transaction_hash": item.TransactionHash,
			}),
			DateCreated: tsp,
			DateUpdated: tsp,
		}

		c.Notes = append(c.Notes, note)
	}

	return nil
}

func (c *moralisCrawler) Work(param crawler.WorkParam) error {
	chainType := GetChainType(param.NetworkID)
	if chainType == Unknown {
		return fmt.Errorf("unsupported network: %s", chainType)
	}

	networkSymbol := chainType.GetNetworkSymbol()

	// those two should be expected to be equal actually
	owner := rss3uri.NewAccountInstance(param.OwnerID, param.OwnerPlatformID.Symbol()).UriString()
	author := rss3uri.NewAccountInstance(param.Identity, constants.PlatformSymbolEthereum).UriString()

	// err := c.setNFTTransfers(param, owner, author, networkSymbol, chainType)
	// if err != nil {
	// 	logger.Errorf("fail to set nft transfers in db: %v", err)

	// 	return err
	// }

	// err = c.setERC20(param, owner, author, networkSymbol, chainType)
	// if err != nil {
	// 	logger.Errorf("fail to set erc20 in db: %v", err)

	// 	return err
	// }

	err := c.setETH(param, owner, author, networkSymbol, chainType)
	if err != nil {
		logger.Errorf("fail to set eth in db: %v", err)

		return err
	}

	// check duplicates in assets
	for i := 0; i < len(c.Assets); i++ {
		for j := i + 1; j < len(c.Assets); j++ {
			if c.Assets[i].Identifier == c.Assets[j].Identifier {
				logger.Errorf("Duplicate asset found: %v!!! This is temporarily removed.", c.Assets[i].Identifier)
				c.Assets = append(c.Assets[:j], c.Assets[j+1:]...)
				j--
			}
		}
	}

	// check duplicates in notes
	for i := 0; i < len(c.Notes); i++ {
		for j := i + 1; j < len(c.Notes); j++ {
			if c.Notes[i].Identifier == c.Notes[j].Identifier {
				logger.Errorf("Duplicate note found: %v!!! This is temporarily removed.", c.Notes[i].Identifier)
				c.Notes = append(c.Notes[:j], c.Notes[j+1:]...)
				j--
			}
		}
	}

	if err := utils.CompleteMimeTypesForItems(c.Notes, c.Assets, c.Profiles); err != nil {
		logger.Error("moralis complete mime types error:", err)
	}

	return nil
}
