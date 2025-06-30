package clients

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ClusterData struct {
	Id          string   `json:"id"`
	Status      string   `json:"status"`
	Message     string   `json:"message"`
	BrokerCount int      `json:"broker_count"`
	TopicCount  int      `json:"topic_count"`
	Brokers     []string `json:"brokers"`
	TotalSize   int      `json:"total_size"`
}

func newKafkaClient(kafkaConfig cfg.KafkaClusterConfig) (*kgo.Client, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(kafkaConfig.BrokerURLs...),
		kgo.ConsumerGroup(fmt.Sprintf("kontext-%s-consumer", kafkaConfig.Id)))

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

type KafkaClients struct {
	client      *kgo.Client
	adminClient *kadm.Client
}

func GetAllKafkaClients(cfg cfg.KontextConfig) map[string]KafkaClients {
	allClients := make(map[string]KafkaClients)

	for clusterId, clusterConfig := range cfg.KafkaClusters {
		normalClient, err := newKafkaClient(clusterConfig)
		if err != nil {
			log.Fatalf("Unable to create Kafka client for %s cluster: %s", clusterId, err)
		}
		slog.Info(fmt.Sprintf("Created client for %s cluster", clusterId))

		adminClient := newAdminKafkaClient(normalClient)
		slog.Info(fmt.Sprintf("Created admin client for %s cluster", clusterId))

		clConfig := KafkaClients{
			client:      normalClient,
			adminClient: adminClient,
		}
		allClients[clusterId] = clConfig
	}

	return allClients
}

type GetAllClustersResponse struct {
	Clusters     []ClusterData `json:"clusters"`
	ClusterCount int           `json:"cluster_count"`
}

func GetAllClusters(ctx context.Context, clients map[string]KafkaClients) GetAllClustersResponse {
	results := GetAllClustersResponse{
		Clusters: make([]ClusterData, 0),
	}
	var wg1 sync.WaitGroup
	var mu sync.Mutex

	for clusterName, kClients := range clients {
		wg1.Add(1)
		go func(name string, kClients KafkaClients) {
			defer wg1.Done()

			var wg2 sync.WaitGroup
			var metadata kadm.Metadata
			var metaErr error
			var logDirs kadm.DescribedAllLogDirs
			var logDirsErr error

			wg2.Add(2)

			go func() {
				slog.Info(fmt.Sprintf("%s - Starting metadata retrieval  %s", clusterName, time.Now()))
				defer wg2.Done()
				metadata, metaErr = kClients.adminClient.Metadata(ctx)
				slog.Info(fmt.Sprintf("%s - ending metadata retrieval %s", clusterName, time.Now()))

			}()

			go func() {
				slog.Info(fmt.Sprintf("%s - Starting log retrieval %s", clusterName, time.Now()))
				defer wg2.Done()
				logDirs, logDirsErr = kClients.adminClient.DescribeAllLogDirs(ctx, nil)
				slog.Info(fmt.Sprintf("%s - ending log retrieval %s", clusterName, time.Now()))
			}()

			wg2.Wait()

			status := "connected"
			message := "Saul Goodman"

			if metaErr != nil {
				status = "error"
				message = fmt.Sprintf("Unable to retrieve metadata: %s", metaErr.Error())
			}

			if logDirsErr != nil {
				status = "error"
				message = fmt.Sprintf("Unable to retrieve cluster storage size: %s", logDirsErr.Error())
			}

			var totalClusterSize int64

			for _, brokerLogDirs := range logDirs {
				if brokerLogDirs.Error() != nil {
					slog.Warn(fmt.Sprintf("Error retrieving log directories for brokers in cluster %s", clusterName))
					continue
				}

				for _, logDir := range brokerLogDirs {
					for _, partitionMap := range logDir.Topics {
						for _, partitionData := range partitionMap {
							totalClusterSize += partitionData.Size
						}
					}

				}
			}

			brokerCount := 0
			if metadata.Brokers != nil {
				brokerCount = len(metadata.Brokers)
			}

			topicCount := 0
			if metadata.Topics != nil {
				topicCount = len(metadata.Topics)
			}
			mu.Lock()
			results.Clusters = append(results.Clusters, ClusterData{
				Id:          name,
				Status:      status,
				Message:     message,
				BrokerCount: brokerCount,
				TopicCount:  topicCount,
				// ! TODO - this is refetching the config, pass this through as state or build from metadata
				Brokers:   cfg.GetConfig().KafkaClusters[name].BrokerURLs,
				TotalSize: int(totalClusterSize),
			})
			mu.Unlock()
		}(clusterName, kClients)
	}

	wg1.Wait()

	results.ClusterCount = len(results.Clusters)
	return results
}
