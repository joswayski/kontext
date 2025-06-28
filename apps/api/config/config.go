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

const defaultPort = "8080"

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		slog.Warn(fmt.Sprintf("No PORT environment variable found, using default port %s", defaultPort))
		return defaultPort
	}

	return port
}
