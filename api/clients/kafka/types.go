package clients

import (
	"github.com/joswayski/kontext/api/types"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// The client, admin client, and config for a cluster
type KafkaCluster struct {
	client      *kgo.Client
	adminClient *kadm.Client
	config      types.KafkaClusterConfig
}

// All clusters with their client, admin client, and config
type AllKafkaClusters map[string]KafkaCluster
