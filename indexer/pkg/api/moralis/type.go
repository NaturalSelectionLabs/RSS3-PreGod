package moralis

import (
	"fmt"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
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

// NFTItem store all indexed NFTs from moralis api.
type NFTItem struct {
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

func (i NFTTransferItem) GetTsp() (time.Time, error) {
	return time.Parse(time.RFC3339, i.BlockTimestamp)
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
	urls := []string{}
	if transactionHash != nil {
		urls = append(urls, GetTxHashURL(network, *transactionHash))
	}

	switch network {
	case constants.NetworkSymbolEthereum:
		if transactionHash != nil {
			urls = append(urls, "https://etherscan.io/tx/"+(*transactionHash))
		}

		urls = append(urls, "https://etherscan.io/nft/"+contractAddress+"/"+tokenId)
		urls = append(urls, "https://opensea.io/assets/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolPolygon:
		if transactionHash != nil {
			urls = append(urls, "https://polygonscan.com/tx/"+(*transactionHash))
		}

		urls = append(urls, "https://polygonscan.com/nft/"+contractAddress+"/"+tokenId)
		urls = append(urls, "https://opensea.io/assets/matic/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolBNBChain:
		if transactionHash != nil {
			urls = append(urls, "https://bscscan.com/tx/"+(*transactionHash))
		}

		urls = append(urls, "https://bscscan.com/nft/"+contractAddress+"/"+tokenId)
	case constants.NetworkSymbolAvalanche:
		if transactionHash != nil {
			urls = append(urls, "https://avascan.info/blockchain/c/tx/"+(*transactionHash))
		}
	case constants.NetworkSymbolFantom:
		if transactionHash != nil {
			urls = append(urls, "https://ftmscan.com/tx/"+(*transactionHash))
		}

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
