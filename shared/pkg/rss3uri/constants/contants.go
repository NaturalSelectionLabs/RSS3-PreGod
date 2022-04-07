package constants

import "fmt"

type ID interface {
	Int() int
	Name() Name
	Symbol() Symbol // returns name if no symbol defined in protocol
}

type Name interface {
	fmt.Stringer
	ID() ID
	Symbol() Symbol // returns name if no symbol defined in protocol
}

type Symbol interface {
	fmt.Stringer
	ID() ID
	Name() Name
}
