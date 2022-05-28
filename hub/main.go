package main

import (
	"os"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache" // will auto Setup by `init()`
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func init() {
	if err := database.Setup(); err != nil {
		logger.Fatalf("database.Setup err: %v", err)
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
			semconv.ServiceNameKey.String("pregod-hub"),
			semconv.ServiceVersionKey.String("v0.4.0"),
		)),
	))
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
