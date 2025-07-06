package clients

import (
	"context"
	"fmt"
	"sort"
)

type ConsumerGroupInCluster struct {
	Name string `json:"name"`
	/*
	 * Kafka Consumer Group States:
	 * - "Empty": No members in the group
	 * - "PreparingRebalance": Group is preparing to rebalance (members joining/leaving)
	 * - "CompletingRebalance": Group is completing rebalance process
	 * - "Stable": Group is stable with all members assigned partitions
	 * - "Dead": Group has no members and no offsets
	 * - "Unknown": State cannot be determined
	 */
	State        string `json:"state"`
	MembersCount int    `json:"members_count"`
}

type AllConsumerGroupsInCluster = []ConsumerGroupInCluster

func getConsumerGroupsInCluster(ctx context.Context, cluster KafkaCluster) (AllConsumerGroupsInCluster, error) {
	listedGroups, err := cluster.adminClient.ListGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list groups: %w", err)
	}

	describedGroups, err := cluster.adminClient.DescribeGroups(ctx, listedGroups.Groups()...)

	if err != nil {
		return nil, fmt.Errorf("could not describe consumer groups %w", err)
	}

	allConsumerGroups := make(AllConsumerGroupsInCluster, 0)
	for _, group := range describedGroups {
		cg := ConsumerGroupInCluster{
			Name:         group.Group,
			State:        group.State,
			MembersCount: len(group.Members),
		}
		allConsumerGroups = append(allConsumerGroups, cg)
	}

	// Sort alphabetically
	sort.Slice(allConsumerGroups, func(i, j int) bool {
		return allConsumerGroups[i].Name < allConsumerGroups[j].Name
	})

	return allConsumerGroups, nil
}
