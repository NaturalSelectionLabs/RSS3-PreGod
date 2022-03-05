package model

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/kamva/mgm/v3"
)

type Attachment struct {
	Content    string   `json:"content" bson:"content"`
	Address    []string `json:"address" bson:"address"`
	MimeType   string   `json:"mime_type" bson:"mime_type"`
	Type       string   `json:"type" bson:"type"`
	SizeInByte int      `json:"size_in_bytes" bson:"size_in_bytes"`
}

type Object struct {
	mgm.DefaultModel `bson:",inline"`

	Uid        string               `json:"uid" bson:"uid"` // Index: (Uid, ItemType)
	ItemTypeID constants.ItemTypeID `json:"item_type_id" bson:"item_type_id"`

	Authors     []string     `json:"authors" bson:"authors"`
	Title       string       `json:"title" bson:"title"`
	Summary     string       `json:"summary" bson:"summary"`
	Tags        []string     `json:"tags" bson:"tags"`
	Attachments []Attachment `json:"attachments" bson:"attachments"`
}

func NewAttachment(content string, address []string, mimetype string, t string, size_in_bytes int) *Attachment {
	return &Attachment{
		Content:    content,
		Address:    address,
		MimeType:   mimetype,
		Type:       t,
		SizeInByte: size_in_bytes,
	}
}

func NewObject(
	authors []string,
	uid string,
	itemTypeID constants.ItemTypeID,
	title string,
	summary string,
	tags []string,
	attachments []Attachment,
) *Object {
	return &Object{
		Authors:     authors,
		Uid:         uid,
		ItemTypeID:  itemTypeID,
		Title:       title,
		Summary:     summary,
		Tags:        tags,
		Attachments: attachments,
	}
}
