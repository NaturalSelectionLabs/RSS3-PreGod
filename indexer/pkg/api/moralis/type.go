package moralis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
)

type ChainType string

const (
	Unknown ChainType = "unknown"

	ETH     ChainType = "eth"
	BSC     ChainType = "bsc"
	Polygon ChainType = "polygon"
	AVAX    ChainType = "avalanche"
	Fantom  ChainType = "fantom"
)

func GetChainType(network constants.NetworkID) ChainType {
	switch network {
	case constants.NetworkIDEthereum:
		return ETH
	case constants.NetworkIDBNBChain:
		return BSC
	case constants.NetworkIDPolygon:
		return Polygon
	case constants.NetworkIDAvalanche:
		return AVAX
	case constants.NetworkIDFantom:
		return Fantom
	default:
		return Unknown
	}
}

func (mt ChainType) GetNetworkSymbol() constants.NetworkSymbol {
	switch mt {
	case ETH:
		return constants.NetworkSymbolEthereum
	case BSC:
		return constants.NetworkSymbolBNBChain
	case Polygon:
		return constants.NetworkSymbolPolygon
	case AVAX:
		return constants.NetworkSymbolAvalanche
	case Fantom:
		return constants.NetworkSymbolFantom
	default:
		return constants.NetworkSymbolUnknown
	}
}

type MoralisAttributes struct {
	// Speed limit property obtained from Moralis header
	MinRateLimit     int
	MinRateLimitUsed int
}

func SetMoralisAttributes(attributes *MoralisAttributes, response httpx.Response) {
	if attributes == nil {
		return
	}

	MinRateLimitStr := response.Header.Get("x-rate-limit-limit")
	if MinRateLimitStr != "" {
		attributes.MinRateLimit, _ = strconv.Atoi(MinRateLimitStr)
	}

	MinRateLimitUsedStr := response.Header.Get("x-rate-limit-used")
	if MinRateLimitUsedStr != "" {
		attributes.MinRateLimitUsed, _ = strconv.Atoi(MinRateLimitUsedStr)
	}
}

// NFTItem store all indexed NFTs from moralis api.
type NFTItem struct {
	MoralisAttributes

	TokenAddress      string `json:"token_address"`
	TokenId           string `json:"token_id"`
	BlockNumberMinted string `json:"block_number_minted"`
	OwnerOf           string `json:"owner_of"`
	BlockNumber       string `json:"block_number"`
	Amount            string `json:"amount"`
	ContractType      string `json:"contract_type"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	TokenURI          string `json:"token_uri"`
	MetaData          string `json:"metadata"`
	SyncedAt          string `json:"synced_at"`
	IsValid           int64  `json:"is_valid"`
	Syncing           int64  `json:"syncing"`
	Frozen            int64  `json:"frozen"`
}

// an incomplete set of metadata returned from Moralis, mapped only what are needed for now
type NFTMetadata struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attributes  []NFTMetadataAttribute `json:"attributes"`
}

type NFTMetadataAttribute struct {
	TraitType   string `json:"trait_type"`
	DisplayType string `json:"display_type"`
	Value       int64  `json:"value"`
}

type NFTResult struct {
	MoralisAttributes

	Total    int64     `json:"total"`
	Page     int64     `json:"page"`
	PageSize int64     `json:"page_size"`
	Result   []NFTItem `json:"result"`
	Status   string    `json:"status"`
}

func (i NFTItem) String() string {
	return fmt.Sprintf(`TokenAddress: %s, TokenId: %s, OwnerOf: %s, TokenURI: %s`,
		i.TokenAddress, i.TokenId, i.OwnerOf, i.TokenURI)
}

func (i NFTItem) GetAssetProof() string {
	return i.TokenAddress + "-" + i.TokenId
}

// NFTTransferItem store the transfers of NFTS.
type NFTTransferItem struct {
	MoralisAttributes

	BlockNumber      string `json:"block_number"`
	BlockTimestamp   string `json:"block_timestamp"`
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"transaction_hash"`
	TransactionIndex int64  `json:"transaction_index"`
	LogIndex         int64  `json:"log_index"`
	Value            string `json:"value"`
	ContractType     string `json:"contract_type"`
	TransactionType  string `json:"transaction_type"`
	TokenAddress     string `json:"token_address"`
	TokenId          string `json:"token_id"`
	FromAddress      string `json:"from_address"`
	ToAddress        string `json:"to_address"`
	Amount           string `json:"amount"`
	Verified         int64  `json:"verified"`
	Operator         string `json:"operator"`
}

type NFTTransferResult struct {
	MoralisAttributes

	Total       int64             `json:"total"`
	Page        int64             `json:"page"`
	PageSize    int64             `json:"page_size"`
	Result      []NFTTransferItem `json:"result"`
	Cursor      string            `json:"cursor"`
	BlockExists bool              `json:"block_exists"`
}

func (i NFTTransferItem) String() string {
	return fmt.Sprintf(`From: %s, To: %s, TokenAddress: %s, ContractType: %s, TokenId: %s`,
		i.FromAddress, i.ToAddress, i.TokenAddress, i.ContractType, i.TokenId)
}

func (i NFTTransferItem) EqualsToToken(nft NFTItem) bool {
	return i.TokenAddress == nft.TokenAddress && i.TokenId == nft.TokenId
}

