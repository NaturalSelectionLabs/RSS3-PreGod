package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
)

func init() {
	if err := config.Setup(); err != nil {
		log.Fatalf("config.Setup err: %v", err)
	}

	if err := logger.Setup(); err != nil {
		log.Fatalf("logger.Setup err: %v", err)
	}

	if err := cache.Setup(); err != nil {
		logger.Fatalf("cache.Setup err: %v", err)
	}

	if err := database.Setup(); err != nil {
		logger.Fatalf("database.Setup err: %v", err)
	}
}

func main() {
	srv := &web.Server{
		RunMode:      config.Config.Hub.Server.RunMode,
		HttpPort:     config.Config.Hub.Server.HttpPort,
		ReadTimeout:  config.Config.Hub.Server.ReadTimeout,
		WriteTimeout: config.Config.Hub.Server.WriteTimeout,
		Handler:      router.Initialize(),
	}

	defer logger.Logger.Sync()

	srv.Start()
}
