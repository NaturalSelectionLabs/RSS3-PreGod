package constants

type (
	NetworkID     int32
	NetworkSymbol string
)

func (p NetworkID) Symbol() NetworkSymbol {
	if v, ok := networkSymbolMap[p]; ok {
		return v
	}

	return NetworkSymbolUnknown
}

func (p NetworkSymbol) ID() NetworkID {
	if v, ok := networkIDMap[p]; ok {
		return v
	}

	return NetworkIDUnknown
}

func (p NetworkSymbol) String() string {
	return string(p)
}

const (
	NetworkIDUnknown        NetworkID = -1
	NetworkIDCrossbell      NetworkID = 0
	NetworkIDEthereum       NetworkID = 1
	NetworkIDPolygon        NetworkID = 2
	NetworkIDBNBChain       NetworkID = 3
	NetworkIDArbitrum       NetworkID = 4
	NetworkIDAvalanche      NetworkID = 5
	NetworkIDFantom         NetworkID = 6
	NetworkIDGnosisMainnet  NetworkID = 7
	NetworkIDSolanaMainet   NetworkID = 8
	NetworkIDFlowMainnet    NetworkID = 9
	NetworkIDArweaveMainnet NetworkID = 10
	NetworkIDRSS            NetworkID = 11
	NetworkIDTwitter        NetworkID = 12
	NetworkIDMisskey        NetworkID = 13
	NetworkIDJike           NetworkID = 14
	NetworkIDPlayStation    NetworkID = 15
	NetworkIDGitHub         NetworkID = 16
	NetworkIDZksync         NetworkID = 17

	NetworkSymbolUnknown        NetworkSymbol = "unknown"
	NetworkSymbolCrossbell      NetworkSymbol = "crossbell"
	NetworkSymbolEthereum       NetworkSymbol = "ethereum"
	NetworkSymbolPolygon        NetworkSymbol = "polygon"
	NetworkSymbolBNBChain       NetworkSymbol = "bnb"
	NetworkSymbolArbitrum       NetworkSymbol = "arbitrum"
	NetworkSymbolAvalanche      NetworkSymbol = "avalanche"
	NetworkSymbolFantom         NetworkSymbol = "fantom"
	NetworkSymbolGnosisMainnet  NetworkSymbol = "gnosis"
	NetworkSymbolSolanaMainet   NetworkSymbol = "solana_mainnet"
	NetworkSymbolFlowMainnet    NetworkSymbol = "flow_mainnet"
	NetworkSymbolArweaveMainnet NetworkSymbol = "arweave_mainnet"
	NetworkSymbolRSS            NetworkSymbol = "rss"
	NetworkSymbolTwitter        NetworkSymbol = "twitter"
	NetworkSymbolMisskey        NetworkSymbol = "misskey"
	NetworkSymbolJike           NetworkSymbol = "jike"
	NetworkSymbolPlayStation    NetworkSymbol = "playstation"
	NetworkSymbolGitHub         NetworkSymbol = "github"
	NetworkSymbolZksync         NetworkSymbol = "zksync"
)

var (
	networkIDMap = map[NetworkSymbol]NetworkID{
		NetworkSymbolUnknown:        NetworkIDUnknown,
		NetworkSymbolCrossbell:      NetworkIDCrossbell,
		NetworkSymbolEthereum:       NetworkIDEthereum,
		NetworkSymbolPolygon:        NetworkIDPolygon,
		NetworkSymbolBNBChain:       NetworkIDBNBChain,
		NetworkSymbolArbitrum:       NetworkIDArbitrum,
		NetworkSymbolAvalanche:      NetworkIDAvalanche,
		NetworkSymbolFantom:         NetworkIDFantom,
		NetworkSymbolGnosisMainnet:  NetworkIDGnosisMainnet,
		NetworkSymbolSolanaMainet:   NetworkIDSolanaMainet,
		NetworkSymbolFlowMainnet:    NetworkIDFlowMainnet,
		NetworkSymbolArweaveMainnet: NetworkIDArweaveMainnet,
		NetworkSymbolRSS:            NetworkIDRSS,
		NetworkSymbolTwitter:        NetworkIDTwitter,
		NetworkSymbolMisskey:        NetworkIDMisskey,
		NetworkSymbolJike:           NetworkIDJike,
		NetworkSymbolPlayStation:    NetworkIDPlayStation,
		NetworkSymbolGitHub:         NetworkIDGitHub,
	}

	networkSymbolMap = map[NetworkID]NetworkSymbol{
		NetworkIDUnknown:        NetworkSymbolUnknown,
		NetworkIDCrossbell:      NetworkSymbolCrossbell,
		NetworkIDPolygon:        NetworkSymbolPolygon,
		NetworkIDBNBChain:       NetworkSymbolBNBChain,
		NetworkIDArbitrum:       NetworkSymbolArbitrum,
		NetworkIDAvalanche:      NetworkSymbolAvalanche,
		NetworkIDFantom:         NetworkSymbolFantom,
		NetworkIDGnosisMainnet:  NetworkSymbolGnosisMainnet,
		NetworkIDSolanaMainet:   NetworkSymbolSolanaMainet,
		NetworkIDFlowMainnet:    NetworkSymbolFlowMainnet,
		NetworkIDArweaveMainnet: NetworkSymbolArweaveMainnet,
		NetworkIDRSS:            NetworkSymbolRSS,
		NetworkIDTwitter:        NetworkSymbolTwitter,
		NetworkIDMisskey:        NetworkSymbolMisskey,
		NetworkIDJike:           NetworkSymbolJike,
		NetworkIDPlayStation:    NetworkSymbolPlayStation,
		NetworkIDGitHub:         NetworkSymbolGitHub,
	}
)

func IsValidNetworkName(value string) bool {
	id, has := networkIDMap[NetworkSymbol(value)]
	if has && id != NetworkIDUnknown {
		return true
	}

	return false
}

func (id NetworkSymbol) GetID() NetworkID {
	return networkIDMap[NetworkSymbol(id)]
}

func GetEthereumPlatformNetworks() []NetworkID {
	return []NetworkID{
		NetworkIDEthereum,
		NetworkIDPolygon,
		NetworkIDBNBChain,
		NetworkIDArbitrum,
		NetworkIDAvalanche,
		NetworkIDFantom,
		NetworkIDGnosisMainnet,
	}
}
