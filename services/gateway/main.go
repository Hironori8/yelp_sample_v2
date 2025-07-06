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
	businessServiceURL := os.Getenv("BUSINESS_SERVICE_URL")
	if businessServiceURL == "" {
		businessServiceURL = "http://business-service:8081"
	}
	businessURL, _ := url.Parse(businessServiceURL)
	businessProxy := httputil.NewSingleHostReverseProxy(businessURL)

	r.GET("/businesses", func(c *gin.Context) {
		businessProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/businesses/:id", func(c *gin.Context) {
		businessProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Review service routes
	reviewServiceURL := os.Getenv("REVIEW_SERVICE_URL")
	if reviewServiceURL == "" {
		reviewServiceURL = "http://review-service:8082"
	}
	reviewURL, _ := url.Parse(reviewServiceURL)
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

	// Logging service routes
	loggingServiceURL := os.Getenv("LOGGING_SERVICE_URL")
	if loggingServiceURL == "" {
		loggingServiceURL = "http://logging-service:8083"
	}
	loggingURL, _ := url.Parse(loggingServiceURL)
	loggingProxy := httputil.NewSingleHostReverseProxy(loggingURL)

	r.POST("/logs/review-view", func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/logs/user/:user_id/history", func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/logs/business/:business_id/stats", func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("API Gateway starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
