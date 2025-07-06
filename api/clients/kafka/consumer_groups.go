package clients

import (
	"context"
	"fmt"
	"sort"
)

type ConsumerGroupState string

const (
	// No members in the group
	ConsumerGroupStateEmpty ConsumerGroupState = "Empty"
	// Group is preparing to rebalance (members joining/leaving)
	ConsumerGroupStatePreparingRebalance ConsumerGroupState = "PreparingRebalance"
	// Group is completing rebalance process
	ConsumerGroupStateCompletingRebalance ConsumerGroupState = "CompletingRebalance"
	// Group is stable with all members assigned partitions
	ConsumerGroupStateStable ConsumerGroupState = "Stable"
	// Group has no members and no offsets
	ConsumerGroupStateDead ConsumerGroupState = "Dead"
	// State cannot be determined
	ConsumerGroupStateUnknown ConsumerGroupState = "Unknown"
)

type ConsumerGroupInCluster struct {
	Name         string             `json:"name"`
	State        ConsumerGroupState `json:"state"`
	MembersCount int                `json:"members_count"`
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
			State:        ConsumerGroupState(group.State),
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
