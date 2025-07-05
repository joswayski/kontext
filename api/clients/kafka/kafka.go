package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sort"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	cfg "github.com/joswayski/kontext/api/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ClusterMetaData struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	BrokerCount int    `json:"broker_count"`
	TopicCount  int    `json:"topic_count"`
	TotalSize   int64  `json:"total_size"`
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

type KafkaCluster struct {
	client      *kgo.Client
	adminClient *kadm.Client
	config      cfg.KafkaClusterConfig
}

// Returns the normal client, admin client, and configs for all clusters
func GetKafkaClustersFromConfig(cfg cfg.KontextConfig) map[string]KafkaCluster {
	allClusters := make(map[string]KafkaCluster)

	for clusterId, clusterConfig := range cfg.KafkaClusters {
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

type GetMetadataForAllClustersResponse struct {
	Clusters     []ClusterMetaData `json:"clusters"`
	ClusterCount int               `json:"cluster_count"`
}

func getMetadataForCluster(ctx context.Context, cluster KafkaCluster) ClusterMetaData {
	var wg sync.WaitGroup
	wg.Add(2)

	var metadata kadm.Metadata
	var metaErr error
	go func() {
		defer wg.Done()
		metadata, metaErr = cluster.adminClient.Metadata(ctx)
	}()

	var logDirs kadm.DescribedAllLogDirs
	var logDirsErr error
	go func() {
		defer wg.Done()
		logDirs, logDirsErr = cluster.adminClient.DescribeAllLogDirs(ctx, nil)
	}()

	wg.Wait()

	status := "connected"
	message := "Saul Goodman"

	if metaErr != nil {
		msg := fmt.Sprintf("Unable to retrieve metadata: %s. Please check if the cluster is running.", metaErr.Error())
		return ClusterMetaData{
			Id:      cluster.config.Id,
			Status:  "error",
			Message: msg,
		}
	}

	if logDirsErr != nil {
		msg := fmt.Sprintf("Unable to retrieve describe log dirs: %s.", logDirsErr.Error())
		return ClusterMetaData{
			Id:      cluster.config.Id,
			Status:  "error",
			Message: msg,
		}
	}

	var totalClusterSize int64

	brokerCount := 0
	if metadata.Brokers != nil {
		brokerCount = len(metadata.Brokers)
	}

	topicCount := 0
	if metadata.Topics != nil {
		topicCount = len(metadata.Topics)
	}

	for _, brokerLogDirs := range logDirs {
		if brokerLogDirs.Error() != nil {
			msg := fmt.Sprintf("Error retrieving log directories for brokers%s: %s", cluster.config.Id, brokerLogDirs.Error())
			return ClusterMetaData{
				Id:      cluster.config.Id,
				Status:  "error",
				Message: msg,
			}
		}

		for _, logDir := range brokerLogDirs {
			for _, partitionMap := range logDir.Topics {
				for _, partitionData := range partitionMap {
					totalClusterSize += partitionData.Size
				}
			}
		}
	}

	return ClusterMetaData{
		Id:          cluster.config.Id,
		Status:      status,
		Message:     message,
		BrokerCount: brokerCount,
		TopicCount:  topicCount,
		TotalSize:   totalClusterSize,
	}
}

func GetMetadataForAllClusters(ctx context.Context, clients map[string]KafkaCluster) GetMetadataForAllClustersResponse {
	results := GetMetadataForAllClustersResponse{
		Clusters: make([]ClusterMetaData, 0),
	}
	var wg sync.WaitGroup

	resultChan := make(chan ClusterMetaData, len(clients))

	for _, cluster := range clients {
		wg.Add(1)
		go func(c KafkaCluster) {
			defer wg.Done()
			resultChan <- getMetadataForCluster(ctx, c)
		}(cluster)
	}

	wg.Wait()
	close(resultChan)

	for cmd := range resultChan {
		results.Clusters = append(results.Clusters, cmd)
	}

	// Sort clusters alphabetically
	sort.Slice(results.Clusters, func(i, j int) bool {
		return results.Clusters[i].Id < results.Clusters[j].Id
	})

	results.ClusterCount = len(results.Clusters)
	return results
}

type GetClusterByIdResponse struct {
	Metadata ClusterMetaData      `json:"metadata"`
	Brokers  []kadm.BrokerDetails `json:"brokers"`
	Topics   []kadm.TopicDetails  `json:"topics"`
}

func GetClusterById(ctx context.Context, id string, clients map[string]KafkaCluster) (GetClusterByIdResponse, error) {

	cluster, exists := clients[id]
	if !exists {
		return GetClusterByIdResponse{}, fmt.Errorf("cluster '%s' not found", id)
	}

	metadata := getMetadataForCluster(ctx, cluster)
	if metadata.Status == "error" {
		return GetClusterByIdResponse{}, fmt.Errorf("error retrieving metadata: %s", metadata.Message)
	}

	cluster.adminClient.ListGroups(ctx)

	return GetClusterByIdResponse{
		Metadata: metadata,
	}, nil
}

var topics = []string{"orders", "users"}

func CreateTopics(ctx context.Context, clients map[string]KafkaCluster) {
	slog.Info("Creating topics...")
	for _, cluster := range clients {
		_, err := cluster.adminClient.CreateTopics(ctx, 1, 1, nil, topics...)
		if err != nil {
			slog.Warn("Unable to create topics")
			continue
		}
		slog.Info(fmt.Sprintf("Topics created in %s cluster", cluster.config.Id))
	}
}

// TODO temporary will cleanup in next PR
type SampleMessage struct {
	MessageType string      `json:"message_type"`
	Data        interface{} `json:"data"`
}

// TODO temporary will cleanup in next PR
func SeedTopics(ctx context.Context, clients map[string]KafkaCluster) {
	slog.Info("Seeding topics...")

	for _, topic := range topics {
		for _, cluster := range clients {
			sampleMsg := SampleMessage{
				MessageType: gofakeit.Word(),
				Data: map[string]string{
					"name": gofakeit.Name(),
				},
			}

			jsonData, err := json.Marshal(sampleMsg)
			if err != nil {
				slog.Error("Failed to marshal message", "error", err)
				continue
			}

			cluster.client.Produce(ctx, &kgo.Record{
				Topic: topic,
				Key:   []byte(gofakeit.UUID()),
				Value: jsonData,
			}, func(r *kgo.Record, err error) {
				if err != nil {
					slog.Error("Failed to produce message", "error", err, "topic", topic)
				} else {
					slog.Info("Message produced successfully", "topic", topic)
				}
			})
		}
	}
}
