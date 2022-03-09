package constants

type InstanceTypeID int32

const (
	InstanceTypeUnknown InstanceTypeID = 0

	InstanceTypeAccount InstanceTypeID = 1
	InstanceTypeObject  InstanceTypeID = 2
	InstanceTypeAsset   InstanceTypeID = 3
	InstanceTypeNote    InstanceTypeID = 4
)

var InstanceTypeMap = map[InstanceTypeID]string{
	InstanceTypeUnknown: "unknown",

	InstanceTypeAccount: "account",
	InstanceTypeObject:  "object",
	InstanceTypeAsset:   "asset",
	InstanceTypeNote:    "note",
}

// Converts InstanceTypeID to string.
func (id InstanceTypeID) String() string {
	return InstanceTypeMap[id]
}

// Converts string to InstanceTypeID.
func StringToInstanceTypeID(instanceType string) InstanceTypeID {
	for k, v := range InstanceTypeMap {
		if v == instanceType {
			return k
		}
	}

	return InstanceTypeUnknown
}
