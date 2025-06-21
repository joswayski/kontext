package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/types"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// KafkaService manages connections to multiple Kafka clusters
type KafkaService struct {
	clients map[string]*kgo.Client
	mu      sync.RWMutex
}

// ClusterInfo represents detailed information about a Kafka cluster
type ClusterInfo struct {
	ID               string `json:"id"`
	BootstrapServers string `json:"bootstrap_servers"`
	Status           string `json:"status"`
	Error            string `json:"error,omitempty"`
}

// NewKafkaService creates a new Kafka service instance
func NewKafkaService(cfg *config.Config) *KafkaService {
	service := &KafkaService{
		clients: make(map[string]*kgo.Client),
	}

	// Initialize connections to all configured clusters
	for _, cluster := range cfg.Kafka.Clusters {
		service.connectToCluster(cluster)
	}

	return service
}

// connectToCluster establishes a connection to a single Kafka cluster
func (ks *KafkaService) connectToCluster(cluster config.ClusterConfig) {
	clientID := fmt.Sprintf("kontext-%s", cluster.Id)

	opts := []kgo.Opt{
		kgo.SeedBrokers(cluster.BootstrapServers...),
		kgo.ClientID(clientID),
		kgo.RequestTimeoutOverhead(5 * time.Second),
		kgo.RetryTimeout(10 * time.Second),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Printf("Failed to create Kafka client for cluster %s: %v", cluster.Id, err)
		return
	}

	ks.mu.Lock()
	ks.clients[cluster.Id] = client
	ks.mu.Unlock()

	log.Printf("Successfully connected to Kafka cluster: %s", cluster.Id)
}

// GetClusterInfo returns information about all clusters including their connection status
func (ks *KafkaService) GetClusterInfo(cfg *config.Config) []ClusterInfo {
	var clusterInfos []ClusterInfo

	for _, cluster := range cfg.Kafka.Clusters {
		info := ClusterInfo{
			ID:               strings.ToLower(cluster.Id),
			BootstrapServers: strings.Join(cluster.BootstrapServers, ","),
		}

		ks.mu.RLock()
		_, exists := ks.clients[cluster.Id]
		ks.mu.RUnlock()

		if !exists {
			info.Status = "disconnected"
			info.Error = "Client not initialized"
		} else {
			info.Status = "connected"
		}

		clusterInfos = append(clusterInfos, info)
	}

	return clusterInfos
}

func (ks *KafkaService) Close() {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	for name, client := range ks.clients {
		log.Printf("Closing connection to cluster: %s", name)
		client.Close()
	}
	ks.clients = make(map[string]*kgo.Client)
}

func (ks *KafkaService) GetClient(clusterName string) (*kgo.Client, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	client, exists := ks.clients[clusterName]
	return client, exists
}

// GetTopics retrieves all topics from a specific cluster
func (ks *KafkaService) GetTopics(clusterId string) ([]types.TopicResponse, error) {
	client, exists := ks.GetClient(clusterId)
	if !exists {
		return nil, fmt.Errorf("cluster %s not found or not connected", clusterId)
	}

	// Create metadata request to get topic information
	req := kmsg.NewMetadataRequest()
	req.Topics = nil // nil means all topics

	// Send the request
	resp, err := req.RequestWith(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %v", err)
	}

	var topics []types.TopicResponse

	for _, topic := range resp.Topics {
		// Skip internal topics unless specifically requested
		if topic.IsInternal && !strings.HasPrefix(*topic.Topic, "__") {
			continue
		}

		// Get partition count and replication factor
		partitions := int32(len(topic.Partitions))
		var replicationFactor int16
		if len(topic.Partitions) > 0 {
			replicationFactor = int16(len(topic.Partitions[0].Replicas))
		}

		// Get topic configs
		configs, err := ks.getTopicConfigs(client, *topic.Topic)
		if err != nil {
			log.Printf("Warning: failed to get configs for topic %s: %v", *topic.Topic, err)
		}

		topics = append(topics, types.TopicResponse{
			Name:              *topic.Topic,
			Partitions:        partitions,
			ReplicationFactor: replicationFactor,
			Configs:           configs,
			IsInternal:        topic.IsInternal,
		})
	}

	return topics, nil
}

// getTopicConfigs retrieves configuration for a specific topic
func (ks *KafkaService) getTopicConfigs(client *kgo.Client, topicName string) ([]string, error) {
	req := kmsg.NewDescribeConfigsRequest()
	req.Resources = []kmsg.DescribeConfigsRequestResource{
		{
			ResourceType: kmsg.ConfigResourceTypeTopic,
			ResourceName: topicName,
		},
	}

	resp, err := req.RequestWith(context.Background(), client)
	if err != nil {
		return nil, err
	}

	var configs []string
	if len(resp.Resources) > 0 {
		for _, config := range resp.Resources[0].Configs {
			configs = append(configs, fmt.Sprintf("%s=%s", config.Name, config.Value))
		}
	}

	return configs, nil
}
