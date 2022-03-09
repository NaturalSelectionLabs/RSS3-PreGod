package constants

type PlatformID int

func (p PlatformID) Symbol() PlatformSymbol {
	if v, ok := platformSymbolMap[p]; ok {
		return v
	}

	return PlatformSymbolUnknown
}

func (p PlatformID) IsSignable() bool {
	_, ok := signablePlatformSymbolMap[p]

	return ok
}

type PlatformSymbol string

func (p PlatformSymbol) ID() PlatformID {
	if v, ok := platformIDMap[p]; ok {
		return v
	}

	return PlatformIDUnknown
}

const (
	PlatformIDUnknown     PlatformID = 0
	PlatformIDEthereum    PlatformID = 1
	PlatformIDSolana      PlatformID = 2
	PlatformIDFlow        PlatformID = 3
	PlatformIDArweave     PlatformID = 4
	PlatformIDRSS         PlatformID = 5
	PlatformIDTwitter     PlatformID = 6
	PlatformIDMisskey     PlatformID = 7
	PlatformIDJike        PlatformID = 8
	PlatformIDPlayStation PlatformID = 9
	PlatformIDGitHub      PlatformID = 10

	PlatformSymbolUnknown     PlatformSymbol = "unknown"
	PlatformSymbolEthereum    PlatformSymbol = "ethereum"
	PlatformSymbolSolana      PlatformSymbol = "solana"
	PlatformSymbolFlow        PlatformSymbol = "flow"
	PlatformSymbolArweave     PlatformSymbol = "arweave"
	PlatformSymbolRSS         PlatformSymbol = "rss"
	PlatformSymbolTwitter     PlatformSymbol = "twitter"
	PlatformSymbolMisskey     PlatformSymbol = "misskey"
	PlatformSymbolJike        PlatformSymbol = "jike"
	PlatformSymbolPlayStation PlatformSymbol = "playstation"
	PlatformSymbolGitHub      PlatformSymbol = "github"
)

var (
	platformSymbolMap = map[PlatformID]PlatformSymbol{
		PlatformIDUnknown:     PlatformSymbolUnknown,
		PlatformIDEthereum:    PlatformSymbolEthereum,
		PlatformIDSolana:      PlatformSymbolSolana,
		PlatformIDFlow:        PlatformSymbolFlow,
		PlatformIDArweave:     PlatformSymbolArweave,
		PlatformIDRSS:         PlatformSymbolRSS,
		PlatformIDTwitter:     PlatformSymbolTwitter,
		PlatformIDMisskey:     PlatformSymbolMisskey,
		PlatformIDJike:        PlatformSymbolJike,
		PlatformIDPlayStation: PlatformSymbolPlayStation,
		PlatformIDGitHub:      PlatformSymbolGitHub,
	}
	platformIDMap = map[PlatformSymbol]PlatformID{}

	signablePlatformSymbolMap = map[PlatformID]PlatformSymbol{
		PlatformIDEthereum: PlatformSymbolEthereum,
		PlatformIDSolana:   PlatformSymbolSolana,
		PlatformIDFlow:     PlatformSymbolFlow,
	}
)

func IsValidPlatformSymbol(value string) bool {
	id, has := platformIDMap[PlatformSymbol(value)]
	if has && id != PlatformIDUnknown {
		return true
	}

	return false
}

func init() {
	for id, symbol := range platformSymbolMap {
		platformIDMap[symbol] = id
	}
}

func (p PlatformID) GetNetwork() []NetworkID {
	switch p {
	case PlatformIDEthereum:
		return append(GetEthereumPlatformNetworks(), NetworkIDArweaveMainnet)
	case PlatformIDArweave:
		return []NetworkID{NetworkIDArweaveMainnet}
	case PlatformIDTwitter:
		return []NetworkID{NetworkIDTwitter}
	case PlatformIDMisskey:
		return []NetworkID{NetworkIDMisskey}
	case PlatformIDJike:
		return []NetworkID{NetworkIDJike}
	default:
		return nil
	}
}
