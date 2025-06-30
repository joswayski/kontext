package handlers

import (
	kafka "github.com/joswayski/kontext/apps/api/clients/kafka"
	"github.com/joswayski/kontext/apps/api/utils"
)

type Handler struct {
	KafkaClusters map[string]kafka.KafkaCluster
	Routes        []utils.Route
}
