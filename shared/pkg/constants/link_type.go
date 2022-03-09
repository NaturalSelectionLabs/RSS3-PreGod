package constants

type LinkTypeID int32

const (
	LinkTypeUnknown LinkTypeID = 0

	LinkTypeFollowing  LinkTypeID = 1
	LinkTypeComment    LinkTypeID = 2
	LinkTypeLike       LinkTypeID = 3
	LinkTypeCollection LinkTypeID = 4
)

var LinkTypeMap = map[LinkTypeID]string{
	LinkTypeUnknown: "unknown",

	LinkTypeFollowing:  "following",
	LinkTypeComment:    "comment",
	LinkTypeLike:       "like",
	LinkTypeCollection: "collection",
}

// Converts LinkTypeID to string.
func (id LinkTypeID) String() string {
	return LinkTypeMap[id]
}

// Converts string to LinkTypeID.
func StringToLinkTypeID(LinkType string) LinkTypeID {
	for k, v := range LinkTypeMap {
		if v == LinkType {
			return k
		}
	}

	return LinkTypeUnknown
}
