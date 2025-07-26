package kafka

import (
	"context"
	"fmt"
)

type GetClusterByIdResponse struct {
	Metadata       ClusterMetaData            `json:"metadata"`
	Brokers        []string                   `json:"brokers"` // URLs for all brokers
	Topics         []TopicInCluster           `json:"topics"`
	ConsumerGroups AllConsumerGroupsInCluster `json:"consumer_groups"`
}

// Returns detailed cluster information, including brokers, topics, and consumer groups
func GetClusterById(ctx context.Context, id string, clients AllKafkaClusters) (GetClusterByIdResponse, error) {
	cluster, exists := clients[id]
	if !exists {
		return GetClusterByIdResponse{}, fmt.Errorf("cluster '%s' not found", id)
	}

	metadata := getMetadataForCluster(ctx, cluster)
	if metadata.Status == "error" {
		return GetClusterByIdResponse{}, fmt.Errorf("error retrieving metadata: %s", metadata.Message)
	}

	consumerGroups, err := getConsumerGroupsInCluster(ctx, cluster)
	if err != nil {
		return GetClusterByIdResponse{}, fmt.Errorf("could not describe groups: %w", err)
	}

	allTopicData, _ := GetTopicsInCluster(ctx, cluster)

	return GetClusterByIdResponse{
		Metadata:       metadata,
		ConsumerGroups: consumerGroups,
		Brokers:        cluster.Config.BrokerURLs,
		Topics:         allTopicData.Topics,
	}, nil
}
