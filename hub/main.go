package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/router"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/gin-gonic/gin"
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

	if err := db.Setup(); err != nil {
		logger.Fatalf("db.Setup err: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		logger.Fatalf("db.AutoMigrate err: %v", err)
	}
}

func main() {
	gin.SetMode(config.Config.HubServer.RunMode)

	addr := net.JoinHostPort("localhost", strconv.Itoa(config.Config.HubServer.HttpPort))

	server := &http.Server{
		Addr:           addr,
		Handler:        router.InitRouter(),
		ReadTimeout:    config.Config.HubServer.ReadTimeout,
		WriteTimeout:   config.Config.HubServer.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	logger.Infof("Start http server listening on http://%s", addr)
	defer logger.Logger.Sync()

	go server.ListenAndServe()

	gracefullyExit(server)
}

func gracefullyExit(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-quit

	logger.Info("Shutdown due to a signal: ", sig)

	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // with a 5s timeout
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Shutdown error:", err)
	}

	logger.Info("Shutdown server successfully in ", time.Since(now))
}
