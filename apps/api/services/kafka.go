package services

import (
	"fmt"
	"log/slog"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewKafkaClient(kafkaConfig cfg.KafkaClusterConfig) *kgo.Client {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(fmt.Sprintf("kontext-%s-consumer", kafkaConfig.Id)))

	if err != nil {
		slog.Error("Could not")
	}
	return cl
}
