package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
	Cors     cors.Config
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout int
	GinMode         string
}

type DatabaseConfig struct {
	Url string
}

// Represents a single Kafka cluster to be monitored
type ClusterConfig struct {
	Name             string   // "PRODUCTION", set in .env with KAFKA_PRODUCTION_CLUSTER=kafka-production-1:9092,kafka-production-2:9092
	BootstrapServers []string // ["kafka-production-1:9092", "kafka-production-2:9092"]
}

type KafkaConfig struct {
	Clusters []ClusterConfig
}

func Load() *Config {
	envFiles := []string{
		".env",
		"../.env",
		"../../.env", // Root directory
	}

	var loaded bool
	for _, envFile := range envFiles {
		err := godotenv.Load(envFile)
		if err == nil {
			log.Printf("Loaded .env file from: %s", envFile)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Printf("Warning: No .env file found in any of the expected locations")
	}

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("API_PORT", "4000"),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 10),
			GinMode:         getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Url: getEnv("DATABASE_URL", ""),
		},
		Kafka: getKafkaConfig(),
		Cors:  getCorsConfig(),
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim whitespace
		values := strings.Split(value, ",")
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}
		return values
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getCorsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = getEnvAsSlice("FRONTEND_URLS", []string{"http://localhost:3000"})
	corsConfig.AllowCredentials = true
	return corsConfig
}

const KAFKA_PREFIX = "KAFKA_"
const KAFKA_SUFFIX = "_CLUSTER"

func getKafkaConfig() KafkaConfig {
	kafkaConfig := KafkaConfig{
		Clusters: []ClusterConfig{},
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Check if the environment variable key matches the pattern KAFKA_*_CLUSTER
		if strings.HasPrefix(key, KAFKA_PREFIX) && strings.HasSuffix(key, KAFKA_SUFFIX) {
			// Extract the cluster name by removing KAFKA_ prefix and _CLUSTER suffix
			clusterName := strings.TrimPrefix(key, KAFKA_PREFIX)
			clusterName = strings.TrimSuffix(clusterName, KAFKA_SUFFIX)

			// Split the comma-separated bootstrap servers
			bootstrapServers := strings.Split(value, ",")
			for i, v := range bootstrapServers {
				bootstrapServers[i] = strings.TrimSpace(v)
			}

			kafkaConfig.Clusters = append(kafkaConfig.Clusters, ClusterConfig{
				Name:             clusterName,
				BootstrapServers: bootstrapServers,
			})
		}
	}

	if len(kafkaConfig.Clusters) == 0 {
		log.Println("No KAFKA_*_CLUSTER variables found, using default localhost:9092")
		kafkaConfig.Clusters = append(kafkaConfig.Clusters, ClusterConfig{
			Name:             "DEFAULT",
			BootstrapServers: []string{"localhost:9092"},
		})
	}

	return kafkaConfig
}
