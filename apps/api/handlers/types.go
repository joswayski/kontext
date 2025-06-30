package handlers

import (
	services "github.com/joswayski/kontext/apps/api/clients/kafka"
	"github.com/joswayski/kontext/apps/api/utils"
)

type Handler struct {
	KafkaClients map[string]services.KafkaClients
	Routes       []utils.Route
}
