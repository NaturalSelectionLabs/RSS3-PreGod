package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"

type BatchGetNodeListRequest struct {
	AddressList    []string `json:"address_list"`    // account address list
	Limit          int      `json:"limit"`           // amount of data per page
	LastIdentifier string   `json:"last_identifier"` // the last identifier of the previous page

	InstanceList []rss3uri.Instance `json:"-"` // parsed from address list
}
