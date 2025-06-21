package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClusters(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello!",
	})
}
