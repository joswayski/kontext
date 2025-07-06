package types

// Global, app wide config
type KontextConfig struct {
	Port                string
	KafkaClusterConfigs AllKafkaClusterConfigs
}

// The config for a cluster including the broker URLs and the ID
type KafkaClusterConfig struct {
	BrokerURLs []string
	// Id of the cluster, taken from the broker URL(s), lowercased
	Id string
}

// All configs for all clusters
type AllKafkaClusterConfigs map[string]KafkaClusterConfig
