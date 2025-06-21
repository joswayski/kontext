package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout int
}

type DatabaseConfig struct {
	Url string
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Brokers       []string
	AdminURL      string
	SchemaURL     string
	ConsumerGroup string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}
	log.Println("Loaded .env file")

	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "4000"),
			ShutdownTimeout: getEnvAsInt("SHUTDOWN_TIMEOUT", 10),
		},
		Database: DatabaseConfig{
			Url: getEnv("DATABASE_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers:       getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			AdminURL:      getEnv("REDPANDA_ADMIN_URL", "http://localhost:9644"),
			SchemaURL:     getEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "kontext-consumer"),
		},
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
