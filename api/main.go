package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joswayski/kontext/api/routes"
	config "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.GetConfig()
	kafkaClusters := kafka.GetKafkaClustersFromConfig(*cfg)

	r := routes.GetRoutes(kafkaClusters)
	srv := &http.Server{
		Addr:    ":" + cfg.ApiPort,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("Starting API server", "port", cfg.ApiPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error running API server", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	slog.Warn("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(shutdownCtx)

	g.Go(func() error {
		slog.Info("Shutting down HTTP server...")
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		slog.Info("Closing Kafka clusters...")
		if err := kafkaClusters.Close(ctx); err != nil {
			return fmt.Errorf("kafka shutdown: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("Shutdown completed with errors", "error", err)
	} else {
		slog.Info("Shutdown completed cleanly")
	}
}
