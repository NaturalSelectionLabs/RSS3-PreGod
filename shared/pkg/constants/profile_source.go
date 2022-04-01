package constants

type ProfileSourceID int

func (p ProfileSourceID) Int() int {
	return int(p)
}

func (p ProfileSourceID) Name() ProfileSourceName {
	if name, exist := profileSourceNameMap[p]; exist {
		return name
	}

	return ProfileSourceNameUnknown
}

type ProfileSourceName string

func (p ProfileSourceName) ID() ProfileSourceID {
	if id, exist := profileSourceIDMap[p]; exist {
		return id
	}

	return ProfileSourceIDUnknown
}

func (p ProfileSourceName) String() string {
	return string(p)
}

var (
	ProfileSourceIDUnknown   ProfileSourceID = -1
	ProfileSourceIDCrossbell ProfileSourceID = 0
	ProfileSourceIDENS       ProfileSourceID = 1
	ProfileSourceIDLens      ProfileSourceID = 2

	ProfileSourceNameUnknown   ProfileSourceName = "unknown"
	ProfileSourceNameCrossbell ProfileSourceName = "Crossbell"
	ProfileSourceNameENS       ProfileSourceName = "ENS"
	ProfileSourceNameLens      ProfileSourceName = "Lens"

	profileSourceNameMap = map[ProfileSourceID]ProfileSourceName{
		ProfileSourceIDUnknown:   ProfileSourceNameUnknown,
		ProfileSourceIDCrossbell: ProfileSourceNameCrossbell,
		ProfileSourceIDENS:       ProfileSourceNameENS,
		ProfileSourceIDLens:      ProfileSourceNameLens,
	}
	profileSourceIDMap = map[ProfileSourceName]ProfileSourceID{}
)

func init() {
	for id, name := range profileSourceNameMap {
		profileSourceIDMap[name] = id
	}
}
