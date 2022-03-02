package constants

type NetworkNameID int32
type NetworkName string

const (
	NetworkNameID_Unknown NetworkNameID = 0

	NetworkNameID_Ethereum  NetworkNameID = 1
	NetworkNameID_Polygon   NetworkNameID = 2
	NetworkNameID_BNB       NetworkNameID = 3
	NetworkNameID_Arbitrum  NetworkNameID = 4
	NetworkNameID_Avalanche NetworkNameID = 5
	NetworkNameID_Fantom    NetworkNameID = 6
	NetworkNameID_Gnosis    NetworkNameID = 7
	NetworkNameID_Solana    NetworkNameID = 8
	NetworkNameID_Flow      NetworkNameID = 9
	NetworkNameID_Arweave   NetworkNameID = 10

	NetworkNameID_Twitter NetworkNameID = 11
	NetworkNameID_Misskey NetworkNameID = 12
	NetworkNameID_Jike    NetworkNameID = 13
)

const (
	NetworkName_Unknown NetworkName = "Unknown"

	NetworkName_Ethereum  NetworkName = "Ethereum Mainnet"
	NetworkName_Polygon   NetworkName = "Polygon"
	NetworkName_BNB       NetworkName = "BNB Chain"
	NetworkName_Arbitrum  NetworkName = "Arbitrum"
	NetworkName_Avalanche NetworkName = "Avalanche"
	NetworkName_Fantom    NetworkName = "Fantom"
	NetworkName_Gnosis    NetworkName = "Gnosis Mainnet"
	NetworkName_Solana    NetworkName = "Solana Mainnet"
	NetworkName_Flow      NetworkName = "Flow Mainnet"
	NetworkName_Arweave   NetworkName = "Arweave Mainnet"

	NetworkName_Twitter NetworkName = "Twitter"
	NetworkName_Misskey NetworkName = "Misskey"
	NetworkName_Jike    NetworkName = "Jike"
)
