package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	fmt.Println("Starting server on port ", port)

	r := gin.Default()
	r.GET("/", handlers.RootHandler)
	r.GET("/health", handlers.HealthHandler)
	r.NoRoute(handlers.NotFoundHandler(r))

	r.Run("0.0.0.0:" + port)
}
