package constants

// ID

type AssetSourceID int

var _ ID = AssetSourceID(0)

func (a AssetSourceID) Int() int {
	return int(a)
}

func (a AssetSourceID) Name() Name {
	return a.Symbol().Name()
}

func (a AssetSourceID) Symbol() Symbol {
	if value, ok := assetSourceNameMap[a]; ok {
		return value
	}

	return AssetSourceNameUnknown
}

// Symbol

//type AssetSourceSymbol string

//var _ Symbol = AssetSourceSymbol("")

// Name

// map

var (
	AssetSourceIDUnknown           AssetSourceID = -1
	AssetSourceIDCrossbell         AssetSourceID = 0
	AssetSourceIDEthereumNFT       AssetSourceID = 1
	AssetSourceIDSolanaNFT         AssetSourceID = 2
	AssetSourceIDFlowNFT           AssetSourceID = 3
	AssetSourceIDPlayStationTrophy AssetSourceID = 4
	AssetSourceIDGitHubAchievement AssetSourceID = 5

	AssetSourceNameUnknown           AssetSourceName = "Unknown"
	AssetSourceNameCrossbell         AssetSourceName = "Crossbell"
	AssetSourceNameEthereumNFT       AssetSourceName = "Ethereum NFT"
	AssetSourceNameSolanaNFT         AssetSourceName = "Solana NFT"
	AssetSourceNameFlowNFT           AssetSourceName = "Flow NFT"
	AssetSourceNamePlayStationTrophy AssetSourceName = "PlayStation Trophy"
	AssetSourceNameGitHubAchievement AssetSourceName = "GitHub Achievement"

	assetSourceNameMap = map[AssetSourceID]AssetSourceName{
		AssetSourceIDUnknown:           AssetSourceNameUnknown,
		AssetSourceIDCrossbell:         AssetSourceNameCrossbell,
		AssetSourceIDEthereumNFT:       AssetSourceNameEthereumNFT,
		AssetSourceIDSolanaNFT:         AssetSourceNameSolanaNFT,
		AssetSourceIDFlowNFT:           AssetSourceNameFlowNFT,
		AssetSourceIDPlayStationTrophy: AssetSourceNamePlayStationTrophy,
		AssetSourceIDGitHubAchievement: AssetSourceNameGitHubAchievement,
	}
	assetSourceIDMap = map[AssetSourceName]AssetSourceID{}
)

func init() {
	for id, name := range assetSourceNameMap {
		assetSourceIDMap[name] = id
	}
}
