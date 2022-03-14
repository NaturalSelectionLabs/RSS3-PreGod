package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router"
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
		log.Fatalf("config.Setup err: %v", err)
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
		RunMode:      config.Config.HubServer.RunMode,
		HttpPort:     config.Config.HubServer.HttpPort,
		ReadTimeout:  config.Config.HubServer.ReadTimeout,
		WriteTimeout: config.Config.HubServer.WriteTimeout,
		Handler:      router.InitRouter(),
	}

	addr := srv.Start()

	logger.Infof("Start http server listening on http://%s", addr)
	defer logger.Logger.Sync()
}
