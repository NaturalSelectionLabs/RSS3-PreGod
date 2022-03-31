package constants

type ProfileSourceID int

func (p ProfileSourceID) Int() int {
	return int(p)
}

func (p ProfileSourceID) Name() ProfileSourceName {
	return nameMap[p]
}

type ProfileSourceName string

func (p ProfileSourceName) String() string {
	return string(p)
}

var (
	ProfileSourceIDCrossbell ProfileSourceID = 0
	ProfileSourceIDENS       ProfileSourceID = 1
	ProfileSourceIDLens      ProfileSourceID = 2

	ProfileSourceNameCrossbell ProfileSourceName = "Crossbell"
	ProfileSourceNameENS       ProfileSourceName = "ENS"
	ProfileSourceNameLens      ProfileSourceName = "Lens"

	nameMap = map[ProfileSourceID]ProfileSourceName{
		ProfileSourceIDCrossbell: ProfileSourceNameCrossbell,
		ProfileSourceIDENS:       ProfileSourceNameENS,
		ProfileSourceIDLens:      ProfileSourceNameLens,
	}
	idMap = map[ProfileSourceName]ProfileSourceID{}
)

func init() {
	for id, name := range nameMap {
		idMap[name] = id
	}
}
