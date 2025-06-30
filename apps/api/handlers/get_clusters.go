package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/apps/api/clients/kafka"
)

func (h *Handler) GetClustersHandler(c *gin.Context) {
	results := kafka.GetAllClusters(c.Request.Context(), h.KafkaClients)
	c.JSON(http.StatusOK, results)
}
