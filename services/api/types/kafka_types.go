package types

// TopicResponse represents the response structure for a Kafka topic
type TopicResponse struct {
	Name              string   `json:"name"`
	Partitions        int32    `json:"partitions"`
	ReplicationFactor int16    `json:"replication_factor"`
	Configs           []string `json:"configs,omitempty"`
	LeaderEpoch       int32    `json:"leader_epoch,omitempty"`
	IsInternal        bool     `json:"is_internal"`
}
