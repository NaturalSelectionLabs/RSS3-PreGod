package main

import (
	"context"
	"log"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arweave"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/crossbell"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func init() {
	if err := cache.Setup(); err != nil {
		log.Fatalf("cache.Setup err: %v", err)
	}

	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
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

	defer logger.Logger.Sync()

	srv.Start()

	return nil
}

// RunAutoUpdater runs every 10 minutes
func RunAutoUpdater(cmd *cobra.Command, args []string) error {
	logger.Info("Start refreshing recent visiters' data")

	return autoupdater.RunRecentVisitQueue(context.Background())
}

func RunAutoCrawler(cmd *cobra.Command, args []string) error {
	eg := errgroup.Group{}

	// GitCoin crawler
	eg.Go(func() error {
		return gitcoin.Start(gitcoin.ETH, gitcoin.Polygon, gitcoin.ZkSync)
	})

	// Arweave crawler
	ar := arweave.NewCrawler(arweave.MirrorUploader, arweave.DefaultCrawlConfig)

	eg.Go(ar.Start)

	// Crossbell crawler
	eg.Go(crossbell.Run)

	return eg.Wait()
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
	rootCmd.AddCommand(&cobra.Command{
		Use:  "autocrawler",
		RunE: RunAutoCrawler,
	})

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
