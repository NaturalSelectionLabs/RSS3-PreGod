package status

// Gets error message from Code.
func GetMsg(code Code) Message {
	msg, ok := messageMap[code]
	if ok {
		return msg
	}

	return messageMap[CodeError]
}
