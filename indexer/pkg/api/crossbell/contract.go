package crossbell

import (
	"errors"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/cryptox"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ContractAddressProfile = "0xDAFf56464108DA447970FBB2D3cf911991CeC034"

	// MethodCreateProfile bd5f69cb929c5721ec11cf16c853c0ed5b3b6ae00c9b2d1482b5f52838537df3
	MethodCreateProfile = cryptox.Keccak256([]byte("createProfile((address,string,string,address,bytes))"))
	// MethodSetPrimaryProfile 295cb43e0f15a50d8c6b90c94ed246fde29af3e5a7f46fa7ba43a5338c39ccf9
	MethodSetPrimaryProfile = cryptox.Keccak256([]byte("setPrimaryProfileId(uint256)"))

	// EventTransfer ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
	EventTransfer = cryptox.Keccak256([]byte("Transfer(address,address,uint256)"))
	// EventProfileCreated a5802a04162552328d75eaac538a033704a7c3beab65d0a83e52da1c8c9b7cdf
	EventProfileCreated = cryptox.Keccak256([]byte("ProfileCreated(uint256,address,address,string,uint256)"))
	// EventLinkProfile bc914995d574dd9ef2df364e4eee2b85deda3ba35d054a62425fba1b97275716
	EventLinkProfile = cryptox.Keccak256([]byte("LinkProfile(address,uint256,uint256,bytes32,uint256)"))
)

var (
	ErrInvalidData = errors.New("invalid data")
)

const (
	lenValue = 32
)

func ParseData(data []byte) (topics []common.Hash, err error) {
	if len(data)%lenValue != 0 {
		return nil, ErrInvalidData
	}

	size := len(data) / lenValue

	topics = make([]common.Hash, size)
	for i := 0; i < size; i++ {
		topics[i] = common.BytesToHash(data[i*lenValue : i*lenValue+lenValue])
	}

	return topics, nil
}
