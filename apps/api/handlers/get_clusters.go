package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ClusterStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (h *Handler) GetClustersHandler(c *gin.Context) {
	results := make(map[string]ClusterStatus)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for clusterName, kafkaClient := range h.KafkaClients {
		wg.Add(1)
		go func(name string, kafkaClient *kgo.Client) {
			defer wg.Done()
			ping := kafkaClient.Ping(c.Request.Context())
			healthy := ping == nil
			status := "connected"
			message := "Saul Goodman"
			if !healthy {
				status = "error"
				message = ping.Error()
			}

			mu.Lock()
			results[clusterName] = ClusterStatus{
				Status:  status,
				Message: message,
			}
			mu.Unlock()

		}(clusterName, kafkaClient)

	}

	wg.Wait()
	c.JSON(http.StatusOK, results)
}
