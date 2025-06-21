package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/services"
)

// ClusterResponse represents the response structure for a Kafka cluster
type ClusterResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	BootstrapServers string `json:"bootstrapServers"`
	Status           string `json:"status"`
}

func GetClusters(ctx *gin.Context, cfg *config.Config, kafkaService *services.KafkaService) {
	// Get cluster info from the Kafka service
	clusterInfos := kafkaService.GetClusterInfo(cfg)

	var clusters []ClusterResponse

	// Map ClusterInfo to ClusterResponse
	for _, info := range clusterInfos {
		clusters = append(clusters, ClusterResponse{
			ID:               info.ID,
			Name:             info.Name,
			BootstrapServers: info.BootstrapServers,
			Status:           info.Status,
		})
	}

	ctx.JSON(http.StatusOK, clusters)
}
