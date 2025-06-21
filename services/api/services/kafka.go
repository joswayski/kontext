package services

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/joswayski/kontext/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

// KafkaService manages connections to multiple Kafka clusters
type KafkaService struct {
	clients map[string]*kgo.Client
	mu      sync.RWMutex
}

// ClusterInfo represents detailed information about a Kafka cluster
type ClusterInfo struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	BootstrapServers string    `json:"bootstrapServers"`
	Status           string    `json:"status"`
	LastChecked      time.Time `json:"lastChecked"`
	Error            string    `json:"error,omitempty"`
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
	clientID := fmt.Sprintf("kontext-%s", cluster.Name)

	opts := []kgo.Opt{
		kgo.SeedBrokers(cluster.BootstrapServers...),
		kgo.ClientID(clientID),
		kgo.RequestTimeoutOverhead(5 * time.Second),
		kgo.RetryTimeout(10 * time.Second),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Printf("Failed to create Kafka client for cluster %s: %v", cluster.Name, err)
		return
	}

	ks.mu.Lock()
	ks.clients[cluster.Name] = client
	ks.mu.Unlock()

	log.Printf("Successfully connected to Kafka cluster: %s", cluster.Name)
}

// GetClusterInfo returns information about all clusters including their connection status
func (ks *KafkaService) GetClusterInfo(cfg *config.Config) []ClusterInfo {
	var clusterInfos []ClusterInfo

	for _, cluster := range cfg.Kafka.Clusters {
		info := ClusterInfo{
			ID:               strings.ToLower(cluster.Name),
			Name:             cluster.Name,
			BootstrapServers: cluster.BootstrapServers[0], // For now, just show the first server
			LastChecked:      time.Now(),
		}

		ks.mu.RLock()
		_, exists := ks.clients[cluster.Name]
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
