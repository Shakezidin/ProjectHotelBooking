package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewHotelCategories returns a list of hotel categories.
func ViewHotelCategories(c *gin.Context) {
	var categories []models.HotelCategory
	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// AddHotelCategory adds a new hotel category.
func AddHotelCategory(c *gin.Context) {
	var category models.HotelCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	validationErr := validate.Struct(category)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	record := Init.DB.Create(&category)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating category",
			"error":   record.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category created successfully"})
}

// DeleteHotelCategory deletes a hotel category by ID.
func DeleteHotelCategory(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("id", "")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'id' parameter"})
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	if err := Init.DB.Where("id = ?", uint(categoryID)).Delete(&models.HotelCategory{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Deletion successful"})
}
