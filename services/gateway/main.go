package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Gateway is running!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Business service routes
	businessURL, _ := url.Parse("http://business-service:8081")
	businessProxy := httputil.NewSingleHostReverseProxy(businessURL)

	r.GET("/businesses", func(c *gin.Context) {
		businessProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/businesses/:id", func(c *gin.Context) {
		businessProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Review service routes
	reviewURL, _ := url.Parse("http://review-service:8082")
	reviewProxy := httputil.NewSingleHostReverseProxy(reviewURL)

	r.GET("/businesses/:id/reviews", func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/businesses/:id/reviews", func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/reviews", func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/reviews/:id", func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("API Gateway starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
