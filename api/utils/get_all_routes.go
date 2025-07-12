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
	if route == "/" {
		return "Returns basic endpoint info. You're here now!"
	}

	if route == "/api/v1/clusters" {
		return "Returns the cluster IDs along with connectivity information"
	}

	if route == "/api/v1/clusters/:clusterId" {
		return "Returns the topics and consumers of the cluster"
	}

	if route == "/api/v1/clusters/:clusterId/topics" {
		return "Returns the topics in the cluster, along with the consumer groups"
	}

	if route == "/api/v1/health" {
		return "Simple healthcheck!"
	}

	return "No description found! Would you be so kind to make a PR here? https://github.com/joswayski/kontext/blob/main/api/utils/get_all_routes.go"
}
