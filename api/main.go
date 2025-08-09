package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joswayski/kontext/api/routes"
	config "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
	"golang.org/x/sync/errgroup"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

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

	g, gCtx := errgroup.WithContext(shutdownCtx)

	g.Go(func() error {
		slog.Info("Shutting down HTTP server")
		if err := srv.Shutdown(gCtx); err != nil {
			return fmt.Errorf("http shutdown error: %w", err)
		}
		slog.Info("HTTP server shut down")
		return nil
	})

	g.Go(func() error {
		slog.Info("Closing Kafka clients")
		if err := kafkaClusters.Close(gCtx); err != nil {
			return fmt.Errorf("kafka shutdown error: %w", err)
		}
		slog.Info("Kafka clients closed")
		return nil
	})

	if err := g.Wait(); err != nil {
		slog.Error("Shutdown completed with errors", "error", err)
	} else {
		slog.Info("Shutdown completed cleanly")
	}
}
