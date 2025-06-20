package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler creates a 404 handler that shows available routes
func NotFoundHandler(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		var routes []string

		// Get all registered routes from the router
		for _, route := range router.Routes() {
			routes = append(routes, route.Method+" "+route.Path)
		}

		c.JSON(http.StatusNotFound, gin.H{
			"message":          fmt.Sprintf("Route '%s %s' not found - did you mean one of these instead?", c.Request.Method, c.Request.URL.Path),
			"available_routes": routes,
		})
	}
}
