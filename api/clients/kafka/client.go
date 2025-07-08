package clients

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/joswayski/kontext/api/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// The client, admin client, and config for a cluster
type KafkaCluster struct {
	Client      *kgo.Client
	AdminClient *kadm.Client
	Config      config.KafkaClusterConfig
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

// Returns the normal client, admin client, and configs for all clusters
func GetKafkaClustersFromConfig(cfg config.KontextConfig) AllKafkaClusters {
	allClusters := make(AllKafkaClusters)

	for clusterId, clusterConfig := range cfg.KafkaClusterConfigs {
		normalClient, err := newKafkaClient(clusterConfig)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}

		adminClient := kadm.NewClient(normalClient)

		slog.Info(fmt.Sprintf("Created clients for %s cluster", clusterId))

		allClusters[clusterId] = KafkaCluster{
			Client:      normalClient,
			AdminClient: adminClient,
			Config:      clusterConfig,
		}
	}

	return allClusters
}

func (clusters AllKafkaClusters) Close() {
	var wg sync.WaitGroup
	for id, cluster := range clusters {
		wg.Add(1)

		go func(id string, cluster KafkaCluster) {
			defer wg.Done()
			slog.Warn(fmt.Sprintf("Shutting down Kafka client for %s cluster", id))
			cluster.Client.Close()
			slog.Warn(fmt.Sprintf("Kafka client for %s cluster shut down at", id))
		}(id, cluster)
	}

	wg.Wait()
}
