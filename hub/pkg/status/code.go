package status

type (
	Code    uint
	Message string
)

const (
	CodeSuccess       Code = 200
	CodeError         Code = 500
	CodeInvalidParams Code = 400
)

const (
	MessageSuccess       Message = "ok"
	MessageError         Message = "error"
	MessageInvalidParams Message = "invalid params"
)

var messageMap = map[Code]Message{
	CodeSuccess:       MessageSuccess,
	CodeError:         MessageError,
	CodeInvalidParams: MessageInvalidParams,
}

func GetMessage(code Code) Message {
	message, ok := messageMap[code]
	if ok {
		return message
	}

	return messageMap[CodeError]
}
