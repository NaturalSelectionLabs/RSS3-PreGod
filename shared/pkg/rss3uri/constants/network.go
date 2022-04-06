package constants

var _ ID = NetworkID(0)

type NetworkID int

func (n NetworkID) Int() int {
	return int(n)
}

func (n NetworkID) Name() Name {
	return n.Symbol().Name()
}

func (n NetworkID) Symbol() Symbol {
	if value, ok := networkIDToSymbolMap[n]; ok {
		return value
	}

	return NetworkSymbolUnknown
}

var _ Symbol = NetworkSymbol("")

type NetworkSymbol string

func (n NetworkSymbol) String() string {
	return string(n)
}

func (n NetworkSymbol) ID() ID {
	if value, ok := networkSymbolToIDMap[n]; ok {
		return value
	}

	return NetworkIDUnknown
}

func (n NetworkSymbol) Name() Name {
	if value, ok := networkSymbolToNameMap[n]; ok {
		return value
	}

	return NetworkNameUnknown
}

var _ Name = NetworkName("")

type NetworkName string

func (n NetworkName) String() string {
	return string(n)
}

func (n NetworkName) ID() ID {
	return n.Symbol().ID()
}

func (n NetworkName) Symbol() Symbol {
	if value, ok := networkNameToSymbolMap[n]; ok {
		return value
	}

	return NetworkSymbolUnknown
}

