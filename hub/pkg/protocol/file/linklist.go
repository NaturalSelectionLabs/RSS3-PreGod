package file

import "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"

type LinkList struct {
	protocol.ListSignedBase

	Total int            `json:"total"`
	List  []LinkListItem `json:"list"`
}

type LinkListItem struct {
	IdentifierTarget string `json:"identifier_target"`
	Type             string `json:"type"`
}
