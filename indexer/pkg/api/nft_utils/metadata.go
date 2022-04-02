package nft_utils

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/valyala/fastjson"
)

type Metadata struct {
	Name         string
	Description  string
	ExternalLink string
	Attributes   string
	Object       string
	Preview      string
}

func ParseNFTMetadata(metadata string) (Metadata, error) {
	var parser fastjson.Parser
	v, err := parser.Parse(metadata)

	if err != nil {
		return Metadata{}, fmt.Errorf("got error: %v when parsing nft metadata: [%s]", err, metadata)
	}

	name := v.GetStringBytes("name")
	description := v.GetStringBytes("description")
	externalLink := v.GetStringBytes("external_link")
	preview := v.GetStringBytes("image")

	if len(preview) == 0 {
		preview = v.GetStringBytes("image_url")
	}

	object := v.GetStringBytes("animation_url")

	attributes := v.GetStringBytes("attributes")
	if len(attributes) == 0 {
		attributes = v.GetStringBytes("traits")
	}

	return Metadata{
		Name:         string(name),
		Description:  string(description),
		ExternalLink: string(externalLink),
		Preview:      string(preview),
		Object:       string(object),
		Attributes:   string(attributes),
	}, nil
}

func getCommAtt(meta Metadata) []datatype.Attachment {
	att := []datatype.Attachment{}

	if len(meta.ExternalLink) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "external_url",
			Content:  meta.ExternalLink,
			MimeType: "text/uri-list",
		})
	}

	if len(meta.Preview) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "preview",
			Content:  meta.Preview,
			MimeType: "", // TODO
		})
	}

	if len(meta.Object) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "object",
			Content:  meta.Object,
			MimeType: "", // TODO
		})
	}

	if len(meta.Attributes) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "attributes",
			Content:  meta.Attributes, //TODO: extract trait_type/value
			MimeType: "",              // TODO
		})
	}

	// TODO: make other unparsed values

	return att
}

// Convert metadata to attachment of asset
func Meta2AssetAtt(meta Metadata) []datatype.Attachment {
	att := []datatype.Attachment{}
	if len(meta.Name) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "name",
			Content:  meta.Name,
			MimeType: "text/plain",
		})
	}

	if len(meta.Description) != 0 {
		att = append(att, datatype.Attachment{
			Type:     "description",
			Content:  meta.Description,
			MimeType: "text/plain",
		})
	}

	att = append(att, getCommAtt(meta)...)

	return att
}

// Convert metadata to attachment of note
func Meta2NoteAtt(meta Metadata) []datatype.Attachment {
	att := []datatype.Attachment{}

	att = append(att, getCommAtt(meta)...)

	return att
}
