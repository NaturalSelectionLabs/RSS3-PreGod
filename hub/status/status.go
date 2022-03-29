package status

type Code int

func (c Code) Message() Message {
	message, exist := messageMap[c]
	if exist {
		return message
	}

	return MessageUnknown
}

func (c Code) Int() int {
	return int(c)
}

type Message string

func (m Message) Code() Code {
	code, exist := codeMap[m]
	if exist {
		return code
	}

	return CodeUnknown
}

func (m Message) String() string {
	return string(m)
}

const (
	// Base error
	CodeUnknown Code = -1
	CodeOK      Code = 0

	// System error
	CodeDatabaseError Code = 20001

	// Service error
	CodeInvalidParams     Code = 10001
	CodeInvalidInstance   Code = 10002
	CodeInstanceNotExist  Code = 10003
	CodeIndexNotExist     Code = 10004
	CodeLinkListNotExist  Code = 10005
	CodeLinkNotExist      Code = 10006
	CodeFileFieldError    Code = 10007
	CodeSignatureNotExist Code = 10008

	// Base error
	MessageUnknown Message = "unknown"
	MessageOK      Message = "ok"

	// System error
	MessageDatabaseError Message = "database error"

	// Service error
	MessageInvalidParams    Message = "invalid params"
	MessageInvalidInstance  Message = "invalid instance"
	MessageInstanceNotExist Message = "instance not exist"
	MessageIndexNotExist    Message = "index not exist"
	MessageLinkListNotExist Message = "link list not exist"
	MessageLinkNotExist     Message = "link not exist"
	MessageFileFieldError   Message = "file field error"
)

var (
	messageMap = map[Code]Message{
		CodeUnknown: MessageUnknown,
		CodeOK:      MessageOK,

		CodeDatabaseError: MessageDatabaseError,

		CodeInvalidParams:    MessageInvalidParams,
		CodeInvalidInstance:  MessageInvalidInstance,
		CodeInstanceNotExist: MessageInstanceNotExist,
		CodeIndexNotExist:    MessageIndexNotExist,
		CodeLinkListNotExist: MessageLinkListNotExist,
		CodeLinkNotExist:     MessageLinkNotExist,
		CodeFileFieldError:   MessageFileFieldError,
	}
	codeMap = map[Message]Code{}
)

func init() {
	for name, id := range messageMap {
		codeMap[id] = name
	}
}
