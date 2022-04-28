package nft_utils

import (
	"fmt"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/httpx"
	lop "github.com/samber/lo/parallel"
	"github.com/valyala/fastjson"
	"golang.org/x/sync/errgroup"
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
		return Metadata{}, nil
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

	attributesV := v.Get("attributes")
	if attributesV == nil {
		attributesV = v.Get("traits")
	}

	attributes := ""
	if attributesV != nil {
		attributes = attributesV.String()
	}

	return Metadata{
		Name:         string(name),
		Description:  string(description),
		ExternalLink: string(externalLink),
		Preview:      string(preview),
		Object:       string(object),
		Attributes:   attributes,
	}, nil
}

// Mainly used for formatting ipfs url to http url
func FormatUrl(url string) string {
	// 1. is data url?
	if strings.Contains(url, "data:") {
		return url
	}

	// 2. is ipfs?
	if strings.HasPrefix(url, "ipfs://") {
		cid := strings.Split(url, "ipfs://")[1]
		ret := "https://cloudflare-ipfs.com/ipfs/" + cid

		return ret
	}

	// TODO: need a smarter way to check if it is a ipfs gateway url
	if strings.Contains(url, "/ipfs/") {
		cid := strings.Split(url, "/ipfs/")[1]
		ret := "https://cloudflare-ipfs.com/ipfs/" + cid

		return ret
	}

	// 3. normal url
	return url
}

func getContentHeader(uri string) *httpx.ContentHeader {
	// 1. is base64 encode?
	if strings.Contains(uri, ";base64,") {
		return &httpx.ContentHeader{
			MIMEType:   strings.Split(uri, ",")[0],
			SizeInByte: len([]byte(strings.Split(uri, ",")[1])),
		}
	}

	// 2. format ipfs
	uri = FormatUrl(uri)

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
			Type:     "external_url",
			Content:  meta.ExternalLink,
			MimeType: "text/uri-list",
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

	if len(meta.Attributes) != 0 || meta.Attributes == `""` {
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

func CompleteMimeTypesForItems(notes []model.Note, assets []model.Asset, profiles []model.Profile) error {
	// complete attachments in parallel
	g := new(errgroup.Group)

	g.Go(func() error {
		lop.ForEach(notes, func(note model.Note, i int) {
			if note.Attachments != nil {
				as, err := database.UnwrapJSON[datatype.Attachments](note.Attachments)
				if err != nil {
					return
				}
				CompleteMimeTypes(as)
				notes[i].Attachments = database.MustWrapJSON(as)
			}
		})

		return nil
	})

	g.Go(func() error {
		lop.ForEach(assets, func(asset model.Asset, i int) {
			if asset.Attachments != nil {
				as, err := database.UnwrapJSON[datatype.Attachments](asset.Attachments)
				if err != nil {
					return
				}
				CompleteMimeTypes(as)
				assets[i].Attachments = database.MustWrapJSON(as)
			}
		})

		return nil
	})

	g.Go(func() error {
		lop.ForEach(profiles, func(profile model.Profile, i int) {
			if profile.Attachments != nil {
				as, err := database.UnwrapJSON[datatype.Attachments](profile.Attachments)
				if err != nil {
					return
				}
				CompleteMimeTypes(as)
				profiles[i].Attachments = database.MustWrapJSON(as)
			}
		})

		return nil
	})

	return g.Wait()
}
