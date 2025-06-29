package clients

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"

	cfg "github.com/joswayski/kontext/apps/api/config"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ClusterStatus struct {
	Id       string        `json:"id"`
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Metadata kadm.Metadata `json:"metadata"`
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
	Clusters []ClusterStatus `json:"clusters"`
}

func GetAllClusters(ctx context.Context, clients map[string]KafkaClients) GetAllClustersResponse {
	results := GetAllClustersResponse{
		Clusters: make([]ClusterStatus, 0),
	}
	var wg sync.WaitGroup

	for clusterName, kClients := range clients {
		wg.Add(1)
		go func(name string, client KafkaClients) {
			defer wg.Done()
			ping := kClients.client.Ping(ctx)
			healthy := ping == nil
			status := "connected"
			message := "Saul Goodman"
			if !healthy {
				results.Clusters = append(results.Clusters, ClusterStatus{
					Id:      clusterName,
					Status:  "error",
					Message: fmt.Sprintf("Unable to connect to cluster %s - error: %s", name, ping.Error()),
				})
				return
			}

			meta, err := kClients.adminClient.Metadata(ctx)

			if err != nil {
				status = "error"
				message = fmt.Sprintf("Connected to cluster but unable to retrieve metadata: %s", err.Error())

				results.Clusters = append(results.Clusters, ClusterStatus{
					Id:      clusterName,
					Status:  status,
					Message: message,
				})
				return
			}

			results.Clusters = append(results.Clusters, ClusterStatus{
				Id:       clusterName,
				Status:   status,
				Message:  message,
				Metadata: meta,
			})
		}(clusterName, kClients)
	}

	wg.Wait()
	return results
}

func formatMetadata(meta kadm.Metadata) {

}
