package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	// Base error
	CodeUnknown = -1

	// System error
	CodeNoRouter = 10001
	CodeNoMethod = 10002
	CodeDatabase = 10003
	CodeIndexer  = 10004

	// Service error
	CodeInvalidParams = 20001

	// Base error
	ErrorUnknown = errors.New("unknown")

	// System error
	ErrorNoRouter = errors.New("no router")
	ErrorNoMethod = errors.New("no method")
	ErrorDatabase = errors.New("database error")
	ErrorIndexer  = errors.New("indexer error")

	// Service error
	ErrorInvalidParams = errors.New("invalid params")
)

var (
	errorMap = map[int]error{
		CodeUnknown: ErrorUnknown,

		CodeNoRouter: ErrorNoRouter,
		CodeNoMethod: ErrorNoMethod,
		CodeDatabase: ErrorDatabase,
		CodeIndexer:  ErrorIndexer,

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

func SetError(c *gin.Context, err error, value error) {
	code := ErrorToCode(err)
	c.Set("code", code)
	_ = c.Error(fmt.Errorf("%w: %s", err, value.Error()))
}

func init() {
	for code, id := range errorMap {
		codeMap[id.Error()] = code
	}
}
