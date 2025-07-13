package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/brianvoe/gofakeit"
	"github.com/twmb/franz-go/pkg/kgo"
)

type TopicInCluster struct {
	Name      string `json:"name"`
	TotalSize int    `json:"total_size"`
}

type AllTopicsInCluster struct {
	Topics    []TopicInCluster `json:"topics"`
	TotalSize int              `json:"total_size"`
}

type DetailedTopic struct {
	TopicInCluster
	ConsumerGroups []ConsumerGroupInCluster `json:"consumer_groups"`
}

type GetTopicsByClusterResult struct {
	Topics []DetailedTopic `json:"topics"`
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

	if allTopicsError != nil {
		return GetTopicsByClusterResult{}, fmt.Errorf("unable to retrieve topics %s", allTopicsError.Error())
	}

	if topcisAndConsumerGroupsError != nil {
		return GetTopicsByClusterResult{}, fmt.Errorf("unable to retrieve consumer groups for topics %s", topcisAndConsumerGroupsError.Error())
	}

	finalTopicList := make([]DetailedTopic, 0)

	for _, topic := range allTopics.Topics {
		detailedTopic := DetailedTopic{
			TopicInCluster: topic,
			ConsumerGroups: topicsAndConsumerGroups[topic.Name],
		}
		finalTopicList = append(finalTopicList, detailedTopic)

	}

	return GetTopicsByClusterResult{
		Topics: finalTopicList,
	}, nil
}

type AllConsumerGroupsInTopics = map[string][]ConsumerGroupInCluster

func getConsumerGroupsForAllTopics(ctx context.Context, cluster KafkaCluster) (AllConsumerGroupsInTopics, error) {
	allGroups, err := cluster.adminClient.DescribeGroups(ctx)

	if err != nil {
		return nil, err
	}

	allTopics := make(AllConsumerGroupsInTopics)

	for _, group := range allGroups {
		topics := group.AssignedPartitions().Topics()

		for _, topic := range topics {
			cg := ConsumerGroupInCluster{
				Name:         group.Group,
				State:        ConsumerGroupState(group.State),
				MembersCount: len(group.Members),
			}
			allTopics[topic] = append(allTopics[topic], cg)
		}
	}

	return allTopics, nil
}

func GetTopicsInCluster(ctx context.Context, cluster KafkaCluster) (AllTopicsInCluster, error) {
	topics, err := cluster.adminClient.ListTopics(ctx)
	if err != nil {
		return AllTopicsInCluster{}, fmt.Errorf("unable to retrieve topics in cluster %s", err.Error())
	}

	topicSizeData, topicSizeErr := GetTopicSizes(ctx, cluster)
	if topicSizeErr != nil {
		return AllTopicsInCluster{}, fmt.Errorf("unable to retrieve topic sizes in cluster %s", topicSizeErr.Error())
	}

	sortedTopics := topics.Sorted()

	allTopics := make([]TopicInCluster, 0)

	for _, topic := range sortedTopics {
		allTopics = append(allTopics, TopicInCluster{
			Name:      topic.Topic,
			TotalSize: topicSizeData.Topics[topic.Topic],
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
	logDirs, logDirsErr := cluster.adminClient.DescribeAllLogDirs(ctx, nil)

	if logDirsErr != nil {
		return GetSizesForEachTopicResult{}, fmt.Errorf("unable to retrieve sizes of topics %s", logDirsErr.Error())
	}

	finalResult := GetSizesForEachTopicResult{
		Topics:    make(map[string]int),
		TotalSize: 0,
	}

	// Skip internal topics
	listedTopics, topicsErr := cluster.adminClient.ListTopics(ctx)
	if topicsErr != nil {
		return GetSizesForEachTopicResult{}, fmt.Errorf("unable to retrieve sizes of topics %s", topicsErr.Error())

	}
	topics := listedTopics.TopicsSet()

	for _, brokerLogDirs := range logDirs {
		if brokerLogDirs.Error() != nil {
			return GetSizesForEachTopicResult{}, fmt.Errorf("error retrieving log directories for brokers%s: %s", cluster.config.Id, brokerLogDirs.Error())
		}

		for _, logDir := range brokerLogDirs {
			for _, partitionMap := range logDir.Topics {
				for _, partitionData := range partitionMap {
					// Skip internal topics
					_, exists := topics[partitionData.Topic]
					if !exists {
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

// TODO - temporary - will cleanup in a separate PR
var topics = []string{"orders", "users"}

// TODO - temporary - will cleanup in a separate PR
func CreateTopics(ctx context.Context, clients AllKafkaClusters) {
	// TODO check if topic exists first
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

// TODO - temporary - will cleanup in a separate PR
type SampleMessage struct {
	MessageType string      `json:"message_type"`
	Data        interface{} `json:"data"`
}

// TODO - temporary - will cleanup in a separate PR
func SeedTopics(ctx context.Context, clients AllKafkaClusters) {
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

			cluster.Client.Produce(ctx, &kgo.Record{
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
