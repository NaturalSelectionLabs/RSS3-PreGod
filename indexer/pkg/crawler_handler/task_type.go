package crawler_handler

type ProcessTaskType int

const (
	ProcessTaskTypeNontType       ProcessTaskType = iota
	ProcessTaskTypeItemStorage    ProcessTaskType = 1
	ProcessTaskTypeUserBioStorage ProcessTaskType = 2
)

type ProcessTaskErrorCode int

const (
	ProcessTaskErrorCodeSuccess             ProcessTaskErrorCode = iota
	ProcessTaskErrorCodeNotFoundData        ProcessTaskErrorCode = 1
	ProcessTaskErrorCodeNotSupportedNetwork ProcessTaskErrorCode = 2
)
