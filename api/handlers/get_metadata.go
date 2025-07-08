package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/api/clients/kafka"
)

func (h *Handler) GetMetadata(c *gin.Context) {

	results, err := kafka.GetMetadata(c.Request.Context(), h.KafkaClusters)

	if err != nil {
		slog.Error("An error ocurred")
		c.JSON(http.StatusInternalServerError, results)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok", "data": results})

}
