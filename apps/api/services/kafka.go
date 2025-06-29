package services

import (
	"fmt"
	"log/slog"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

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
			slog.Warn(fmt.Sprintf("Unable to create Kafka client for %s cluster", clusterId))
		}
		allClients[clusterId] = client
	}

	return allClients
}
