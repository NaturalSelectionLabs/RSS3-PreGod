package moralis

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type MoralisChainType string

const (
	Unknown MoralisChainType = "unknown"

	ETH     MoralisChainType = "eth"
	BSC     MoralisChainType = "bsc"
	Polygon MoralisChainType = "polygon"
	AVAX    MoralisChainType = "avalanche"
	Fantom  MoralisChainType = "fantom"
)

func GetChainType(network constants.NetworkName) MoralisChainType {
	if network == constants.NetworkName_Ethereum {
		return ETH
	} else if network == constants.NetworkName_BNB {
		return BSC
	} else if network == constants.NetworkName_Polygon {
		return Polygon
	} else if network == constants.NetworkName_Avalanche {
		return AVAX
	} else if network == constants.NetworkName_Fantom {
		return Fantom
	} else {
		return Unknown
	}
}

func (mt MoralisChainType) GetNFTItemTypeID() constants.ItemTypeID {
	if mt == ETH {
		return constants.ItemType_Ethereum_Nft
	} else if mt == BSC {
		return constants.ItemType_Bsc_Nft
	} else if mt == Polygon {
		return constants.ItemType_Polygon_Nft
	} else if mt == AVAX {
		return constants.ItemType_Avax_Nft
	} else if mt == Fantom {
		return constants.ItemType_Fantom_Nft
	} else {
		return constants.ItemType_Unknown
	}
}
