package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/services"
)

func GetTopics(ctx *gin.Context, cfg *config.Config, kafkaService *services.KafkaService) {
	clusterId := ctx.Param("clusterId")
	if clusterId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cluster ID is required"})
		return
	}

	// Normalize cluster ID to match the format used in config
	clusterId = strings.ToUpper(clusterId)

	// Get topics from the Kafka service
	topics, err := kafkaService.GetTopics(clusterId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, topics)
}
