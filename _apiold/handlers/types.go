package handlers

import (
	"github.com/joswayski/kontext/api/utils"
	kafka "github.com/joswayski/kontext/pkg/kafka"
)

type Handler struct {
	KafkaClusters kafka.AllKafkaClusters
	Routes        []utils.Route
}
