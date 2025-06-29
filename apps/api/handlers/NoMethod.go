package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/apps/api/utils"
)

func (h *Handler) GetNoMethodHandler(c *gin.Context) {
	var alternatives []utils.Route

	currentPath := c.Request.URL.Path
	currentMethod := c.Request.Method

	// Check if path exists but method doesn't (405 Method Not Allowed)
	for _, route := range h.Routes {
		if route.Path == currentPath && route.Method != currentMethod {
			alternatives = append(alternatives, route)
		}
	}

	// If we found alternatives, it's a 405 Method Not Allowed
	if len(alternatives) > 0 {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": fmt.Sprintf("Method '%s' not allowed for path '%s'. Did you mean one of these?", currentMethod, currentPath),
			"routes":  alternatives,
		})
		return
	}
}
