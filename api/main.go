package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	kafka "github.com/joswayski/kontext/api/clients/kafka"
	config "github.com/joswayski/kontext/api/config"
	"github.com/joswayski/kontext/api/routes"
)

func startServer(srv *http.Server, cfg config.KontextConfig) {
	slog.Info("Starting API server on port " + cfg.Port)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error running API server", "error", err)
		os.Exit(1)
	}
}

func startConsumers(allClusters kafka.AllKafkaClusters) {
	slog.Info("Starting consumers")
	if len(allClusters) == 0 {
		slog.Warn("No Kafka clusters configured - consumers shutting down.")
		return
	}

	// for _, cluster := range kafkaClients {}
}

func waitForShutdown(srv *http.Server) {
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

	slog.Info("API server shutdown complete")
}

func main() {
	cfg := config.GetConfig()
	kafkaClusters := kafka.GetKafkaClustersFromConfig(*cfg)

	r := routes.GetRoutes(kafkaClusters)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	ctx := context.Background()

	var topicWg sync.WaitGroup
	topicWg.Add(1)
	go func() {
		defer topicWg.Done()
		kafka.CreateTopics(ctx, kafkaClusters)
	}()
	topicWg.Wait()

	// go kafka.SeedTopics(ctx, kafkaClusters)
	// go startConsumers(kafkaClusters) // TODO temporary
	go startServer(srv, *cfg)

	waitForShutdown(srv)
}
