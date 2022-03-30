package protocol

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/datatype"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
)

type Index struct {
	SignedBase

	Profile IndexProfile `json:"profile"`
	Links   IndexLinks   `json:"links"`
	Items   IndexItems   `json:"items"`
}

func NewIndex(instance rss3uri.Instance) *Index {
	identifier := rss3uri.New(instance).String()

	return &Index{
		SignedBase: SignedBase{
			Base: Base{
				Version:    Version,
				Identifier: identifier,
			},
		},
		Links: IndexLinks{
			Identifiers:    []IndexLinkIdentifier{},
			IdentifierBack: fmt.Sprintf("%s/list/backlink", identifier),
		},
	}
}

type IndexProfile struct {
	Name        *string                 `json:"name"`
	Avatars     []string                `json:"avatars"`
	Bio         *string                 `json:"bio"`
	Attachments IndexProfileAttachments `json:"attachments"`
	Accounts    []IndexProfileAccount   `json:"accounts"`
}

type IndexProfileAttachments []IndexProfileAttachment

func (ipas *IndexProfileAttachments) ToDBStruct() datatype.Attachments {
	attachments := make([]datatype.Attachment, len(*ipas))
	for _, ipa := range *ipas {
		attachments = append(attachments, ipa.ToDBStruct())
	}

	return attachments
}

type IndexProfileAttachment struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	Address     string `json:"address"`
	MimeType    string `json:"mime_type"`
	SizeInBytes int    `json:"size_in_bytes"` // max: 4GB
}

func (ipa *IndexProfileAttachment) ToDBStruct() datatype.Attachment {
	return datatype.Attachment{
		Type:        ipa.Type,
		Content:     ipa.Content,
		Address:     ipa.Address,
		MimeType:    ipa.MimeType,
		SizeInBytes: ipa.SizeInBytes,
	}
}

type IndexProfileAccount struct {
	Identifier string `json:"identifier"`
	Signature  string `json:"signature,omitempty"`
}

func NewIndexProfileAccount(
	prefixID constants.PrefixID,
	platformIdentity string,
	platformID constants.PlatformID,
	signature string,
) *IndexProfileAccount {
	return &IndexProfileAccount{
		Identifier: rss3uri.New(&rss3uri.PlatformInstance{
			Prefix:   prefixID.String(),
			Identity: platformIdentity,
			Platform: platformID.Symbol(),
		}).String(),
		Signature: signature,
	}
}

func (i *Index) AddProfileAccount(
	platformAccountID string,
	platformID constants.PlatformID,
	signature string,
) {
	i.Profile.Accounts = append(i.Profile.Accounts, *NewIndexProfileAccount(constants.PrefixIDAccount, platformAccountID, platformID, signature))
}

type IndexLinks struct {
	Identifiers    []IndexLinkIdentifier `json:"identifiers"`
	IdentifierBack string                `json:"identifier_back"`
}

type IndexLinkIdentifier struct {
	Type             string `json:"type"`
	IdentifierCustom string `json:"identifier_custom"`
	Identifier       string `json:"identifier"`
}

func NewIndexLinkIdentifier(identifier string, linkType constants.LinkTypeID, maxPageIndex int) *IndexLinkIdentifier {
	return &IndexLinkIdentifier{
		Type:             linkType.String(),
		IdentifierCustom: fmt.Sprintf("%s/list/link/%s/%d", identifier, linkType, maxPageIndex),
		Identifier:       fmt.Sprintf("%s/list/link/%s", identifier, linkType),
	}
}

func (i *Index) AddLinkIdentifier(linkType constants.LinkTypeID, maxPageIndex int) {
	i.Links.Identifiers = append(i.Links.Identifiers, *NewIndexLinkIdentifier(i.Identifier, linkType, maxPageIndex))
}

type IndexItems struct {
	Notes  IndexItemsNotes  `json:"notes"`
	Assets IndexItemsAssets `json:"assets"`
}

type IndexItemsNotes struct {
	IdentifierCustom string                  `json:"identifier_custom"`
	Identifier       string                  `json:"identifier"`
	Filters          *IndexItemsNotesFilters `json:"filters"`
}

type IndexItemsNotesFilters struct {
	AllowList []string `json:"allowlist"`
	BlockList []string `json:"blocklist"`
}

type IndexItemsAssets struct {
	IdentifierCustom string                   `json:"identifier_custom"`
	Identifier       string                   `json:"identifier"`
	Filters          *IndexItemsAssetsFilters `json:"filters"`
}

type IndexItemsAssetsFilters struct {
	AllowList []string `json:"allowlist"`
	BlockList []string `json:"blocklist"`
}
