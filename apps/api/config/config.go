package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func GetConfig() *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		slog.Warn("No .env file found at project root!")
	}

	return &Config{
		Port: getPort(),
	}
}

const envPort = "API_PORT"
const defaultPort = "8080"

func getPort() string {
	port := os.Getenv(envPort)

	if port == "" {
		slog.Warn(fmt.Sprintf("No %s environment variable found, using default port %s", envPort, defaultPort))
		return defaultPort
	}

	return port
}
