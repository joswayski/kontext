package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joswayski/kontext/config"
	"github.com/joswayski/kontext/handlers"
	"github.com/joswayski/kontext/services"
)

func main() {
	cfg := config.Load()
	gin.SetMode(cfg.Server.GinMode)

	kafkaService := services.NewKafkaService(cfg)
	defer kafkaService.Close()

	r := gin.Default()
	r.Use(cors.New(cfg.Cors))

	// Intentional
	r.GET("/", handlers.RootHandler)
	r.GET("/health", handlers.HealthHandler)

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("", handlers.RootHandler)
		apiV1.GET("/health", handlers.HealthHandler)
		apiV1.GET("/clusters", func(ctx *gin.Context) {
			handlers.GetClusters(ctx, cfg, kafkaService)
		})
	}

	r.NoRoute(handlers.NotFoundHandler(r))
	r.NoMethod(handlers.NotFoundHandler(r))

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
