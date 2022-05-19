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
	NetworkIDZkSync         NetworkID = 17

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
	NetworkSymbolZkSync         NetworkSymbol = "zksync"
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

	networkSymbolMap = map[NetworkID]NetworkSymbol{}
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
		NetworkIDCrossbell,
		NetworkIDEthereum,
		NetworkIDPolygon,
		NetworkIDBNBChain,
		NetworkIDArbitrum,
		NetworkIDAvalanche,
		NetworkIDFantom,
		NetworkIDGnosisMainnet,
	}
}

func GetNetworkList(platformID PlatformID) []NetworkID {
	switch platformID {
	case PlatformIDEthereum:
		return []NetworkID{
			NetworkIDCrossbell,     // 0
			NetworkIDEthereum,      // 1
			NetworkIDPolygon,       // 2
			NetworkIDBNBChain,      // 3
			NetworkIDArbitrum,      // 4
			NetworkIDArbitrum,      // 5
			NetworkIDFantom,        // 6
			NetworkIDGnosisMainnet, // 7
		}
	case PlatformIDSolana:
		return []NetworkID{
			NetworkIDSolanaMainet, // 8
		}
	case PlatformIDFlow:
		return []NetworkID{
			NetworkIDFlowMainnet, // 9
		}
	case PlatformIDArweave:
		return []NetworkID{
			NetworkIDArweaveMainnet, // 10
		}
	case PlatformIDRSS:
		return []NetworkID{
			NetworkIDRSS, // 11
		}
	case PlatformIDTwitter:
		return []NetworkID{
			NetworkIDTwitter, // 12
		}
	case PlatformIDMisskey:
		return []NetworkID{
			NetworkIDMisskey, // 13
		}
	case PlatformIDJike:
		return []NetworkID{
			NetworkIDJike, // 14
		}
	case PlatformIDPlayStation:
		return []NetworkID{
			NetworkIDPlayStation, // 15
		}
	case PlatformIDGitHub:
		return []NetworkID{
			NetworkIDGitHub, // 16
		}
	}

	return []NetworkID{}
}

func init() {
	for id, name := range networkIDMap {
		networkSymbolMap[name] = id
	}
}
