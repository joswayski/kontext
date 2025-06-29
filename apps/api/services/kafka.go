package services

import (
	"fmt"
	"log/slog"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewKafkaClient(kafkaConfig cfg.KafkaClusterConfig) (*kgo.Client, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(fmt.Sprintf("kontext-%s-consumer", kafkaConfig.Id)))

	if err != nil {
		slog.Error(fmt.Sprintf("Could not get Kafka client for %s cluster. Error: %s", kafkaConfig.Id, err))
		return nil, err
	}
	return cl, nil
}
