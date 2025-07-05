package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/api/clients/kafka"
)

func (h *Handler) GetClusterByIdHandler(c *gin.Context) {
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

	results, err := kafka.GetClusterById(c.Request.Context(), cid, h.KafkaClusters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
