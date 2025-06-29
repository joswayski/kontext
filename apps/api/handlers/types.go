package handlers

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Handler struct {
	KafkaClients map[string]*kgo.Client
}
