package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"

type BatchGetNodeListRequest struct {
	AddressList    []string `json:"addresses"`       // account address list
	Limit          int      `json:"limit"`           // amount of data per page
	LastIdentifier string   `json:"last_identifier"` // the last identifier of the previous page
	Tags           []string `json:"tags"`
	ExcludeTags    []string `json:"exclude_tags"`
	ItemSources    []string `json:"item_sources"`
	Networks       []string `json:"networks"`
	Latest         bool     `json:"latest"`

	InstanceList []rss3uri.Instance `json:"-"` // parsed from address list
}
