package protocol

// These need to be changed to generics after go 1.18 is released.

type ListSignedBase struct {
	SignedBase

	IdentifierNext string `json:"identifier_next"`
}

type ListUnsignedBase struct {
	UnsignedBase

	IdentifierNext string `json:"identifier_next"`
}

type ItemPageList struct {
	ListSignedBase

	List []Item `json:"list"`
}

type ItemList struct {
	ListUnsignedBase

	List []Item `json:"list"`
}

type LinkList struct {
	ListSignedBase

	List []string `json:"list"`
}

type BacklinkList struct {
	ListUnsignedBase

	List []string `json:"list"`
}
