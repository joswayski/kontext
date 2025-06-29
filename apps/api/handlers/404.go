package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetNotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "The route you're looking for was not found. Perhaps you wanted one of these?", "routes": h.Routes})

}
