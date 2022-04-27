package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"

type BatchGetNodeListRequest struct {
	AddressList []string `json:"address_list"` // account address list
	Page        int      `json:"page"`         // which page to query, start from 1
	Limit       int      `json:"limit"`        // amount of data per page

	InstanceList []rss3uri.Instance `json:"-"` // parsed from address list
}
