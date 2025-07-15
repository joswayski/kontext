package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	glideConfig "github.com/joswayski/kontext/glide/config"
	globalConfig "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Eventually we will move this to it's own repo but for now, its just a folder inside the main Kontext repo
// This should simulate a ridesharing app called Glide which sends events into the clusters defined in /config/config.go
func main() {
	globalConfig := globalConfig.GetConfig()
	kafkaClustersAndClients := kafka.GetKafkaClustersFromConfig(*globalConfig)

	// TODO start a goroutine for all of these

	for clusterName, clusterDataConfig := range glideConfig.ClusterDataConfig {
		clusterConfig, exists := kafkaClustersAndClients[clusterName]

		if !exists {
			log.Fatalf("the cluster %s was configured in the Glide application but was not present in the environment variables", clusterName)
		}

		topicsInCluster, err := clusterConfig.AdminClient.ListTopics(context.Background())

		if err != nil {
			log.Fatalf("unable to retrieve topic metadata in %s cluster with error: %s", clusterName, err.Error())
		}

		// For each topic in the cluster
		for topicName, topicConfig := range clusterDataConfig.Topics {
			// First check if the topic exists in cluster. If not, create it
			topicExists := false
			for _, topic := range topicsInCluster.TopicsList().Topics() {
				if topic == topicName {
					topicExists = true
					break
				}
			}

			if !topicExists {
				// Create the topic
				_, err := clusterConfig.AdminClient.CreateTopic(
					context.Background(),
					2,   // Partitions
					1,   // Replication Factor
					nil, // Topic Config
					topicName)

				if err != nil {
					log.Fatalf("error when creating topic '%s' in cluster '%s': %s", topicName, clusterName, err.Error())
				}
			}

			sampleMessage := topicConfig.CreateMessage()
			msg, err := json.Marshal(sampleMessage)
			if err != nil {
				slog.Warn(fmt.Sprintf("unable to marshall message: %v - %s", sampleMessage, err.Error()))
			}
			record := &kgo.Record{Topic: topicName, Value: []byte(msg)}

			for {
				clusterConfig.Client.Produce(context.Background(), record, func(*kgo.Record, error) {
					if err != nil {
						slog.Warn(fmt.Sprintf("failed to produce message into topic %s in cluster %s: %v with error: %s", topicName, clusterName, sampleMessage, err.Error()))
					} else {
						slog.Info(fmt.Sprintf("message produced succesfully into %s topic at %s", topicName, time.Now()))
					}
				})
			}

		}

	}

	slog.Info("Starting Glide application!")

}
