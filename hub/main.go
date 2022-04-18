package main

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache" // will auto Setup by `init()`
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/es"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
)

func init() {
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
