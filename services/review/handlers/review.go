package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"yelp_sample_v2/database"
	"yelp_sample_v2/models"

	"github.com/gin-gonic/gin"
)

func GetBusinessReviews(c *gin.Context) {
	id := c.Param("id")

	var business models.Business
	if err := database.DB.First(&business, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	var reviews []models.Review

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	if err := database.DB.Where("business_id = ?", id).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	// Log review views asynchronously
	go logReviewViews(c, id, reviews)

	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context) {
	businessID := c.Param("id")

	var business models.Business
	if err := database.DB.First(&business, businessID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	var request struct {
		Rating int    `json:"rating" binding:"required,min=1,max=5"`
		Text   string `json:"text"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	businessIDInt, _ := strconv.Atoi(businessID)

	review := models.Review{
		BusinessID: uint(businessIDInt),
		UserID:     1,
		Rating:     request.Rating,
		Text:       request.Text,
	}

	if err := database.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	updateBusinessRating(uint(businessIDInt))

	c.JSON(http.StatusCreated, review)
}

func GetReviews(c *gin.Context) {
	var reviews []models.Review

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	if err := database.DB.Preload("User").
		Preload("Business").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func GetReview(c *gin.Context) {
	id := c.Param("id")

	var review models.Review
	if err := database.DB.Preload("User").
		Preload("Business").
		First(&review, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Log single review view asynchronously
	go logSingleReviewView(c, strconv.Itoa(int(review.BusinessID)), review.ID)

	c.JSON(http.StatusOK, review)
}

func updateBusinessRating(businessID uint) {
	var avgRating float32
	var count int64

	database.DB.Model(&models.Review{}).
		Where("business_id = ?", businessID).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Row().Scan(&avgRating, &count)

	database.DB.Model(&models.Business{}).
		Where("id = ?", businessID).
		Updates(map[string]interface{}{
			"rating":       avgRating,
			"review_count": count,
		})
}

// logReviewViews sends review view logs to the logging service
func logReviewViews(c *gin.Context, businessID string, reviews []models.Review) {
	for _, review := range reviews {
		logSingleReviewView(c, businessID, review.ID)
	}
}

// logSingleReviewView sends a single review view log to the logging service
func logSingleReviewView(c *gin.Context, businessID string, reviewID uint) {
	loggingServiceURL := os.Getenv("LOGGING_SERVICE_URL")
	if loggingServiceURL == "" {
		loggingServiceURL = "http://logging-service:8083"
	}

	// TODO: 認証システム実装後の修正予定
	// 開発用: 固定ユーザーID（PostgreSQLのusers.id=1）
	userID := 1

	businessIDInt, _ := strconv.Atoi(businessID)

	logData := map[string]interface{}{
		"user_id":     userID,
		"business_id": businessIDInt,
		"review_id":   int(reviewID),
		"ip_address":  c.ClientIP(),
		"user_agent":  c.GetHeader("User-Agent"),
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		return // Fail silently for logging
	}

	// Send to logging service (fire and forget)
	go func() {
		resp, err := http.Post(
			loggingServiceURL+"/logs/review-view",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err == nil {
			resp.Body.Close()
		}
	}()
}
