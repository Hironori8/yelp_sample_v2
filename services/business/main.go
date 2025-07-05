package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"business/handlers"
	"yelp_sample_v2/database"

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