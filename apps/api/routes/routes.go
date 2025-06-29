package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/apps/api/handlers"
	services "github.com/joswayski/kontext/apps/api/services/kafka"
	"github.com/joswayski/kontext/apps/api/utils"
)

func GetRoutes(kafkaClients map[string]services.KafkaClients) *gin.Engine {
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
