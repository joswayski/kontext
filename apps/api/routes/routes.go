package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/apps/api/handlers"
	"github.com/twmb/franz-go/pkg/kgo"
)

func GetRoutes(kafkaClients map[string]*kgo.Client) *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := handlers.Handler{
		KafkaClients: kafkaClients,
	}

	r.GET("", handlers.RootHandler)
	r.GET("/health", handlers.RootHandler)
	r.GET("/api/v1/clusters", h.GetClustersHandler)

	return r
}
