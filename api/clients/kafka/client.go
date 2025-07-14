package clients

import (
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/joswayski/kontext/pkg/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// The client, admin client, and config for a cluster
type KafkaCluster struct {
	Client      *kgo.Client
	adminClient *kadm.Client
	config      config.KafkaClusterConfig
}

// All clusters with their client, admin client, and config
type AllKafkaClusters map[string]KafkaCluster

// Returns the normal client, admin client, and configs for all clusters
func GetKafkaClustersFromConfig(cfg config.KontextConfig) AllKafkaClusters {
	allClusters := make(AllKafkaClusters)

	for clusterId, clusterConfig := range cfg.KafkaClusterConfigs {
		groupId := fmt.Sprintf("kontext-%s-consumer", clusterConfig.Id)
		normalClient, err := kgo.NewClient(
			kgo.SeedBrokers(clusterConfig.BrokerURLs...),
			kgo.ConsumerGroup(groupId),
			kgo.ClientID(groupId),
			kgo.ConsumeTopics(topics...),
		)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}

		adminClient := kadm.NewClient(normalClient)

		slog.Info(fmt.Sprintf("Created clients for %s cluster", clusterId))

		allClusters[clusterId] = KafkaCluster{
			Client:      normalClient,
			adminClient: adminClient,
			config:      clusterConfig,
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
