package handlers

import (
	kafka "github.com/joswayski/kontext/api/clients/kafka"
	"github.com/joswayski/kontext/api/utils"
)

type Handler struct {
	KafkaClusters kafka.AllKafkaClusters
	Routes        []utils.Route
}
