package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/brianvoe/gofakeit"
	"github.com/twmb/franz-go/pkg/kgo"
)

type TopicInCluster struct {
	Name            string `json:"name"`
	PartitionsCount int    `json:"partitions_count"`
}

type AllTopicsInCluster = []TopicInCluster

type DetailedTopic struct {
	TopicInCluster
	ConsumerGroups []string `json:"consumer_groups"`
}

type GetTopicsByClusterResult struct {
	Topics []DetailedTopic `json:"topics"`
}

func GetTopicsByCluster(ctx context.Context, clients AllKafkaClusters, clusterId string) (GetTopicsByClusterResult, error) {
	allTopics, err := getTopicsInCluster(ctx, clients[clusterId])
	if err != nil {
		return GetTopicsByClusterResult{}, fmt.Errorf("unable to retrieve topics %s", err.Error())
	}

	topicsAndConsumerGroups, cgErr := getConsumerGroupsForAllTopics(ctx, clients[clusterId])

	if cgErr != nil {
		return GetTopicsByClusterResult{}, fmt.Errorf("unable to retrieve consumer groups for topics %s", err.Error())
	}

	finalTopicList := make([]DetailedTopic, 0)

	for _, topic := range allTopics {
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

type AllConsumerGroupsInTopics = map[string][]string

func getConsumerGroupsForAllTopics(ctx context.Context, cluster KafkaCluster) (AllConsumerGroupsInTopics, error) {
	allGroups, err := cluster.adminClient.DescribeGroups(ctx)

	if err != nil {
		return nil, err
	}

	allTopics := make(map[string][]string)

	for _, group := range allGroups {
		topics := group.AssignedPartitions().Topics()

		for _, topic := range topics {
			allTopics[topic] = append(allTopics[topic], group.Group)
		}
	}

	return allTopics, nil
}

func getTopicsInCluster(ctx context.Context, cluster KafkaCluster) (AllTopicsInCluster, error) {
	topics, err := cluster.adminClient.ListTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve toics in cluster %s", err.Error())
	}

	sortedTopics := topics.Sorted()

	allTopics := make(AllTopicsInCluster, 0)
	for _, topic := range sortedTopics {

		allTopics = append(allTopics, TopicInCluster{
			Name:            topic.Topic,
			PartitionsCount: int(len(topic.Partitions)),
		})
	}

	return allTopics, nil
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
