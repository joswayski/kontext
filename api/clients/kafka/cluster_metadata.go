package clients

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/twmb/franz-go/pkg/kadm"
)

type ClusterMetaData struct {
	Id                 string `json:"id"`
	Status             string `json:"status"`
	Message            string `json:"message,omitempty"`
	BrokerCount        int    `json:"broker_count"`
	TopicCount         int    `json:"topic_count"`
	ConsumerGroupCount int    `json:"consumer_group_count"`
	TotalSize          int64  `json:"total_size"`
}

type GetMetadataForAllClustersResponse struct {
	Clusters     []ClusterMetaData `json:"clusters"`
	ClusterCount int               `json:"cluster_count"`
}

func getMetadataForCluster(ctx context.Context, cluster KafkaCluster) ClusterMetaData {
	var wg sync.WaitGroup
	wg.Add(3)

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

	var consumerGroups kadm.ListedGroups
	var consumerGroupsError error

	go func() {
		defer wg.Done()
		consumerGroups, consumerGroupsError = cluster.adminClient.ListGroups(ctx)
	}()
	wg.Wait()

	status := "connected"

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

	if consumerGroupsError != nil {
		msg := fmt.Sprintf("Unable to retrieve consumer groups: %s.", consumerGroupsError.Error())
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
		slog.Info(fmt.Sprintf("Cluster %s - All topics:", cluster.config.Id))
		for _, topic := range metadata.Topics {
			slog.Info(fmt.Sprintf("  Topic: %s, Internal: %t", topic.Topic, topic.IsInternal))
			if !topic.IsInternal {
				// In the future I might revisit this but for now,
				// I only actually care about the 'main' topics
				topicCount += 1
			}
		}
		slog.Info(fmt.Sprintf("Cluster %s - Non-internal topic count: %d", cluster.config.Id, topicCount))
	}

	consumerGroupCount := 0
	if consumerGroups != nil {
		consumerGroupCount = len(consumerGroups.Groups())
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
		Id:                 cluster.config.Id,
		Status:             status,
		BrokerCount:        brokerCount,
		TopicCount:         topicCount,
		ConsumerGroupCount: consumerGroupCount,
		TotalSize:          totalClusterSize,
	}
}

func GetMetadataForAllClusters(ctx context.Context, clients AllKafkaClusters) GetMetadataForAllClustersResponse {
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
