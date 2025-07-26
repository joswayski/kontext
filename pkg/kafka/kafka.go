package kafka

import (
	"fmt"
	"log"
	"log/slog"
	"sync"

	config "github.com/joswayski/kontext/pkg/config"
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

// Returns the normal client, admin client, and configs for all clusters
func GetKafkaClustersFromConfig(cfg config.KontextConfig) AllKafkaClusters {
	allClusters := make(AllKafkaClusters)

	for clusterId, clusterConfig := range cfg.KafkaClusterConfigs {
		groupId := fmt.Sprintf("kontext-%s-consumer", clusterConfig.Id)

		// Create a single client for both producing and consuming
		normalClient, err := kgo.NewClient(
			kgo.SeedBrokers(clusterConfig.BrokerURLs...),
			kgo.ConsumerGroup(groupId),
			kgo.ClientID(groupId),
		)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}

		adminClient := kadm.NewClient(normalClient)

		slog.Info(fmt.Sprintf("Created client for %s cluster", clusterId))

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
			slog.Warn(fmt.Sprintf("Kafka client for %s cluster shut down", id))
		}(id, cluster)
	}

	wg.Wait()
}
