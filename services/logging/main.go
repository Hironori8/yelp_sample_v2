package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"logging/cassandra"
	"logging/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to Cassandra
	if err := cassandra.Connect(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()

	// Run migrations
	if err := cassandra.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Setup Gin
	r := gin.Default()

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Logging service is running!",
		})
	})

	r.GET("/health", handlers.HealthCheck)

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// Logging endpoints
	r.POST("/logs/review-view", handlers.LogReviewView)
	r.GET("/logs/user/:user_id/history", handlers.GetUserViewHistory)
	r.GET("/logs/business/:business_id/stats", handlers.GetBusinessViewStats)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Setup graceful shutdown
	go func() {
		fmt.Printf("Logging service starting on port %s\n", port)
		if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Logging service shutting down...")
}
