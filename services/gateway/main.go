package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key"
	}
	return secret
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user info in context and add headers for downstream services
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Request.Header.Set("X-User-ID", strconv.FormatUint(uint64(claims.UserID), 10))
		c.Request.Header.Set("X-User-Email", claims.Email)
		c.Next()
	}
}

func optionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(getJWTSecret()), nil
			})

			if err == nil && token.Valid {
				// Store user info in context and add headers for downstream services
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Request.Header.Set("X-User-ID", strconv.FormatUint(uint64(claims.UserID), 10))
				c.Request.Header.Set("X-User-Email", claims.Email)
				fmt.Printf("Gateway: Set X-User-ID header to %d\n", claims.UserID)
			}
		}
		// Continue regardless of token validity (optional auth)
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// Public routes
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

	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// Auth service routes (public)
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8084"
	}
	authURL, _ := url.Parse(authServiceURL)
	authProxy := httputil.NewSingleHostReverseProxy(authURL)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
		authGroup.POST("/login", func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
		authGroup.POST("/logout", func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
		authGroup.GET("/me", authMiddleware(), func(c *gin.Context) {
			authProxy.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Business service routes
	businessServiceURL := os.Getenv("BUSINESS_SERVICE_URL")
	if businessServiceURL == "" {
		businessServiceURL = "http://business-service:8081"
	}
	businessURL, _ := url.Parse(businessServiceURL)
	businessProxy := httputil.NewSingleHostReverseProxy(businessURL)

	// Public business routes
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

	// Review routes (public but with optional auth for logging)
	r.GET("/businesses/:id/reviews", optionalAuthMiddleware(), func(c *gin.Context) {
		// Debug: check if user info is in context
		if userID, exists := c.Get("user_id"); exists {
			fmt.Printf("Gateway: Found user_id in context: %v\n", userID)
		} else {
			fmt.Printf("Gateway: No user_id in context\n")
		}
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Protected review routes
	r.POST("/businesses/:id/reviews", authMiddleware(), func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/reviews", authMiddleware(), func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/reviews/:id", authMiddleware(), func(c *gin.Context) {
		reviewProxy.ServeHTTP(c.Writer, c.Request)
	})

	// Logging service routes
	loggingServiceURL := os.Getenv("LOGGING_SERVICE_URL")
	if loggingServiceURL == "" {
		loggingServiceURL = "http://logging-service:8083"
	}
	loggingURL, _ := url.Parse(loggingServiceURL)
	loggingProxy := httputil.NewSingleHostReverseProxy(loggingURL)

	// Protected logging routes
	r.POST("/logs/review-view", authMiddleware(), func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/logs/user/:user_id/history", authMiddleware(), func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/logs/business/:business_id/stats", authMiddleware(), func(c *gin.Context) {
		loggingProxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("API Gateway starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
