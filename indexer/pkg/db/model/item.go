package model

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attachment struct {
	Content    string    `json:"content" bson:"content"`
	Address    []string  `json:"address" bson:"address"`
	MimeType   string    `json:"mime_type" bson:"mime_type"`
	Type       string    `json:"type" bson:"type"`
	SizeInByte int       `json:"size_in_bytes" bson:"size_in_bytes"`
	SyncAt     time.Time `json:"sync_at" bson:"sync_at"`
}

type Item struct {
	mgm.DefaultModel `bson:",inline"`

	ItemId            ObjectId           `json:"item_id" bson:"item_id"` // Index
	Metadata          Metadata           `json:"metadata" bson:"metadata"`
	Tags              constants.ItemTags `json:"tags" bson:"tags"`
	Authors           []string           `json:"authors" bson:"authors"`
	Title             string             `json:"title,omitempty" bson:"title"`
	Summary           string             `json:"summary" bson:"summary"`
	Attachments       []Attachment       `json:"attachments" bson:"attachments"`
	PlatformCreatedAt time.Time          `json:"date_created" bson:"date_created"`
}

func NewAttachment(
	content string,
	address []string,
	mimetype string,
	t string,
	size_in_bytes int,
	sync_at time.Time) *Attachment {
	return &Attachment{
		Content:    content,
		Address:    address,
		MimeType:   mimetype,
		Type:       t,
		SizeInByte: size_in_bytes,
		SyncAt:     sync_at,
	}
}

func NewItem(
	networkId constants.NetworkID,
	proof string,
	metadata Metadata,
	tags constants.ItemTags,
	authors []string,
	title string,
	summary string,
	attachments []Attachment,
	platformCreatedAt time.Time,
) *Item {
	return &Item{
		DefaultModel: mgm.DefaultModel{
			IDField: mgm.IDField{
				ID: primitive.NewObjectID(),
			},
			DateFields: mgm.DateFields{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		ItemId: ObjectId{
			NetworkID: networkId,
			Proof:     proof,
			SyncAt:    time.Now(),
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
