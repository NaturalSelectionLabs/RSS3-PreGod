package main

import (
	"context"
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arweave"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
	"github.com/spf13/cobra"
)

func init() {
	if err := cache.Setup(); err != nil {
		log.Fatalf("cache.Setup err: %v", err)
	}

	if err := db.Setup(); err != nil {
		log.Fatalf("web.Setup err: %v", err)
	}
}

func RunHTTPServer(cmd *cobra.Command, args []string) error {
	srv := &web.Server{
		RunMode:      config.Config.Indexer.Server.RunMode,
		HttpPort:     config.Config.Indexer.Server.HttpPort,
		ReadTimeout:  config.Config.Indexer.Server.ReadTimeout,
		WriteTimeout: config.Config.Indexer.Server.WriteTimeout,
		Handler:      router.InitRouter(),
	}

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
	ethParam := gitcoin.NewParam(1, 10000, 10, 10, 600)
	polygonParam := gitcoin.NewParam(1, 10000, 10, 10, 600)
	zkParam := gitcoin.NewParam(1, 10000, 10, 10, 600)
	gc := gitcoin.NewGitcoinCrawler(ethParam, polygonParam, zkParam)

	go gc.PolygonStart()
	go gc.EthStart()
	go gc.ZkStart()

	defer logger.Logger.Sync()

	srv.Start()

	return nil
}

// runs every 10 minutes
func RunAutoUpdater(cmd *cobra.Command, args []string) error {
	return autoupdater.RunRecentVisitQueue(context.Background())
}

var rootCmd = &cobra.Command{Use: "indexer"}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "httpsvc",
		RunE: RunHTTPServer,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:  "autoupdater",
		RunE: RunAutoUpdater,
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
