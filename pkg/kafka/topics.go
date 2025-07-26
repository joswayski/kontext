package kafka

import (
	"context"
	"fmt"
	"sync"
)

type TopicInCluster struct {
	Name               string `json:"name"`
	TotalSize          int    `json:"total_size"`
	ConsumerGroupCount int    `json:"consumer_group_count"`
}

type AllTopicsInCluster struct {
	Topics []TopicInCluster `json:"topics"`
	// For all topics
	TotalSize int `json:"total_size"`
}

type DetailedTopic struct {
	TopicInCluster
	ConsumerGroups []ConsumerGroupInCluster `json:"consumer_groups"`
}

type GetTopicsByClusterResult struct {
	Topics     []DetailedTopic `json:"topics"`
	TopicCount int             `json:"topic_count"`
}

func GetTopicsByCluster(ctx context.Context, clients AllKafkaClusters, clusterId string) (GetTopicsByClusterResult, error) {

	var wg sync.WaitGroup
	wg.Add(2)

	var allTopics AllTopicsInCluster
	var allTopicsError error

	go func() {
		defer wg.Done()
		allTopics, allTopicsError = GetTopicsInCluster(ctx, clients[clusterId])
	}()

	var topicsAndConsumerGroups AllConsumerGroupsInTopics
	var topcisAndConsumerGroupsError error

	go func() {
		defer wg.Done()
		topicsAndConsumerGroups, topcisAndConsumerGroupsError = getConsumerGroupsForAllTopics(ctx, clients[clusterId])
	}()

	wg.Wait()

	finalTopicList := make([]DetailedTopic, 0)

	if allTopicsError != nil {
		return GetTopicsByClusterResult{Topics: finalTopicList}, fmt.Errorf("unable to retrieve topics %s", allTopicsError.Error())
	}

	if topcisAndConsumerGroupsError != nil {
		return GetTopicsByClusterResult{Topics: finalTopicList}, fmt.Errorf("unable to retrieve consumer groups for topics %s", topcisAndConsumerGroupsError.Error())
	}

	topicCount := 0
	for _, topic := range allTopics.Topics {
		consumerGroups := topicsAndConsumerGroups[topic.Name]
		if consumerGroups == nil {
			consumerGroups = make([]ConsumerGroupInCluster, 0)
		}
		detailedTopic := DetailedTopic{
			TopicInCluster: topic,
			ConsumerGroups: consumerGroups,
		}
		finalTopicList = append(finalTopicList, detailedTopic)
		topicCount += 1
	}

	return GetTopicsByClusterResult{
		Topics:     finalTopicList,
		TopicCount: topicCount,
	}, nil
}

type AllConsumerGroupsInTopics = map[string][]ConsumerGroupInCluster

func getConsumerGroupsForAllTopics(ctx context.Context, cluster KafkaCluster) (AllConsumerGroupsInTopics, error) {
	allGroups, err := cluster.AdminClient.DescribeGroups(ctx)

	if err != nil {
		return nil, err
	}

	allTopics := make(AllConsumerGroupsInTopics)

	for _, group := range allGroups {
		topics := group.AssignedPartitions().Topics()

		for _, topic := range topics {
			cg := ConsumerGroupInCluster{
				Name:        group.Group,
				State:       ConsumerGroupState(group.State),
				MemberCount: len(group.Members),
			}
			allTopics[topic] = append(allTopics[topic], cg)
		}
	}

	return allTopics, nil
}

func GetTopicsInCluster(ctx context.Context, cluster KafkaCluster) (AllTopicsInCluster, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	var topicSizeData GetSizesForEachTopicResult
	var topicSizeDataError error
	go func() {
		defer wg.Done()
		topicSizeData, topicSizeDataError = GetTopicSizes(ctx, cluster)
	}()

	var consumerGroupsInTopics AllConsumerGroupsInTopics
	var consumerGroupsInTopicsError error

	go func() {
		defer wg.Done()
		consumerGroupsInTopics, consumerGroupsInTopicsError = getConsumerGroupsForAllTopics(ctx, cluster)
	}()

	if topicSizeDataError != nil {
		return AllTopicsInCluster{}, fmt.Errorf("unable to retrieve topic sizes in cluster %s", topicSizeDataError.Error())
	}

	if consumerGroupsInTopicsError != nil {
		return AllTopicsInCluster{}, fmt.Errorf("unable to retrieve consumer groups per topic %s", consumerGroupsInTopicsError.Error())
	}

	wg.Wait()
	allTopics := make([]TopicInCluster, 0)

	for topic, sizeData := range topicSizeData.Topics {
		allTopics = append(allTopics, TopicInCluster{
			Name:               topic,
			TotalSize:          sizeData,
			ConsumerGroupCount: len(consumerGroupsInTopics[topic]),
		})
	}

	return AllTopicsInCluster{
		TotalSize: topicSizeData.TotalSize,
		Topics:    allTopics,
	}, nil
}

type GetSizesForEachTopicResult struct {
	Topics    map[string]int `json:"topics"`
	TotalSize int            `json:"total_size"`
}

func GetTopicSizes(ctx context.Context, cluster KafkaCluster) (GetSizesForEachTopicResult, error) {
	logDirs, logDirsErr := cluster.AdminClient.DescribeAllLogDirs(ctx, nil)

	if logDirsErr != nil {
		return GetSizesForEachTopicResult{}, fmt.Errorf("unable to retrieve sizes of topics %s", logDirsErr.Error())
	}

	finalResult := GetSizesForEachTopicResult{
		Topics:    make(map[string]int),
		TotalSize: 0,
	}

	// This doesn't pull internal topics, we want to ignore them for now
	listedTopics, topicsErr := cluster.AdminClient.ListTopics(ctx)
	if topicsErr != nil {
		return GetSizesForEachTopicResult{}, fmt.Errorf("unable to retrieve sizes of topics %s", topicsErr.Error())

	}
	topics := listedTopics.TopicsSet()

	for _, brokerLogDirs := range logDirs {
		if brokerLogDirs.Error() != nil {
			return GetSizesForEachTopicResult{}, fmt.Errorf("error retrieving log directories for brokers%s: %s", cluster.Config.Id, brokerLogDirs.Error())
		}

		for _, logDir := range brokerLogDirs {
			for _, partitionMap := range logDir.Topics {
				for _, partitionData := range partitionMap {
					_, exists := topics[partitionData.Topic]
					if !exists {
						// Skip internal topics not in the list above
						continue
					}
					finalResult.Topics[partitionData.Topic] += int(partitionData.Size)
					finalResult.TotalSize += int(partitionData.Size)
				}
			}
		}
	}

	return finalResult, nil
}
