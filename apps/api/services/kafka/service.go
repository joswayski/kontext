package services

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ClusterStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newKafkaClient(kafkaConfig cfg.KafkaClusterConfig) (*kgo.Client, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(fmt.Sprintf("kontext-%s-consumer", kafkaConfig.Id)))

	if err != nil {
		slog.Error(fmt.Sprintf("Could not get Kafka client for %s cluster. Error: %s", kafkaConfig.Id, err))
		return nil, err
	}
	return cl, nil
}

func GetAllKafkaClients(cfg cfg.KontextConfig) map[string]*kgo.Client {
	allClients := make(map[string]*kgo.Client)

	for clusterId, clusterConfig := range cfg.KafkaClusters {
		client, err := newKafkaClient(clusterConfig)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}
		slog.Info(fmt.Sprintf("Created client for %s cluster", clusterId))
		allClients[clusterId] = client
	}

	return allClients
}

func GetClusterStatuses(ctx context.Context, clients map[string]*kgo.Client) map[string]ClusterStatus {
	results := make(map[string]ClusterStatus)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for clusterName, kafkaClient := range clients {
		wg.Add(1)
		go func(name string, client *kgo.Client) {
			defer wg.Done()
			ping := client.Ping(ctx)
			healthy := ping == nil
			status := "connected"
			message := "Saul Goodman"
			if !healthy {
				status = "error"
				message = ping.Error()
			}

			mu.Lock()
			results[name] = ClusterStatus{
				Status:  status,
				Message: message,
			}
			mu.Unlock()
		}(clusterName, kafkaClient)
	}

	wg.Wait()
	return results
}
