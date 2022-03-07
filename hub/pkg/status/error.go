package status

import "errors"

var (
	ErrorInvalidURI       = errors.New("invalid uri")
	ErrorNotFound         = errors.New("not found")
	ErrorMethodNotAllowed = errors.New("method not allowed")
	ErrorAccountNotFound  = errors.New("account not found")
)
