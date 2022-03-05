package constants

type ItemTypeID int32

const (
	ItemType_Unknown ItemTypeID = 0

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

	// AutoNoteType
	// 'evm_mirror_entry' | 'twitter_tweet' | 'misskey_note' | 'jike_node'.
	ItemType_Evm_Mirror_Entry ItemTypeID = 9
	ItemType_Twitter_Tweet    ItemTypeID = 10
	ItemType_Misskey_Note     ItemTypeID = 11
	ItemType_Jike_Node        ItemTypeID = 12
)

var ItemTypeMap = map[ItemTypeID]string{
	ItemType_Unknown: "unknown",

	ItemType_Custom:       "custom",
	ItemType_Gitcoin:      "evm_gitcoin_donation",
	ItemType_Xdai_Poap:    "evm_xdai_poap",
	ItemType_Bsc_Nft:      "evm_bsc_nft",
	ItemType_Ethereum_Nft: "evm_ethereum_nft",
	ItemType_Polygon_Nft:  "evm_polygon_nft",
	ItemType_Avax_Nft:     "evm_avax_nft",
	ItemType_Fantom_Nft:   "evm_fantom_nft",

	ItemType_Evm_Mirror_Entry: "evm_mirror_entry",
	ItemType_Twitter_Tweet:    "twitter_tweet",
	ItemType_Misskey_Note:     "misskey_note",
	ItemType_Jike_Node:        "jike_node",
}

// Converts ItemTypeID to string.
func (id ItemTypeID) String() string {
	return ItemTypeMap[id]
}

// Converts string to ItemTypeID.
func StringToItemTypeID(ItemType string) ItemTypeID {
	for k, v := range ItemTypeMap {
		if v == ItemType {
			return k
		}
	}

	return ItemType_Unknown
}
