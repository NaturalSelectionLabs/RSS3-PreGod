package protocol

type IndexFile struct {
	SignedBase
	Agent   Agent   `json:"agent"`
	Profile Profile `json:"profile"`
	Links   Links   `json:"links"`
	Items   Items   `json:"items"`
}

type LinkListFile struct {
	SignedBase
	Total int                `json:"total"`
	List  []LinkListFileItem `json:"list,omitempty"`
}

type LinkListFileItem struct {
	Type             string `json:"type"`
	IdentifierTarget string `json:"identifier_target"`
}
