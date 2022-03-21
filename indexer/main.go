package main

import (
	"log"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arweave"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/processor"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
	"github.com/RichardKnop/machinery/v1/tasks"
	jsoniter "github.com/json-iterator/go"
)

var jsoni = jsoniter.ConfigCompatibleWithStandardLibrary

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

	if err := processor.Setup(); err != nil {
		log.Fatalf("processor.Setup err: %v", err)
	}
}

func dispatchTasks(ti time.Duration) {
	// TODO: Get all accounts
	instances := []rss3uri.PlatformInstance{}
	for _, i := range instances {
		for _, n := range i.Platform.ID().GetNetwork() {
			time.Sleep(ti)

			// marshal WorkParam to string so it's supported by machinery
			param, err := jsoni.MarshalToString(crawler.WorkParam{Identity: i.Identity, PlatformID: i.Platform.ID(), NetworkID: n})

			if err != nil {
				logger.Errorf("dispatchTasks WorkParam mashalling error: %v", err)

				return
			}

			crawlerTask := tasks.Signature{
				// the name is defined by RegisterTasks() in processor/processor.go
				Name: "dispatch",
				Args: []tasks.Arg{
					{
						Type:  "string",
						Value: param,
					},
				},
			}

			processor.SendTask(crawlerTask)
		}
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

	// arweave crawler
	ar := arweave.NewArCrawler(
		1,
		500,
		10,
		2,
		600,
		"Ky1c1Kkt-jZ9sY1hvLF5nCf6WWdBhIU5Un_BMYh-t3c")
	ar.Start()

	// gitcoin crawler
	ethParam := gitcoin.Param{
		FromHeight:    1,
		Step:          10000,
		MinStep:       10,
		Confirmations: 10,
		SleepInterval: 600,
	}

	polygonParam := gitcoin.Param{
		FromHeight:    1,
		Step:          10000,
		MinStep:       10,
		Confirmations: 10,
		SleepInterval: 600,
	}

	zkParam := gitcoin.Param{
		FromHeight:    1,
		Step:          10000,
		MinStep:       10,
		Confirmations: 10,
		SleepInterval: 600,
	}
	gc := gitcoin.NewGitcoinCrawler(ethParam, polygonParam, zkParam)

	gc.PolygonStart()
	gc.EthStart()
	gc.ZkStart()

	defer logger.Logger.Sync()

	// TODO: adjust interval
	dispatchTasks(time.Minute)
}
