package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"

	"github.com/brianvoe/gofakeit"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type TopicsInCluster struct {
	Name            string `json:"name"`
	PartitionsCount int    `json:"partitions_count"`
}

type AllTopicsInCluster = []TopicsInCluster

func Test(ctx context.Context, clients AllKafkaClusters) (kadm.DescribedGroups, error) {
	v, err := clients["production"].adminClient.DescribeGroups(ctx)

	if err != nil {
		return nil, err
	}

	return v, nil
}
func getTopicsInCluster(ctx context.Context, cluster KafkaCluster) (AllTopicsInCluster, error) {
	topics, err := cluster.adminClient.ListTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve toics in cluster %s", err.Error())
	}

	allTopics := make(AllTopicsInCluster, 0)
	for _, topic := range topics {

		allTopics = append(allTopics, TopicsInCluster{
			Name:            topic.Topic,
			PartitionsCount: int(len(topic.Partitions)),
		})
	}

	// Sort alphabetically
	sort.Slice(allTopics, func(i, j int) bool {
		return allTopics[i].Name < allTopics[j].Name
	})
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
