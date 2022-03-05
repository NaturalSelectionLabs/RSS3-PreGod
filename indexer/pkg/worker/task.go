package worker

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

// TODO: move this to indexer/pkg/api/moralis/task.go
type Task struct {
	Identity string
	Network  constants.NetworkName
}
