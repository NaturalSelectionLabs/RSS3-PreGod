package constants

type LinkSourceID int

func (p LinkSourceID) Int() int {
	return int(p)
}

func (p LinkSourceID) Name() LinkSourceName {
	if name, exist := linkSourceNameMap[p]; exist {
		return name
	}

	return LinkSourceNameUnknown
}

type LinkSourceName string

func (p LinkSourceName) ID() LinkSourceID {
	if id, exist := linkSourceIDMap[p]; exist {
		return id
	}

	return LinkSourceIDUnknown
}

func (p LinkSourceName) String() string {
	return string(p)
}

var (
	LinkSourceIDUnknown   LinkSourceID = -1
	LinkSourceIDCrossbell LinkSourceID = 0
	LinkSourceIDLens      LinkSourceID = 1

	LinkSourceNameUnknown   LinkSourceName = "Unknown"
	LinkSourceNameCrossbell LinkSourceName = "Crossbell"
	LinkSourceNameLens      LinkSourceName = "Lens"

	linkSourceNameMap = map[LinkSourceID]LinkSourceName{
		LinkSourceIDUnknown:   LinkSourceNameUnknown,
		LinkSourceIDCrossbell: LinkSourceNameCrossbell,
		LinkSourceIDLens:      LinkSourceNameLens,
	}
	linkSourceIDMap = map[LinkSourceName]LinkSourceID{}
)

func init() {
	for id, name := range linkSourceNameMap {
		linkSourceIDMap[name] = id
	}
}
