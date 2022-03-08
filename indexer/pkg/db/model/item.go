package model

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/kamva/mgm/v3"
)

type ItemId struct {
	NetworkId constants.NetworkID `json:"network_id" bson:"network_id"`
	Proof     string              `json:"proof" bson:"proof"`
}
type Attachment struct {
	Content    string    `json:"content" bson:"content"`
	Address    []string  `json:"address" bson:"address"`
	MimeType   string    `json:"mime_type" bson:"mime_type"`
	Type       string    `json:"type" bson:"type"`
	SizeInByte int       `json:"size_in_bytes" bson:"size_in_bytes"`
	SyncAt     time.Time `json:"sync_at" bson:"sync_at"`
}

type Metadata map[string]interface{}

type Item struct {
	mgm.DefaultModel `bson:",inline"`

	ItemId            ItemId             `json:"item_id" bson:"item_id"` // Index
	Metadata          Metadata           `json:"metadata" bson:"metadata"`
	Tags              constants.ItemTags `json:"tags" bson:"tags"`
	Authors           []string           `json:"authors" bson:"authors"`
	Title             string             `json:"title" bson:"title"`
	Summary           string             `json:"summary" bson:"summary"`
	Attachments       []Attachment       `json:"attachments" bson:"attachments"`
	PlatformCreatedAt time.Time          `json:"date_created" bson:"date_created"`
}

func NewAttachment(content string, address []string, mimetype string, t string, size_in_bytes int, sync_at time.Time) *Attachment {
	return &Attachment{
		Content:    content,
		Address:    address,
		MimeType:   mimetype,
		Type:       t,
		SizeInByte: size_in_bytes,
		SyncAt:     sync_at,
	}
}

func NewItem(networkId constants.NetworkID, proof string, metadata Metadata,
	tags constants.ItemTags, authors []string, title string, summary string,
	attachments []Attachment, platformCreatedAt time.Time) *Item {
	return &Item{
		ItemId: ItemId{
			NetworkId: networkId,
			Proof:     proof,
		},
		Metadata:          metadata,
		Tags:              tags,
		Authors:           authors,
		Title:             title,
		Summary:           summary,
		Attachments:       attachments,
		PlatformCreatedAt: platformCreatedAt,
	}
}
