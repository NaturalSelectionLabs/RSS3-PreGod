package util

type ErrorCode int
type ErrorMsg string

const (
	ErrorCodeSuccess ErrorCode = iota

	// External API return problem
	ErrorCodeParameterError ErrorCode = 1000

	// Internal implementation issues
	ErrorCodeNotFoundData        ErrorCode = 2000
	ErrorCodeNotSupportedNetwork ErrorCode = 2001
	ErrorCodeGetDataError        ErrorCode = 2002
)

const (
	ErrorMsgParameterError = "Parameter Error"

	ErrorMsgNotFoundData        = "Not found data"
	ErrorMsgNotSupportedNetwork = "Not supported network"
	ErrorMsgGetDataError        = "Get data error"
)

type ErrorBase struct {
	ErrorCode ErrorCode `json:"code"`
	ErrorMsg  ErrorMsg  `json:"msg"`
}

var (
	errorMsgMap = map[ErrorCode]ErrorMsg{
		ErrorCodeSuccess:             "",
		ErrorCodeParameterError:      ErrorMsgParameterError,
		ErrorCodeNotFoundData:        ErrorMsgNotFoundData,
		ErrorCodeNotSupportedNetwork: ErrorMsgNotSupportedNetwork,
	}
)

func GetErrorBase(errorCode ErrorCode) ErrorBase {
	return ErrorBase{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsgMap[errorCode],
	}
}
