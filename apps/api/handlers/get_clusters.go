package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClusterStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (h *Handler) GetClustersHandler(c *gin.Context) {
	results := make(map[string]ClusterStatus)

	slog.Info(fmt.Sprintf("Checking %d kafka clients", len(h.KafkaClients)))
	for k, v := range h.KafkaClients {
		slog.Info(fmt.Sprintf("Checking %s", k))
		ping := v.Ping(c.Request.Context())
		healthy := ping == nil
		status := "connected"
		message := "Saul Goodman"
		if !healthy {
			status = "error"
			message = ping.Error()
		}
		results[k] = ClusterStatus{
			Status:  status,
			Message: message,
		}

	}
	c.JSON(http.StatusOK, results)
}
