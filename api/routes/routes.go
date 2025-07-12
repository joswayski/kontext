package routes

import (
	"github.com/gin-gonic/gin"
	clients "github.com/joswayski/kontext/api/clients/kafka"
	"github.com/joswayski/kontext/api/handlers"
	"github.com/joswayski/kontext/api/utils"
)

func GetRoutes(kafkaClusters clients.AllKafkaClusters) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := handlers.Handler{
		KafkaClusters: kafkaClusters,
	}

	r.GET("", h.GetRootHandler)
	r.GET("/health", h.GetRootHandler)
	r.GET("/api/v1/clusters", h.GetClustersHandler)
	r.GET("/api/v1/clusters/:clusterId", h.GetClusterByIdHandler)
	r.GET("/test", h.Test)
	// r.GET("/api/v1/clusters/:clusterId/topics", h.getTopi)

	r.HandleMethodNotAllowed = true
	r.NoMethod(h.GetNoMethodHandler)
	r.NoRoute(h.GetNotFoundHandler)

	// Get all routes once after they're registered
	allRoutes := utils.GetAllRoutes(r)
	h.Routes = allRoutes

	return r
}
