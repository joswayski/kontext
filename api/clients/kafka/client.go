package clients

import (
	"context"
	"fmt"
	"log"
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

// Creates a normal kafka client
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
		// TODO - temporary
		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

				slog.Info(fmt.Sprintf("Polling kafka prod %s", time.Now()))
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

// Creates an admin client for the cluster to retrieve metadata
func newAdminKafkaClient(kgoClient *kgo.Client) *kadm.Client {
	acl := kadm.NewClient(
		kgoClient,
	)

	return acl
}

// Returns the normal client, admin client, and configs for all clusters
func GetKafkaClustersFromConfig(cfg config.KontextConfig) AllKafkaClusters {
	allClusters := make(AllKafkaClusters)

	for clusterId, clusterConfig := range cfg.KafkaClusterConfigs {
		normalClient, err := newKafkaClient(clusterConfig)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}
		slog.Info(fmt.Sprintf("Created client for %s cluster", clusterId))

		adminClient := newAdminKafkaClient(normalClient)
		slog.Info(fmt.Sprintf("Created admin client for %s cluster", clusterId))

		allClusters[clusterId] = KafkaCluster{
			client:      normalClient,
			adminClient: adminClient,
			config:      clusterConfig,
		}
	}

	return allClusters
}
