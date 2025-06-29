package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/apps/api/services/kafka"
)

func (h *Handler) GetClustersHandler(c *gin.Context) {
	results := kafka.GetClusterStatuses(c.Request.Context(), h.KafkaClients)
	c.JSON(http.StatusOK, results)
}
