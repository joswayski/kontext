package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/services"
)

func GetTopics(ctx *gin.Context, cfg *config.Config, kafkaService *services.KafkaService) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "cluster ID is required"})
		return
	}

	// Get topics from the Kafka service
	topics, err := kafkaService.GetTopics(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the topics array (empty if no topics found)
	ctx.JSON(http.StatusOK, topics)
}
