package constants

type AssetSourceID int

func (p AssetSourceID) Int() int {
	return int(p)
}

func (p AssetSourceID) Name() AssetSourceName {
	if name, exist := assetSourceNameMap[p]; exist {
		return name
	}

	return AssetSourceNameUnknown
}

type AssetSourceName string

func (p AssetSourceName) String() string {
	return string(p)
}

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
