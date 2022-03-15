package main

import (
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router"
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
		log.Fatalf("cache.Setup err: %v", err)
	}

	if err := db.Setup(); err != nil {
		log.Fatalf("db.Setup err: %v", err)
	}
}

func main() {
	srv := &web.Server{
		RunMode:      config.Config.Indexer.Server.RunMode,
		HttpPort:     config.Config.Indexer.Server.HttpPort,
		ReadTimeout:  config.Config.Indexer.Server.ReadTimeout,
		WriteTimeout: config.Config.Indexer.Server.WriteTimeout,
		Handler:      router.InitRouter(),
	}

	addr := srv.Start()

	logger.Infof("Start http server listening on http://%s", addr)
	defer logger.Logger.Sync()
}
