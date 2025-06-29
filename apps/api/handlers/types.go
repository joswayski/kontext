package handlers

import (
	"github.com/joswayski/kontext/apps/api/utils"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler struct {
	KafkaClients map[string]*kgo.Client
	Routes       []utils.Route
}
