package config

// Messages per second for a topic
type MessageRate = float64

// Topic configuration with message rate
type TopicConfig struct {
	MessageRate MessageRate `json:"message_rate"`
	// TODO sample message in here
}

// Topics within a cluster, keyed by topic name
type Topics = map[string]TopicConfig

// Cluster data containing topics
type ClusterData struct {
	Topics Topics `json:"topics"`
}

// Global Cluster Data where the key is the cluster ID
type GlobalClusterData = map[string]ClusterData

var ClusterDataConfig = GlobalClusterData{
	"users": {
		Topics: Topics{
			"user.created": {
				MessageRate: 0.1,
			},
			"user.updated": {
				MessageRate: 0.3,
			},
		},
	},
	"drivers": {
		Topics: Topics{
			"driver.onboarded": {
				MessageRate: 0.05,
			},
			"driver.activated": {
				MessageRate: 0.2,
			},
			"driver.deactivated": {
				MessageRate: 0.2,
			},
			"driver.location.updated": {
				MessageRate: 5.0,
			},
			"driver.rating.submitted": {
				MessageRate: 0.8,
			},
		},
	},

	"rides": {
		Topics: Topics{
			"ride.requested": {
				MessageRate: 2.0,
			},
			"ride.fare.calculated": {
				MessageRate: 2.0,
			},
			"ride.matched": {
				MessageRate: 1.8,
			},
			"ride.started": {
				MessageRate: 1.5,
			},
			"ride.completed": {
				MessageRate: 1.4,
			},
			"ride.cancelled": {
				MessageRate: 0.4,
			},
		},
	},

	"payments": {
		Topics: Topics{
			"payment.method.added": {
				MessageRate: 0.1,
			},
			"payment.method.removed": {
				MessageRate: 0.05,
			},
			"payment.initiated": {
				MessageRate: 1.4,
			},
			"payment.succeeded": {
				MessageRate: 1.3,
			},
			"payment.failed": {
				MessageRate: 0.1,
			},
			"refund.issued": {
				MessageRate: 0.05,
			},
		},
	},
}
