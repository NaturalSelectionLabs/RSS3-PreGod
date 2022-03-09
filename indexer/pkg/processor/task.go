package processor

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

type Task struct {
	Identity   string
	Network    constants.NetworkID
	PlatformID constants.PlatformID // optional
}

func NewTaskQueue() chan *Task {
	return make(chan *Task)
}
