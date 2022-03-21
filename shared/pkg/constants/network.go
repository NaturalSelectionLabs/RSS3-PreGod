package constants

type (
	NetworkID     int32
	NetworkSymbol string
)

const (
	NetworkIDUnknown         NetworkID = 0
	NetworkIDEthereumMainnet NetworkID = 1
	NetworkIDPolygon         NetworkID = 2
	NetworkIDBNBChain        NetworkID = 3
	NetworkIDArbitrum        NetworkID = 4
	NetworkIDAvalanche       NetworkID = 5
	NetworkIDFantom          NetworkID = 6
	NetworkIDGnosisMainnet   NetworkID = 7
	NetworkIDSolanaMainet    NetworkID = 8
	NetworkIDFlowMainnet     NetworkID = 9
	NetworkIDArweaveMainnet  NetworkID = 10
	NetworkIDRSS             NetworkID = 11
	NetworkIDTwitter         NetworkID = 12
	NetworkIDMisskey         NetworkID = 13
	NetworkIDJike            NetworkID = 14
	NetworkIDPlayStation     NetworkID = 15
	NetworkIDGitHub          NetworkID = 16
	NetworkIDZksync          NetworkID = 17

	NetworkSymbolUnknown                       = "unknown"
	NetworkSymbolEthereumMainnet NetworkSymbol = "ethereum_mainnet"
	NetworkSymbolPolygon         NetworkSymbol = "polygon"
	NetworkSymbolBNBChain        NetworkSymbol = "bnb"
	NetworkSymbolArbitrum        NetworkSymbol = "arbitrum"
	NetworkSymbolAvalanche       NetworkSymbol = "avalanche"
	NetworkSymbolFantom          NetworkSymbol = "fantom"
	NetworkSymbolGnosisMainnet   NetworkSymbol = "gnosis"
	NetworkSymbolSolanaMainet    NetworkSymbol = "solana_mainnet"
	NetworkSymbolFlowMainnet     NetworkSymbol = "flow_mainnet"
	NetworkSymbolArweaveMainnet  NetworkSymbol = "arweave_mainnet"
	NetworkSymbolRSS             NetworkSymbol = "rss"
	NetworkSymbolTwitter         NetworkSymbol = "twitter"
	NetworkSymbolMisskey         NetworkSymbol = "misskey"
	NetworkSymbolJike            NetworkSymbol = "jike"
	NetworkSymbolPlayStation     NetworkSymbol = "playstation"
	NetworkSymbolGitHub          NetworkSymbol = "github"
	NetworkSymbolZksync          NetworkSymbol = "zksync"
)

var (
	NetworkIDMap = map[NetworkSymbol]NetworkID{
		NetworkSymbolUnknown:         NetworkIDUnknown,
		NetworkSymbolEthereumMainnet: NetworkIDEthereumMainnet,
		NetworkSymbolPolygon:         NetworkIDPolygon,
		NetworkSymbolBNBChain:        NetworkIDBNBChain,
		NetworkSymbolArbitrum:        NetworkIDArbitrum,
		NetworkSymbolAvalanche:       NetworkIDAvalanche,
		NetworkSymbolFantom:          NetworkIDFantom,
		NetworkSymbolGnosisMainnet:   NetworkIDGnosisMainnet,
		NetworkSymbolSolanaMainet:    NetworkIDSolanaMainet,
		NetworkSymbolFlowMainnet:     NetworkIDFlowMainnet,
		NetworkSymbolArweaveMainnet:  NetworkIDArweaveMainnet,
		NetworkSymbolRSS:             NetworkIDRSS,
		NetworkSymbolTwitter:         NetworkIDTwitter,
		NetworkSymbolMisskey:         NetworkIDMisskey,
		NetworkSymbolJike:            NetworkIDJike,
		NetworkSymbolPlayStation:     NetworkIDPlayStation,
		NetworkSymbolGitHub:          NetworkIDGitHub,
	}
)

func IsValidNetworkName(value string) bool {
	id, has := NetworkIDMap[NetworkSymbol(value)]
	if has && id != NetworkIDUnknown {
		return true
	}

	return false
}

func (id NetworkSymbol) GetID() NetworkID {
	return NetworkIDMap[NetworkSymbol(id)]
}

func GetEthereumPlatformNetworks() []NetworkID {
	return []NetworkID{
		NetworkIDEthereumMainnet,
		NetworkIDPolygon,
		NetworkIDBNBChain,
		NetworkIDArbitrum,
		NetworkIDAvalanche,
		NetworkIDFantom,
		NetworkIDGnosisMainnet,
	}
}
