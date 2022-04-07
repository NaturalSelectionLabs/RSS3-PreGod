package nft_utils

import (
	"fmt"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	lop "github.com/samber/lo/parallel"
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
	if metadata == "" {
		return Metadata{}, fmt.Errorf("metadata is an empty string")
	}

	var parser fastjson.Parser
	v, err := parser.Parse(metadata)

	if err != nil {
		return Metadata{}, fmt.Errorf("got error: %v when parsing nft metadata: %+v", err, metadata)
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

func getContentHeader(uri string) *httpx.ContentHeader {
	// 1. is base64 encode?
	if strings.Contains(uri, ";base64,") {
		return &httpx.ContentHeader{
			MIMEType:   strings.Split(uri, ",")[0],
			SizeInByte: len([]byte(strings.Split(uri, ",")[1])),
		}
	}

	// 2. is ipfs?
	if strings.HasPrefix(uri, "ipfs://") {
		cid := strings.Split(uri, "ipfs://")[1]
		url := "https://cloudflare-ipfs.com/ipfs/" + cid
		contentHeader, _ := httpx.GetContentHeader(url)

		return contentHeader
	}

	// 3. is http?
	if strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://") {
		contentHeader, _ := httpx.GetContentHeader(uri)

		return contentHeader
	}

	return nil
}

func getCommAtt(meta Metadata) []datatype.Attachment {
	var as []datatype.Attachment

	if len(meta.ExternalLink) != 0 {
		as = append(as, datatype.Attachment{
			Type:    "external_url",
			Content: meta.ExternalLink,
		})
	}

	if len(meta.Preview) != 0 {
		as = append(as, datatype.Attachment{
			Type:    "preview",
			Address: meta.Preview,
		})
	}

	if len(meta.Object) != 0 {
		as = append(as, datatype.Attachment{
			Type:    "object",
			Address: meta.Object,
		})
	}

	if len(meta.Attributes) != 0 {
		as = append(as, datatype.Attachment{
			Type:     "attributes",
			Content:  meta.Attributes, //TODO: extract trait_type/value
			MimeType: "text/json",
		})
	}

	// TODO: make other unparsed values

	return as
}

// Convert metadata to attachment of asset.
func Meta2AssetAtt(meta Metadata) []datatype.Attachment {
	var att []datatype.Attachment

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

// Meta2NoteAtt converts metadata to attachment of note.
// Note that this function does NOT include the MimeType and SizeInBytes of the attachment.
// You may need to call CompleteMimeTypes to complete them later
// (for parallel requests).
func Meta2NoteAtt(meta Metadata) []datatype.Attachment {
	var att []datatype.Attachment

	att = append(att, getCommAtt(meta)...)

	return att
}

func CompleteMimeTypes(as []datatype.Attachment) {
	// get mimetypes
	lop.ForEach(as, func(a datatype.Attachment, i int) {
		if a.Address != "" {
			contentHeader := getContentHeader(a.Address)
			if contentHeader != nil {
				as[i].MimeType = contentHeader.MIMEType
				as[i].SizeInBytes = contentHeader.SizeInByte
			}
		}
	})
}
