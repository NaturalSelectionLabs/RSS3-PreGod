package status

// Gets error message from Code.
func GetMsg(code Code) Msg {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
