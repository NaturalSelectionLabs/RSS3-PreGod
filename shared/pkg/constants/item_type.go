package constants

type ItemType string
type ItemTypeID int

const (
	ItemTypeUnknown ItemType = "unknown"

	ItemType_Custom ItemTypeID = 1

	// AutoAssetType
	// 'evm_gitcoin_donation' | 'evm_xdai_poap' | 'evm_bsc_nft' | 'evm_ethereum_nft' | 'evm_polygon_nft'.
	ItemType_Gitcoin      ItemTypeID = 2
	ItemType_Xdai_Poap    ItemTypeID = 3
	ItemType_Bsc_Nft      ItemTypeID = 4
	ItemType_Ethereum_Nft ItemTypeID = 5
	ItemType_Polygon_Nft  ItemTypeID = 6
	ItemType_Avax_Nft     ItemTypeID = 7
	ItemType_Fantom_Nft   ItemTypeID = 8
	ItemType_Arbitrum_Nft ItemTypeID = 9

	// AutoNoteType
	// 'evm_mirror_entry' | 'twitter_tweet' | 'misskey_note' | 'jike_node'.
	ItemType_Evm_Mirror_Entry ItemTypeID = 9
	ItemType_Twitter_Tweet    ItemTypeID = 10
	ItemType_Misskey_Note     ItemTypeID = 11
	ItemType_Jike_Node        ItemTypeID = 12
)
