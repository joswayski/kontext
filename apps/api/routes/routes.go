package routes

import (
	"github.com/gin-gonic/gin"
	clients "github.com/joswayski/kontext/apps/api/clients/kafka"
	"github.com/joswayski/kontext/apps/api/handlers"
	"github.com/joswayski/kontext/apps/api/utils"
)

func GetRoutes(kafkaClients map[string]clients.KafkaClients) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := handlers.Handler{
		KafkaClients: kafkaClients,
	}

	r.GET("", h.GetRootHandler)
	r.GET("/health", h.GetRootHandler)
	r.GET("/api/v1/clusters", h.GetClustersHandler)

	r.HandleMethodNotAllowed = true
	r.NoMethod(h.GetNoMethodHandler)
	r.NoRoute(h.GetNotFoundHandler)

	// Get all routes once after they're registered
	allRoutes := utils.GetAllRoutes(r)
	h.Routes = allRoutes

	return r
}