func (i NFTItem) GetUid() string {
	return fmt.Sprintf("%s.%s", i.TokenAddress, i.TokenId)
}

func (i NFTTransferItem) GetUid() string {
	return fmt.Sprintf("%s.%s", i.TokenAddress, i.TokenId)
}

type ERC20Transfer struct {
	MoralisAttributes

	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Cursor   string              `json:"cursor"`
	Result   []ERC20TransferItem `json:"result"`
}

type ERC20TransferItem struct {
	MoralisAttributes

	TransactionHash string `json:"transaction_hash"`
	TokenAddress    string `json:"address"`
	BlockTimestamp  string `json:"block_timestamp"`
	BlockNumber     string `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	ToAddress       string `json:"to_address"`
	FromAddress     string `json:"from_address"`
	Value           string `json:"value"`
}

type Erc20TokenMetaDataItem struct {
	MoralisAttributes

	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    string `json:"decimals"`
	Logo        string `json:"logo"`
	LogoHash    string `json:"logo_hash"`
	Thumbnail   string `json:"thumbnail"`
	BlockNumber string `json:"block_number"`
	Validated   int    `json:"validated"`
	CreatedAt   string `json:"created_at"`
}

func (i ERC20TransferItem) String() string {
	return fmt.Sprintf(`From: %s, To: %s`,
		i.FromAddress, i.ToAddress)
}

type GetLogsItem struct {
	TransactionHash string `json:"transaction_hash"`
	Address         string `json:"address"`
	BlockTimestamp  string `json:"block_timestamp"`
	BlockNumber     string `json:"block_number"`
	BlockHash       string `json:"block_hash"`
	Data            string `json:"data"`
	Topic0          string `json:"topic0"`
	Topic1          string `json:"topic1"`
	Topic2          string `json:"topic2"`
	Topic3          string `json:"topic3"`
}

func (i GetLogsItem) String() string {
	return fmt.Sprintf(`TransactionHash: %s, TokenAddress: %s, Data: %s, Topic0: %s, Topic1: %s, Topic2:%s, Topic3: %s`,
		i.TransactionHash, i.Address, i.Data, i.Topic0, i.Topic1, i.Topic2, i.Topic3)
}

type GetLogsResult struct {
	MoralisAttributes

	Total    int64         `json:"total"`
	Page     int64         `json:"page"`
	PageSize int64         `json:"page_size"`
	Result   []GetLogsItem `json:"result"`
}

// Returns related urls based on the network and contract tx hash.
func GetTxRelatedURLs(
	network constants.NetworkSymbol,
	contractAddress string,
	tokenId string,
	transactionHash *string,
) []string {
	var urls []string
	if transactionHash != nil {
		urls = append(urls, GetTxHashURL(network, *transactionHash))
	}

	switch network {
	case constants.NetworkSymbolEthereum:
		urls = append(urls, "https://etherscan.io/nft/"+contractAddress+"/"+tokenId)
		urls = append(urls, "https://opensea.io/assets/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolPolygon:
		urls = append(urls, "https://polygonscan.com/token/"+contractAddress)
		urls = append(urls, "https://opensea.io/assets/matic/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolBNBChain:
		urls = append(urls, "https://bscscan.com/nft/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolAvalanche:
	case constants.NetworkSymbolFantom:
		urls = append(urls, "https://ftmscan.com/nft/"+contractAddress+"/"+tokenId)
	}

	return urls
}

// Returns related urls based on the network and contract tx hash.
func GetTxHashURL(
	network constants.NetworkSymbol,
	transactionHash string,
) string {
	switch network {
	case constants.NetworkSymbolEthereum:
		return "https://etherscan.io/tx/" + (transactionHash)

	case constants.NetworkSymbolPolygon:
		return "https://polygonscan.com/tx/" + (transactionHash)

	case constants.NetworkSymbolBNBChain:
		return "https://bscscan.com/tx/" + (transactionHash)

	case constants.NetworkSymbolAvalanche:
		return "https://avascan.info/blockchain/c/tx/" + (transactionHash)
	case constants.NetworkSymbolFantom:
		return "https://ftmscan.com/tx/" + (transactionHash)
	default:
		return ""
	}
}

// for ENS only
type ENSTextRecord struct {
	Domain      string
	Description string
	Text        map[string]string
	Avatar      string
	Attachments datatype.Attachments
	CreatedAt   time.Time
	TxHash      string
}

// returns a list of recommended keys for a given ENS domain, as per https://app.ens.domains/
// this is a combination of Global Keys and Service Keys, see: https://eips.ethereum.org/EIPS/eip-634
func getTextRecordKeyList() []string {
	// nolint:lll // this is a list of keys
	return []string{"email", "url", "avatar", "description", "notice", "keywords", "com.discord", "com.github", "com.reddit", "com.twitter", "org.telegram", "eth.ens.delegate"}
}

func GetTsp(blockTimestamp string) (time.Time, error) {
	if t, err := time.Parse(time.RFC1123, blockTimestamp); err != nil {
		// try another format
		if t, err = time.Parse(time.RFC3339, blockTimestamp); err != nil {
			// try another format
			return time.Parse("2006-01-02T15:04:05.000Z", blockTimestamp)
		} else {
			return t, err
		}
	} else {
		return t, err
	}
}
