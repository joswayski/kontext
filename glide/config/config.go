package config

import (
	msg "github.com/joswayski/kontext/glide/config/messages"
)

// Messages per second for a topic
type MessageRate = float64

// Topic configuration with message rate
type TopicConfig struct {
	MessageRate   MessageRate        `json:"message_rate"`
	CreateMessage func() interface{} `json:"-"`
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
				MessageRate:   0.1,
				CreateMessage: msg.GenerateUserCreatedMessage,
			},
			"user.updated": {
				MessageRate:   0.3,
				CreateMessage: msg.GenerateUserUpdatedMessage,
			},
		},
	},
	"drivers": {
		Topics: Topics{
			"driver.onboarded": {
				MessageRate:   0.05,
				CreateMessage: msg.GenerateDriverOnboardedMessage,
			},
			"driver.activated": {
				MessageRate:   0.2,
				CreateMessage: msg.GenerateDriverActivatedMessage,
			},
			"driver.deactivated": {
				MessageRate:   0.2,
				CreateMessage: msg.GenerateDriverDeactivatedMessage,
			},
			"driver.location.updated": {
				MessageRate:   5.0,
				CreateMessage: msg.GenerateDriverLocationUpdatedMessage,
			},
			"driver.rating.submitted": {
				MessageRate:   0.8,
				CreateMessage: msg.GenerateDriverRatingSubmittedMessage,
			},
		},
	},

	"rides": {
		Topics: Topics{
			"ride.requested": {
				MessageRate:   2.0,
				CreateMessage: msg.GenerateRideRequestedMessage,
			},
			"ride.fare.calculated": {
				MessageRate:   2.0,
				CreateMessage: msg.GenerateRideFareCalculatedMessage,
			},
			"ride.matched": {
				MessageRate:   1.8,
				CreateMessage: msg.GenerateRideMatchedMessage,
			},
			"ride.started": {
				MessageRate:   1.5,
				CreateMessage: msg.GenerateRideStartedMessage,
			},
			"ride.completed": {
				MessageRate:   1.4,
				CreateMessage: msg.GenerateRideCompletedMessage,
			},
			"ride.cancelled": {
				MessageRate:   0.4,
				CreateMessage: msg.GenerateRideCancelledMessage,
			},
		},
	},

	"payments": {
		Topics: Topics{
			"payment.method.added": {
				MessageRate:   0.1,
				CreateMessage: msg.GeneratePaymentMethodAddedMessage,
			},
			"payment.method.removed": {
				MessageRate:   0.05,
				CreateMessage: msg.GeneratePaymentMethodRemovedMessage,
			},
			"payment.initiated": {
				MessageRate:   1.4,
				CreateMessage: msg.GeneratePaymentInitiatedMessage,
			},
			"payment.succeeded": {
				MessageRate:   1.3,
				CreateMessage: msg.GeneratePaymentSucceededMessage,
			},
			"payment.failed": {
				MessageRate:   0.1,
				CreateMessage: msg.GeneratePaymentFailedMessage,
			},
			"refund.issued": {
				MessageRate:   0.05,
				CreateMessage: msg.GenerateRefundIssuedMessage,
			},
		},
	},
}
