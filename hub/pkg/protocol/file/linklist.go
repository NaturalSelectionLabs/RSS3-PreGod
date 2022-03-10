package file

import "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"

type LinkList struct {
	protocol.SignedBase
	// TODO Define identifier_next
	Total int            `json:"total"`
	List  []LinkListItem `json:"list"`
}

type LinkListItem struct {
	IdentifierTarget string `json:"identifier_target"`
	Type             string `json:"type"`
}
