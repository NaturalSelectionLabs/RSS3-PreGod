package protocol

type BackLinkList struct {
	ListUnsignedBase

	Total int            `json:"total"`
	List  []LinkListItem `json:"list"`
}
