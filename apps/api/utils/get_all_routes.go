package utils

import "github.com/gin-gonic/gin"

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func GetAllRoutes(r *gin.Engine) []Route {
	currentRoutes := r.Routes()
	result := make([]Route, 0)

	for _, route := range currentRoutes {
		result = append(result, Route{
			Method: route.Method,
			Path:   route.Path,
		})
	}

	return result
}
