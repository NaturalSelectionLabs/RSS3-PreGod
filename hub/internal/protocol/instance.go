package protocol

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

const (
	InstanceProfiles  = "profiles"
	InstanceLinks     = "links"
	InstanceBacklinks = "backlinks"
	InstanceAssets    = "assets"
	InstanceNotes     = "notes"
)

type InstanceItem struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

func NewInstanceItemList(instance rss3uri.Instance) []InstanceItem {
	uri := rss3uri.New(instance).String()

	return []InstanceItem{
		{
			Type:       InstanceProfiles,
			Identifier: fmt.Sprintf("%s/profiles", uri),
		},
		{
			Type:       InstanceLinks,
			Identifier: fmt.Sprintf("%s/links", uri),
		},
		{
			Type:       InstanceBacklinks,
			Identifier: fmt.Sprintf("%s/backlinks", uri),
		},
		{
			Type:       InstanceAssets,
			Identifier: fmt.Sprintf("%s/assets", uri),
		},
		{
			Type:       InstanceNotes,
			Identifier: fmt.Sprintf("%s/notes", uri),
		},
	}
}
