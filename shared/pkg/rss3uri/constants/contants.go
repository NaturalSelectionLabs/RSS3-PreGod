package constants

import "fmt"

type ID interface {
	Int() int
	Name() Name
	Symbol() Symbol
}

type Name interface {
	fmt.Stringer
	ID() ID
	Symbol() Symbol
}

type Symbol interface {
	fmt.Stringer
	ID() ID
	Name() Name
}
