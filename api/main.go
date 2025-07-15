package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joswayski/kontext/api/routes"
	config "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
)

func startServer(srv *http.Server, cfg config.KontextConfig) {
	slog.Info("Starting API server on port " + cfg.ApiPort)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error running API server", "error", err)
		os.Exit(1)
	}
}

func awaitShutdownSignal(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Warn("Shutting down API server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
		os.Exit(1)
	}

	slog.Info("API server shutdown complete, cleaning up resources...")
}

func main() {
	cfg := config.GetConfig()
	kafkaClusters := kafka.GetKafkaClustersFromConfig(*cfg)

	r := routes.GetRoutes(kafkaClusters)

	srv := &http.Server{
		Addr:    ":" + cfg.ApiPort,
		Handler: r,
	}

	go startServer(srv, *cfg)

	awaitShutdownSignal(srv)
	kafkaClusters.Close()
}
