package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/pkg/kafka"
)

func (h *Handler) GetClustersHandler(c *gin.Context) {
	results := kafka.GetMetadataForAllClusters(c.Request.Context(), h.KafkaClusters)
	c.JSON(http.StatusOK, results)
}
