package api

import "errors"

var (
	// Base error
	CodeUnknown = -1

	// System error
	CodeNoRouter      = 10001
	CodeNoMethod      = 10002
	CodeNotFound      = 10003
	CodeDatabaseError = 10004

	// Service error
	CodeInvalidParams = 20001

	// Base error
	ErrorUnknown = errors.New("unknown")

	// System error
	ErrorNoRouter      = errors.New("no router")
	ErrorNoMethod      = errors.New("no method")
	ErrorNotFound      = errors.New("not found")
	ErrorDatabaseError = errors.New("database error")

	// Service error
	ErrorInvalidParams = errors.New("invalid params")
)

var (
	errorMap = map[int]error{
		CodeUnknown: ErrorUnknown,

		CodeNoRouter:      ErrorNoRouter,
		CodeNoMethod:      ErrorNoMethod,
		CodeNotFound:      ErrorNotFound,
		CodeDatabaseError: ErrorDatabaseError,

		CodeInvalidParams: ErrorInvalidParams,
	}
	codeMap = map[string]int{}
)

func ErrorToCode(err error) int {
	if code, exist := codeMap[err.Error()]; exist {
		return code
	}

	return CodeUnknown
}

func CodeToError(code int) error {
	if err, exist := errorMap[code]; exist {
		return err
	}

	return ErrorUnknown
}

func init() {
	for code, id := range errorMap {
		codeMap[id.Error()] = code
	}
}
