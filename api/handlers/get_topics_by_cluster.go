package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	kafka "github.com/joswayski/kontext/api/clients/kafka"
)

func (h *Handler) Test(c *gin.Context) {
	results, _ := kafka.Test(c.Request.Context(), h.KafkaClusters)

	c.JSON(http.StatusOK, results)
}
