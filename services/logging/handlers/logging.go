package handlers

import (
	"net/http"
	"strconv"
	"time"

	"logging/cassandra"
	"logging/models"

	"github.com/gin-gonic/gin"
	"github.com/scylladb/gocqlx/v2/qb"
)

var reviewViewLogsTable = "yelp_logs.review_view_logs"

func LogReviewView(c *gin.Context) {
	var req models.LogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.LogResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// TODO: 認証システム実装後はPostgreSQLユーザーIDとの適切なマッピングを実装
	userID := 1

	// Create log entry
	log := models.ReviewViewLog{
		UserID:     userID,
		BusinessID: req.BusinessID,
		ReviewID:   req.ReviewID,
		ViewedAt:   time.Now(),
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
	}

	// Insert into Cassandra
	stmt, names := qb.Insert(reviewViewLogsTable).Columns(
		"user_id", "business_id", "review_id", "viewed_at",
		"ip_address", "user_agent",
	).ToCql()

	q := cassandra.Session.Query(stmt, names).BindStruct(&log)
	if err := q.ExecRelease(); err != nil {
		c.JSON(http.StatusInternalServerError, models.LogResponse{
			Success: false,
			Message: "Failed to log review view: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.LogResponse{
		Success: true,
		Message: "Review view logged successfully",
	})
}

func GetUserViewHistory(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	// Query user's view history
	stmt, names := qb.Select(reviewViewLogsTable).
		Where(qb.Eq("user_id")).
		OrderBy("viewed_at", qb.DESC).
		Limit(100).
		ToCql()

	q := cassandra.Session.Query(stmt, names).BindMap(map[string]interface{}{
		"user_id": userID,
	})

	var logs []models.ReviewViewLog
	if err := q.SelectRelease(&logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch view history"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func GetBusinessViewStats(c *gin.Context) {
	businessIDStr := c.Param("business_id")
	if businessIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	// Note: This is a simplified example. In a real scenario, you'd need to create
	// additional tables or use materialized views for efficient business-based queries
	// since Cassandra doesn't support secondary index queries efficiently

	c.JSON(http.StatusOK, gin.H{
		"message":     "Business view stats endpoint - would require additional table design for efficient querying",
		"business_id": businessIDStr,
	})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "logging",
	})
}
