package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	glideConfig "github.com/joswayski/kontext/glide/config"
	globalConfig "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
)

// Eventually we will move this to it's own repo but for now, its just a folder inside the main Kontext repo
// This should simulate a ridesharing app called Glide which sends events into the clusters defined in /config/config.go
func main() {
	globalConfig := globalConfig.GetConfig()
	kafkaClustersAndClients := kafka.GetKafkaClustersFromConfig(*globalConfig)

	for clusterName, clusterDataConfig := range glideConfig.ClusterDataConfig {
		clusterConfig, exists := kafkaClustersAndClients[clusterName]

		if !exists {
			log.Fatal(fmt.Sprintf("the cluster %s was configured in the Glide application but was not present in the environment variables", clusterName))
		}

		topicsInCluster, err := clusterConfig.AdminClient.ListTopics(context.Background())

		if err != nil {
			log.Fatal(fmt.Sprintf("unable to retrieve topic metadata in %s cluster", clusterName))
		}

		// For each topic in the cluster
		for topicName, topicConfig := range clusterDataConfig.Topics {

			// First check if the topic exists. If not, create it
		
		if topicsInCluster.
		
			// Create the topic
		}
			// Publish a sample message at the specified rate into the specified topic

		}
	}

	slog.Info("Starting Glide application!")

}
