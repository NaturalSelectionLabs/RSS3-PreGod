package protocol

type Instance struct {
	SignedBase

	Controller string `json:"controller,omitempty"` // A contract address indicating ownership of the file.

	Profile Profile `json:"profile,omitempty"`
	Links   Links   `json:"links,omitempty"`
	Items   Items   `json:"items,omitempty"`
}

type Profile struct {
	Name     string    `json:"name,omitempty"`
	Avatars  []string  `json:"avatars,omitempty"`
	Bio      string    `json:"bio,omitempty"`
	Banners  []string  `json:"banners,omitempty"`
	Websites []string  `json:"websites,omitempty"`
	Accounts []Account `json:"accounts,omitempty"`
}

type Links struct {
	Identifiers    []LinkIdentifier `json:"identifiers,omitempty"`
	IdentifierBack string           `json:"identifier_back,omitempty"`
}

type LinkIdentifier struct {
	Type             string `json:"type"`
	IdentifierCustom string `json:"identifier_custom"`
	Identifier       string `json:"identifier"`
}

type Items struct {
	Notes  Notes  `json:"notes,omitempty"`
	Assets Assets `json:"assets,omitempty"`
}

type Notes struct {
	Identifier     string `json:"identifier"`
	IdentifierPage string `json:"identifier_page,omitempty"`
}

type Assets struct {
	Identifier     string `json:"identifier"`
	IdentifierPage string `json:"identifier_page,omitempty"`
}

type Account struct {
	Identifier string `json:"identifier"`
	Signature  string `json:"signature,omitempty"` // Signature of `[RSS3] I am adding ${SignableAccount} to my RSS3 instance ${InstanceURI}`
}
