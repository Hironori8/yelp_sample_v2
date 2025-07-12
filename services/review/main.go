package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"review/handlers"
	"review/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Review service is running!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Review routes
	r.GET("/businesses/:id/reviews", handlers.GetBusinessReviews)
	r.POST("/businesses/:id/reviews", handlers.CreateReview)
	r.GET("/reviews", handlers.GetReviews)
	r.GET("/reviews/:id", handlers.GetReview)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Printf("Review service starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}