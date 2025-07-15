package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joswayski/kontext/api/routes"
	config "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

func startServer(srv *http.Server, cfg config.KontextConfig) {
	slog.Info("Starting API server on port " + cfg.ApiPort)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error running API server", "error", err)
		os.Exit(1)
	}
}

// TODO temporary - we need to do this properly now and scale up for more data
func startConsumers(allClusters kafka.AllKafkaClusters) {
	slog.Info("Starting consumers")
	if len(allClusters) == 0 {
		slog.Warn("No Kafka clusters configured - consumers shutting down.")
		return
	}

	for clusterId, cluster := range allClusters {
		// Capture the cluster in a local variable to avoid goroutine closure issues
		clusterCopy := cluster
		clusterIdCopy := clusterId

		// Start consumer goroutine
		go func() {
			slog.Info(fmt.Sprintf("Starting consumer for cluster: %s", clusterIdCopy))

			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

				// Poll for messages (this is how franz-go consumers work)
				fetches := clusterCopy.Client.PollFetches(ctx)
				cancel()

				// Process any fetched records
				if errs := fetches.Errors(); len(errs) > 0 {
					for _, err := range errs {
						slog.Error(fmt.Sprintf("Error polling cluster %s: %v", clusterIdCopy, err))
					}
				}

				// Iterate through all fetched records
				iter := fetches.RecordIter()
				for !iter.Done() {
					record := iter.Next()
					slog.Info(fmt.Sprintf("Consumed message from cluster %s, topic %s: %s",
						clusterIdCopy, record.Topic, string(record.Value)))
				}

				// Small delay before next poll
				time.Sleep(time.Millisecond * 100)
			}
		}()

		// Start producer goroutine
		go func() {
			slog.Info(fmt.Sprintf("Starting producer for cluster: %s", clusterIdCopy))

			for {
				ctx := context.Background()
				var wg sync.WaitGroup
				wg.Add(1)

				var topic string
				if time.Now().UnixNano()%2 == 0 {
					topic = "orders"
				} else {
					topic = "users"
				}

				record := &kgo.Record{
					Topic: topic,
					Value: []byte(fmt.Sprintf("Test message from %s at %s", clusterIdCopy, time.Now().Format(time.RFC3339))),
				}

				clusterCopy.Client.Produce(ctx, record, func(_ *kgo.Record, err error) {
					defer wg.Done()
					if err != nil {
						slog.Error(fmt.Sprintf("Produce error in cluster %s: %v", clusterIdCopy, err))
					} else {
						slog.Info(fmt.Sprintf("Produced message to cluster %s, topic %s", clusterIdCopy, record.Topic))
					}
				})

				wg.Wait()

				// Wait 5 seconds before producing the next message
				time.Sleep(time.Second * 5)
			}
		}()
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
	defer kafkaClusters.Close()

	r := routes.GetRoutes(kafkaClusters)

	srv := &http.Server{
		Addr:    ":" + cfg.ApiPort,
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

	// TODO temporary
	go kafka.SeedTopics(ctx, kafkaClusters)
	go startConsumers(kafkaClusters)
	go startServer(srv, *cfg)

	awaitShutdownSignal(srv)
}
