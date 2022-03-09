package constants

const (
	PrefixIDUnknown  PrefixID = 0
	PrefixIDAccount  PrefixID = 1
	PrefixIDInstance PrefixID = 2
	PrefixIDNote     PrefixID = 3
	PrefixIDAsset    PrefixID = 4

	PrefixNameUnknown  PrefixName = "unknown"
	PrefixNameAccount  PrefixName = "account"
	PrefixNameInstance PrefixName = "instance"
	PrefixNameNote     PrefixName = "note"
	PrefixNameAsset    PrefixName = "asset"
)

var (
	prefixNameMap = map[PrefixID]PrefixName{}
	prefixIDMap   = map[PrefixName]PrefixID{
		PrefixNameUnknown:  PrefixIDUnknown,
		PrefixNameAccount:  PrefixIDAccount,
		PrefixNameInstance: PrefixIDInstance,
		PrefixNameNote:     PrefixIDNote,
		PrefixNameAsset:    PrefixIDAsset,
	}
)

type PrefixName string
type PrefixID int32

func (id PrefixID) String() PrefixName {
	value, has := prefixNameMap[id]
	if has && value != PrefixNameUnknown {
		return value
	}

	return PrefixNameUnknown
}

func IsValidPrefix(value string) bool {
	id, has := prefixIDMap[PrefixName(value)]
	if has && id != PrefixIDUnknown {
		return true
	}

	return false
}

func init() {
	for name, id := range prefixIDMap {
		prefixNameMap[id] = name
	}
}
