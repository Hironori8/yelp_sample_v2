package main

import (
	"auth/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"auth/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	database.Connect()
	database.Migrate()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	// Authentication routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", handlers.GetMe)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	fmt.Printf("Auth service starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
