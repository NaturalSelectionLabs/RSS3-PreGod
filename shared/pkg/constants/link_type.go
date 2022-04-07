package constants

type LinkTypeID int

type LinkTypeName string

func (l LinkTypeID) Int() int {
	return int(l)
}

func (l LinkTypeID) Name() LinkTypeName {
	if name, exist := linkTypeNameMap[l]; exist {
		return name
	}

	return LinkTypeNameUnknown
}

// Converts LinkTypeID to string.
func (l LinkTypeID) String() string {
	return linkTypeNameMap[l].String()
}

func (l LinkTypeName) String() string {
	return string(l)
}

func (l LinkTypeName) ID() LinkTypeID {
	if id, exist := linkTypeIDMap[l]; exist {
		return id
	}

	return LinkTypeUnknown
}

const (
	LinkTypeUnknown LinkTypeID = -1

	LinkTypeFollow     LinkTypeID = 0
	LinkTypeComment    LinkTypeID = 1
	LinkTypeLike       LinkTypeID = 2
	LinkTypeCollection LinkTypeID = 3

	LinkTypeNameUnknown LinkTypeName = "unknown"

	LinkTypeNameFollow     LinkTypeName = "follow"
	LinkTypeNameComment    LinkTypeName = "comment"
	LinkTypeNameLike       LinkTypeName = "like"
	LinkTypeNameCollection LinkTypeName = "collection"
)

var LinkTypeMap = map[LinkTypeID]string{
	LinkTypeUnknown: "unknown",

	LinkTypeFollow:     "follow",
	LinkTypeComment:    "comment",
	LinkTypeLike:       "like",
	LinkTypeCollection: "collection",
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

var (
	linkTypeNameMap = map[LinkTypeID]LinkTypeName{
		LinkTypeUnknown: LinkTypeNameUnknown,

		LinkTypeFollow:     LinkTypeNameFollow,
		LinkTypeComment:    LinkTypeNameComment,
		LinkTypeLike:       LinkTypeNameLike,
		LinkTypeCollection: LinkTypeNameCollection,
	}
	linkTypeIDMap = map[LinkTypeName]LinkTypeID{}
)

func init() {
	for id, name := range linkTypeNameMap {
		linkTypeIDMap[name] = id
	}
}
