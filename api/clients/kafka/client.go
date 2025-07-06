package clients

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	config "github.com/joswayski/kontext/api/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// The client, admin client, and config for a cluster
type KafkaCluster struct {
	client      *kgo.Client
	adminClient *kadm.Client
	config      config.KafkaClusterConfig
}

// All clusters with their client, admin client, and config
type AllKafkaClusters map[string]KafkaCluster

func newKafkaClient(kafkaConfig config.KafkaClusterConfig) (*kgo.Client, error) {
	groupId := fmt.Sprintf("kontext-%s-consumer", kafkaConfig.Id)
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(groupId),
		kgo.ClientID(groupId),
		kgo.ConsumeTopics(topics...),
	)

	// For debugging / testing 2 consumer groups
	// TODO - temporary
	groupId2 := fmt.Sprintf("kontext-%s-consumer-2", kafkaConfig.Id)
	kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(groupId2),
		kgo.ClientID(groupId2),
		kgo.ConsumeTopics(topics...),
	)

	if kafkaConfig.Id == "production" {
		cc := cl.GetConsumeTopics()
		slog.Info(fmt.Sprintf("topic configs %s", cc))

		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

				slog.Info("Polling kafka prod")
				cl.PollFetches(ctx)
				cancel()
			}
		}()
	}

	if err != nil {
		slog.Error(fmt.Sprintf("Could not get Kafka client for %s cluster. Error: %s", kafkaConfig.Id, err))
		return nil, err
	}

	return cl, nil
}

func newAdminKafkaClient(kgoClient *kgo.Client) *kadm.Client {
	acl := kadm.NewClient(
		kgoClient,
	)

	return acl
}
