package file

import "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"

type BackLinkList struct {
	protocol.ListUnsignedBase

	Total int            `json:"total"`
	List  []LinkListItem `json:"list"`
}