var (
	NetworkIDUnknown     NetworkID = -1
	NetworkIDCrossbell   NetworkID = 0
	NetworkIDEthereum    NetworkID = 1
	NetworkIDPolygon     NetworkID = 2
	NetworkIDBNBChain    NetworkID = 3
	NetworkIDArbitrum    NetworkID = 4
	NetworkIDAvalanche   NetworkID = 5
	NetworkIDFantom      NetworkID = 6
	NetworkIDGnosis      NetworkID = 7
	NetworkIDSolana      NetworkID = 8
	NetworkIDFlow        NetworkID = 9
	NetworkIDArweave     NetworkID = 10
	NetworkIDRSS         NetworkID = 11
	NetworkIDTwitter     NetworkID = 12
	NetworkIDMisskey     NetworkID = 13
	NetworkIDJike        NetworkID = 14
	NetworkIDPlayStation NetworkID = 15
	NetworkIDGitHub      NetworkID = 16

	NetworkSymbolUnknown     NetworkSymbol = "unknown"
	NetworkSymbolCrossbell   NetworkSymbol = "crossbell"
	NetworkSymbolEthereum    NetworkSymbol = "ethereum"
	NetworkSymbolPolygon     NetworkSymbol = "polygon"
	NetworkSymbolBNBChain    NetworkSymbol = "bnb"
	NetworkSymbolArbitrum    NetworkSymbol = "arbitrum"
	NetworkSymbolAvalanche   NetworkSymbol = "avalanche"
	NetworkSymbolFantom      NetworkSymbol = "fantom"
	NetworkSymbolGnosis      NetworkSymbol = "gnosis"
	NetworkSymbolSolana      NetworkSymbol = "solana"
	NetworkSymbolFlowMainnet NetworkSymbol = "flow"
	NetworkSymbolArweave     NetworkSymbol = "arweave"
	NetworkSymbolRSS         NetworkSymbol = "rss"
	NetworkSymbolTwitter     NetworkSymbol = "twitter"
	NetworkSymbolMisskey     NetworkSymbol = "misskey"
	NetworkSymbolJike        NetworkSymbol = "jike"
	NetworkSymbolPlayStation NetworkSymbol = "playstation"
	NetworkSymbolGitHub      NetworkSymbol = "github"

	NetworkNameUnknown     NetworkName = "Unknown"
	NetworkNameCrossbell   NetworkName = "Crossbell"
	NetworkNameEthereum    NetworkName = "Ethereum"
	NetworkNamePolygon     NetworkName = "Polygon"
	NetworkNameBNBChain    NetworkName = "Binance Smart Chain"
	NetworkNameArbitrum    NetworkName = "Arbitrum"
	NetworkNameAvalanche   NetworkName = "Avalanche"
	NetworkNameFantom      NetworkName = "Fantom"
	NetworkNameGnosis      NetworkName = "Gnosis"
	NetworkNameSolana      NetworkName = "Solana"
	NetworkNameFlowMainnet NetworkName = "Flow"
	NetworkNameArweave     NetworkName = "Arweave"
	NetworkNameRSS         NetworkName = "RSS"
	NetworkNameTwitter     NetworkName = "Twitter"
	NetworkNameMisskey     NetworkName = "Misskey"
	NetworkNameJike        NetworkName = "Jike"
	NetworkNamePlayStation NetworkName = "PlayStation"
	NetworkNameGitHub      NetworkName = "GitHub"

	networkIDToSymbolMap = map[NetworkID]NetworkSymbol{
		NetworkIDUnknown:     NetworkSymbolUnknown,
		NetworkIDCrossbell:   NetworkSymbolCrossbell,
		NetworkIDEthereum:    NetworkSymbolEthereum,
		NetworkIDPolygon:     NetworkSymbolPolygon,
		NetworkIDBNBChain:    NetworkSymbolBNBChain,
		NetworkIDArbitrum:    NetworkSymbolArbitrum,
		NetworkIDAvalanche:   NetworkSymbolAvalanche,
		NetworkIDFantom:      NetworkSymbolFantom,
		NetworkIDGnosis:      NetworkSymbolGnosis,
		NetworkIDSolana:      NetworkSymbolSolana,
		NetworkIDFlow:        NetworkSymbolFlowMainnet,
		NetworkIDArweave:     NetworkSymbolArweave,
		NetworkIDRSS:         NetworkSymbolRSS,
		NetworkIDTwitter:     NetworkSymbolTwitter,
		NetworkIDMisskey:     NetworkSymbolMisskey,
		NetworkIDJike:        NetworkSymbolJike,
		NetworkIDPlayStation: NetworkSymbolPlayStation,
		NetworkIDGitHub:      NetworkSymbolGitHub,
	}
	networkSymbolToIDMap = map[NetworkSymbol]NetworkID{}

	networkSymbolToNameMap = map[NetworkSymbol]NetworkName{
		NetworkSymbolUnknown:     NetworkNameUnknown,
		NetworkSymbolCrossbell:   NetworkNameCrossbell,
		NetworkSymbolEthereum:    NetworkNameEthereum,
		NetworkSymbolPolygon:     NetworkNamePolygon,
		NetworkSymbolBNBChain:    NetworkNameBNBChain,
		NetworkSymbolArbitrum:    NetworkNameArbitrum,
		NetworkSymbolAvalanche:   NetworkNameAvalanche,
		NetworkSymbolFantom:      NetworkNameFantom,
		NetworkSymbolGnosis:      NetworkNameGnosis,
		NetworkSymbolSolana:      NetworkNameSolana,
		NetworkSymbolFlowMainnet: NetworkNameFlowMainnet,
		NetworkSymbolArweave:     NetworkNameArweave,
		NetworkSymbolRSS:         NetworkNameRSS,
		NetworkSymbolTwitter:     NetworkNameTwitter,
		NetworkSymbolMisskey:     NetworkNameMisskey,
		NetworkSymbolJike:        NetworkNameJike,
		NetworkSymbolPlayStation: NetworkNamePlayStation,
		NetworkSymbolGitHub:      NetworkNameGitHub,
	}
	networkNameToSymbolMap = map[NetworkName]NetworkSymbol{}
)

func init() {
	for id, symbol := range networkIDToSymbolMap {
		networkSymbolToIDMap[symbol] = id
	}

	for symbol, name := range networkSymbolToNameMap {
		networkNameToSymbolMap[name] = symbol
	}
}
