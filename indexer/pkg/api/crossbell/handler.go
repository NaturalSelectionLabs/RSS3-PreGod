package crossbell

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
)

func Run() error {
	instance := New(&Config{
		RPC: config.Config.Crossbell.RPC,
	})

	return instance.Run()
}
