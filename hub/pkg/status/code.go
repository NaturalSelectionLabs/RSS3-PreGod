package status

type (
	Code    uint16
	Message string
)

const (
	CodeSuccess       Code = 200
	CodeError         Code = 500
	CodeInvalidParams Code = 400

	MessageSuccess       = "ok"
	MessageError         = "error"
	MessageInvalidParams = "invalid params"
)

var messageMap = map[Code]Message{
	CodeSuccess:       MessageSuccess,
	CodeError:         MessageError,
	CodeInvalidParams: MessageInvalidParams,
}
