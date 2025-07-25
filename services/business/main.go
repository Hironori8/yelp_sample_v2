package main

import (
	"business/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"business/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Business service is running!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// Business routes
	r.GET("/businesses", handlers.SearchBusinesses)
	r.GET("/businesses/:id", handlers.GetBusiness)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Business service starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
