package status

type Code uint16
type Msg string

const (
	SUCCESS        Code = 200
	ERROR          Code = 500
	INVALID_PARAMS Code = 400
)

const (
	SUCCESS_MSG        Msg = "Ok"
	ERROR_MSG          Msg = "Error"
	INVALID_PARAMS_MSG Msg = "Invalid params"
)

var MsgFlags = map[Code]Msg{
	SUCCESS:        SUCCESS_MSG,
	ERROR:          ERROR_MSG,
	INVALID_PARAMS: INVALID_PARAMS_MSG,
}
