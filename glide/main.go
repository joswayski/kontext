package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	glideConfig "github.com/joswayski/kontext/glide/config"
	globalConfig "github.com/joswayski/kontext/pkg/config"
	kafka "github.com/joswayski/kontext/pkg/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func runProducers(kafkaClustersAndClients kafka.AllKafkaClusters) {
	for clusterName, clusterDataConfig := range glideConfig.ClusterDataConfig {
		go func(clusterName string, clusterDataConfig glideConfig.ClusterData) {
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
				go func(topicName string, topicConfig glideConfig.TopicConfig, topicsInCluster kadm.TopicDetails, clusterConfig kafka.KafkaCluster) {
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

					for {
						sampleMessage := topicConfig.CreateMessage()
						msg, err := json.Marshal(sampleMessage)
						if err != nil {
							slog.Warn(fmt.Sprintf("unable to marshall message: %v - %s", sampleMessage, err.Error()))
						} else {
							record := &kgo.Record{Topic: topicName, Value: []byte(msg)}

							clusterConfig.Client.Produce(context.Background(), record, func(r *kgo.Record, err error) {
								if err != nil {
									slog.Warn(fmt.Sprintf("failed to produce message into topic %s in cluster %s: %v with error: %s", topicName, clusterName, sampleMessage, err.Error()))
								} else {
									slog.Info(fmt.Sprintf("%s produced succesfully", topicName))
								}
							})
						}

						time.Sleep(time.Duration(1000/topicConfig.MessageRate) * time.Millisecond)

					}
				}(topicName, topicConfig, topicsInCluster, clusterConfig)
			}
		}(clusterName, clusterDataConfig)
	}

}

func runConsumers(kafkaClustersAndClients kafka.AllKafkaClusters) {
	for clusterName, clusterDataConfig := range glideConfig.ClusterDataConfig {
		go func(clusterName string, clusterDataConfig glideConfig.ClusterData) {
			clusterConfig, exists := kafkaClustersAndClients[clusterName]
			if !exists {
				log.Fatalf("the cluster %s was configured in the Glide application but was not present in the environment variables", clusterName)
			}

			topicsInCluster, err := clusterConfig.AdminClient.ListTopics(context.Background())

			if err != nil {
				log.Fatalf("unable to retrieve topic metadata in %s cluster with error: %s", clusterName, err.Error())
			}

			for topicName, topicConfig := range clusterDataConfig.Topics {
				go func(topicName string, topicConfig glideConfig.TopicConfig, topicsInCluster kadm.TopicDetails, clusterConfig kafka.KafkaCluster) {
					clusterConfig.Client.AddConsumeTopics(topicName)

					for {
						fetches := clusterConfig.Client.PollFetches(context.Background())
						if errs := fetches.Errors(); len(errs) > 0 {
							slog.Warn(fmt.Sprintf("error when polling for messages in topic %s in cluster %s: %s", topicName, clusterName, errs[0].Err.Error()))
						}

						iter := fetches.RecordIter()
						for !iter.Done() {
							record := iter.Next()
							time.Sleep(time.Duration(80 * time.Millisecond))
							fmt.Println("Consumed message!", string(record.Value))
						}
					}

				}(topicName, topicConfig, topicsInCluster, clusterConfig)
			}
		}(clusterName, clusterDataConfig)
	}
}

// Eventually we will move this to it's own repo but for now, its just a folder inside the main Kontext repo
// This should simulate a ridesharing app called Glide which sends events into the clusters defined in /config/config.go
func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	slog.Info("Starting Glide application!")

	globalConfig := globalConfig.GetConfig()
	kafkaClustersAndClients := kafka.GetKafkaClustersFromConfig(*globalConfig)

	go runProducers(kafkaClustersAndClients)
	go runConsumers(kafkaClustersAndClients)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Warn("Shutting down Glide app!")
	kafkaClustersAndClients.Close()
}
