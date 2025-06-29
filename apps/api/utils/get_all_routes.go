package utils

import "github.com/gin-gonic/gin"

type Route struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
}

func GetAllRoutes(r *gin.Engine) []Route {
	currentRoutes := r.Routes()
	result := make([]Route, 0)

	for _, route := range currentRoutes {
		result = append(result, Route{
			Method:      route.Method,
			Path:        route.Path,
			Description: getRouteDescription(route.Path),
		})
	}

	return result
}

func getRouteDescription(route string) string {
	result := ""
	if route == "/" {
		result = "Returns basic endpoint info."
	}

	if route == "/api/v1/clusters" {
		result = "Returns the cluster IDs along with connectivity information"
	}

	return result
}
