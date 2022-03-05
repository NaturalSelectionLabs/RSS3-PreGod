package rss3_uri

import (
	"errors"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
)

type Authority struct {
	Prefix   constants.PrefixName   `json:"prefix"`
	Identity string                 `json:"identity"`
	Platform constants.PlatformName `json:"platform"`
}

// ParseAuthority parses the uri and returns the authority struct.
// Returns error if the uri is invalid.
func ParseAuthority(uri string) (*Authority, error) {
	s := strings.SplitN(uri, "@", 2)

	if len(s) != 2 {
		return nil, errors.New("no platform name")
	}

	pi := strings.SplitN(s[0], ":", 2)

	if len(pi) != 2 {
		return nil, errors.New("no prefix")
	}

	prefix := pi[0]
	identity := pi[1]
	platform := s[1]

	if !constants.IsValidPrefix(prefix) {
		return nil, errors.New("invalid prefix")
	}

	if !constants.IsValidPlatformName(platform) {
		return nil, errors.New("invalid platform name")
	}

	return &Authority{
		Prefix:   constants.PrefixName(prefix),
		Identity: identity,
		Platform: constants.PlatformName(platform),
	}, nil
}

// URI returns an RSS3 URI for any identity
func URI(identity string) string {
	return "rss3://" + identity
}

// AccountInstanceId returns an account instance id
// example:
// account:0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944@evm
// account:DIYgod@twitter
func AccountInstanceId(identity, accountPlatform string) string {
	return "account" + ":" + identity + "@" + accountPlatform
}

// ItemInstanceId returns an automatically generated Instance id for Item.
// example:
// asset:ethereum-0xacbe98efe2d4d103e221e04c76d7c55db15c8e89-5@evm
func ItemInstanceId(itemType, uniqueId, itemPlatform string) string {
	return itemType + ":" + uniqueId + "@" + itemPlatform
}

// AssetInstanceId returns an asset instance id
// example:
// asset:0xb9619cf4f875cdf0e3ce48b28a1c725bc4f6c0fb-1024@ethereum
func AssetInstanceId(assetAddress, tokenId, assetPlatform string) string {
	return "asset" + ":" + assetAddress + "-" + tokenId + "@" + assetPlatform
}

// NoteInstanceId returns an asset instance id
// example:
// note:5591079b-1f5b-4ae9-8209-51b18f0d3be0@twitter
func NoteInstanceId(noteUUID, itemPlatform string) string {
	return "note" + ":" + noteUUID + "@" + itemPlatform
}

func ItemId(instanceId, itemType, itemUUID string) string {
	return instanceId + "/" + itemType + "/" + itemUUID
}
