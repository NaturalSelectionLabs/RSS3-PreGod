package main

import (
	"context"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/arweave"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/gitcoin"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/subscribe/ens"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func init() {
	if err := cache.Setup(); err != nil {
		log.Fatalf("cache.Setup err: %v", err)
	}

	if err := database.Setup(); err != nil {
		log.Fatalf("database.Setup err: %v", err)
	}

	// TODO
	var exporter trace.SpanExporter

	if config.Config.OpenTelemetry == nil {
		file, err := os.Create("traces.log")
		if err != nil {
			logger.Fatal(err)
		}

		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(file),
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		var err error

		if exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Config.OpenTelemetry.URL))); err != nil {
			logger.Fatal(err)
		}
	}

	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("pregod-indexer"),
			semconv.ServiceVersionKey.String("v0.4.0"),
		)),
	))
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
	srv := &web.Server{
		RunMode:      config.Config.Indexer.Server.RunMode,
		HttpPort:     config.Config.Indexer.Server.HttpPort,
		ReadTimeout:  config.Config.Indexer.Server.ReadTimeout,
		WriteTimeout: config.Config.Indexer.Server.WriteTimeout,
	}

	// zksync
	// go zksync.Start()

	if err := gitcoin.Setup(); err != nil {
		log.Fatalf("gitcoin.Setup err: %v", err)
	}

	// TODO: remove gitcoin crawler for now
	logger.Info("Start crawling gitcoin")
	// gitcoin crawler
	go gitcoin.Start(gitcoin.Polygon)
	go gitcoin.Start(gitcoin.ETH)
	go gitcoin.Start(gitcoin.ZkSync)
	logger.Info("Start crawling arweave")

	// subscribe ens
	go ens.Run()

	//arweave crawler
	ar := arweave.NewCrawler(arweave.MirrorUploader, arweave.DefaultCrawlConfig)

	if err := ar.Start(); err != nil {
		logger.Errorf("arweave crawler start error: %v", err)
	}

	srv.Start()

	return nil
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
