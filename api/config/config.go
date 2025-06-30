package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type KontextConfig struct {
	Port          string
	KafkaClusters map[string]KafkaClusterConfig
}

type KafkaClusterConfig struct {
	BrokerURLs []string
	// Id of the cluster, taken from the broker URL(s), lowercased
	Id string
}

func GetConfig() *KontextConfig {

	err := godotenv.Load("../../.env")
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to load .env file from root: %v", err))
	}

	return &KontextConfig{
		Port:          getPort(),
		KafkaClusters: getKafkaClusters(),
	}
}

const envPort = "API_PORT"
const defaultPort = "4000"

func getPort() string {
	port := os.Getenv(envPort)

	if port == "" {
		slog.Warn(fmt.Sprintf("No %s environment variable found, using default port %s", envPort, defaultPort))
		return defaultPort
	}

	return port
}

const brokerUrlPrefix = "KAFKA_"
const brokerUrlSuffix = "_BROKER_URL"

func getKafkaClusters() map[string]KafkaClusterConfig {
	clusters := make(map[string]KafkaClusterConfig)
	envs := os.Environ()

	for _, env := range envs {

		var key, value string

		parts := strings.Split(env, "=")
		if len(parts) == 2 {
			key = parts[0]
			value = parts[1]
		}

		if strings.HasPrefix(key, brokerUrlPrefix) && strings.HasSuffix(key, brokerUrlSuffix) {
			clusterId := strings.TrimPrefix(key, brokerUrlPrefix)
			clusterId = strings.TrimSuffix(clusterId, brokerUrlSuffix)
			clusterId = strings.ToLower(clusterId)
			urls := strings.Split(value, ",")
			if len(urls) == 0 || value == "" {
				slog.Warn(fmt.Sprintf("A Kafka broker key was set (%s), but no URLs were provided", key))
				continue
			}
			clusters[clusterId] = KafkaClusterConfig{
				Id:         clusterId,
				BrokerURLs: strings.Split(value, ","),
			}
		}
	}

	if len(clusters) == 0 {
		log.Fatal("No Kafka clusters found in environment variables! Make sure to set the KAFKA_<CLUSTER_ID>_BROKER_URL environment variable for each cluster.")
	} else {
		s := ""
		if len(clusters) != 1 {
			s = "s"
		}
		slog.Info(fmt.Sprintf("Found %d cluster%s in the env config!", len(clusters), s))
		idx := 1
		for id := range clusters {
			slog.Info(fmt.Sprintf("%d. %s", idx, id))
			idx++
		}
	}

	return clusters
}
