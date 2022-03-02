package worker

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

type Task struct {
	Identity string
	Network  constants.NetworkName
}
