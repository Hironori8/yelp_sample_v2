package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"yelp_sample_v2/database"
	"yelp_sample_v2/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.Migrate()
	
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Go server is running!",
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
	r.GET("/businesses/:id/reviews", handlers.GetBusinessReviews)
	r.POST("/businesses/:id/reviews", handlers.CreateReview)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}