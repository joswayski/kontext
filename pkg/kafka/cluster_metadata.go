package kafka

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/twmb/franz-go/pkg/kadm"
)

type MetadataStatus string

type ClusterMetaData struct {
	Id     string         `json:"id"`
	Status MetadataStatus `json:"status"`
	// Only shows if there is an error (status === "error")
	Message            string `json:"message,omitempty"`
	BrokerCount        int    `json:"broker_count"`
	TopicCount         int    `json:"topic_count"`
	ConsumerGroupCount int    `json:"consumer_group_count"`
	TotalSize          int64  `json:"total_size"`
}

const (
	MetadataStatusConnected = "connected"
	MetadataStatusError     = "error"
)

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
		metadata, metaErr = cluster.AdminClient.Metadata(ctx)
	}()

	// var logDirs kadm.DescribedAllLogDirs
	// var logDirsErr error
	// go func() {
	// 	defer wg.Done()
	// 	logDirs, logDirsErr = cluster.adminClient.DescribeAllLogDirs(ctx, nil)
	// }()

	var consumerGroups kadm.ListedGroups
	var consumerGroupsError error

	go func() {
		defer wg.Done()
		consumerGroups, consumerGroupsError = cluster.AdminClient.ListGroups(ctx)
	}()

	var topicSizeData GetSizesForEachTopicResult
	var topicSizeError error

	go func() {
		defer wg.Done()
		topicSizeData, topicSizeError = GetTopicSizes(ctx, cluster)
	}()
	wg.Wait()

	if metaErr != nil {
		msg := fmt.Sprintf("Unable to retrieve metadata: %s. Please check if the cluster is running.", metaErr.Error())
		return ClusterMetaData{
			Id:      cluster.Config.Id,
			Status:  MetadataStatusError,
			Message: msg,
		}
	}

	// if logDirsErr != nil {
	// 	msg := fmt.Sprintf("Unable to retrieve describe log dirs: %s.", logDirsErr.Error())
	// 	return ClusterMetaData{
	// 		Id:      cluster.config.Id,
	// 		Status:  MetadataStatusError,
	// 		Message: msg,
	// 	}
	// }

	if consumerGroupsError != nil {
		msg := fmt.Sprintf("Unable to retrieve consumer groups: %s.", consumerGroupsError.Error())
		return ClusterMetaData{
			Id:      cluster.Config.Id,
			Status:  MetadataStatusError,
			Message: msg,
		}
	}

	if topicSizeError != nil {
		msg := fmt.Sprintf("Unable to retrieve topic size data: %s.", consumerGroupsError.Error())
		return ClusterMetaData{
			Id:      cluster.Config.Id,
			Status:  MetadataStatusError,
			Message: msg,
		}
	}

	brokerCount := 0
	if metadata.Brokers != nil {
		brokerCount = len(metadata.Brokers)
	}

	topicCount := 0
	if metadata.Topics != nil {
		for _, topic := range metadata.Topics {
			topicCount += len(metadata.Topics)
			if !topic.IsInternal {
				// In the future I might revisit this but for now,
				// I only actually care about the 'main' topics
				topicCount += 1
			}
		}
	}

	consumerGroupCount := 0
	if consumerGroups != nil {
		consumerGroupCount = len(consumerGroups.Groups())
	}

	// Incase we ever need this in the future, this includes internal topics
	// For now, executive decision is to not include this as its not relevant for business Kontext
	// for _, brokerLogDirs := range logDirs {
	// 	if brokerLogDirs.Error() != nil {
	// 		msg := fmt.Sprintf("Error retrieving log directories for brokers%s: %s", cluster.config.Id, brokerLogDirs.Error())
	// 		return ClusterMetaData{
	// 			Id:      cluster.config.Id,
	// 			Status:  MetadataStatusError,
	// 			Message: msg,
	// 		}
	// 	}

	// 	for _, logDir := range brokerLogDirs {
	// 		for _, partitionMap := range logDir.Topics {
	// 			for _, partitionData := range partitionMap {
	// 				totalClusterSize += partitionData.Size
	// 			}
	// 		}
	// 	}
	// }

	return ClusterMetaData{
		Id:                 cluster.Config.Id,
		Status:             MetadataStatusConnected,
		BrokerCount:        brokerCount,
		TopicCount:         topicCount,
		ConsumerGroupCount: consumerGroupCount,
		TotalSize:          int64(topicSizeData.TotalSize),
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
