package main

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/task/internal/service"
)

var s *service.Service

func main() {
	// load config
	if err := config.Setup(); err != nil {
		logger.Errorf("task: load config error, %v", err)

		return
	}

	s = service.NewService()

	s.SubscribeEns()
}
