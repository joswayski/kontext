package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/api/clients/kafka"
)

func (h *Handler) GetTopicsByCluster(c *gin.Context) {

	cid := c.Param("clusterId")
	if cid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No clusterId provided"})
		return
	}

	_, exists := h.KafkaClusters[cid]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("The provided clusterId of '%s' not found in the environment", cid)})
		return
	}

	results, _ := kafka.GetTopicsByCluster(c.Request.Context(), h.KafkaClusters, cid)

	c.JSON(http.StatusOK, results)
}
