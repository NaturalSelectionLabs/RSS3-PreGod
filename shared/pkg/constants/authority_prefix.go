package constants

type PrefixName string
type PrefixID int32

const (
	PrefixID_Unknown PrefixID = 0

	PrefixID_Account  PrefixID = 1
	PrefixID_Instance PrefixID = 2
	PrefixID_Note     PrefixID = 3
	PrefixID_Asset    PrefixID = 4
)

const (
	PrefixName_Unknown PrefixName = "unknown"

	PrefixName_Account  PrefixName = "account"
	PrefixName_Instance PrefixName = "instance"
	PrefixName_Note     PrefixName = "note"
	PrefixName_Asset    PrefixName = "asset"
)

var PrefixMap = map[PrefixID]PrefixName{
	PrefixID_Unknown: PrefixName_Unknown,

	PrefixID_Account:  PrefixName_Account,
	PrefixID_Instance: PrefixName_Instance,
	PrefixID_Note:     PrefixName_Note,
	PrefixID_Asset:    PrefixName_Asset,
}

func IsValidPrefix(prefix string) bool {
	switch PrefixName(prefix) {
	case PrefixName_Account, PrefixName_Instance, PrefixName_Note, PrefixName_Asset:
		return true
	}

	return false
}

// Converts PrefixID to string.
func (id PrefixID) String() PrefixName {
	v, ok := PrefixMap[id]
	if !ok {
		return PrefixName_Unknown
	}

	return v
}

// Converts PrefixName to PrefixID.
func StringToPrefixID(prefix string) PrefixID {
	for k, v := range PrefixMap {
		if string(v) == prefix {
			return k
		}
	}

	return PrefixID_Unknown
}
