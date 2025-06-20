package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/handlers"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	r.GET("/", handlers.RootHandler)
	r.GET("/health", handlers.HealthHandler)
	r.NoRoute(handlers.NotFoundHandler(r))
	r.NoMethod(handlers.NotFoundHandler(r))

	r.Run(":" + cfg.Server.Port)
}
