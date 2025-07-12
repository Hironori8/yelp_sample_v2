package handlers

import (
	"net/http"
	"strconv"

	"business/database"
	"business/models"

	"github.com/gin-gonic/gin"
)

func SearchBusinesses(c *gin.Context) {
	var businesses []models.Business

	query := database.DB.Model(&models.Business{})

	if category := c.Query("category"); category != "" {
		query = query.Where("category ILIKE ?", "%"+category+"%")
	}

	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if location := c.Query("location"); location != "" {
		query = query.Where("address ILIKE ?", "%"+location+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&businesses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search businesses"})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

func GetBusiness(c *gin.Context) {
	id := c.Param("id")

	var business models.Business
	if err := database.DB.First(&business, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Business not found"})
		return
	}

	c.JSON(http.StatusOK, business)
}
