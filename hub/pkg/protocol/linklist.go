package protocol

type LinkList struct {
	ListSignedBase

	Total int            `json:"total"`
	List  []LinkListItem `json:"list"`
}

type LinkListItem struct {
	IdentifierTarget string `json:"identifier_target"`
	Type             string `json:"type"`
}
