package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Saul Goodman! This is the root route, you probably want one of these other ones:",
		"routes":  h.Routes,
	})
}
