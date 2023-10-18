package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewRoomCategories returns a list of all room categories.
func ViewRoomCategories(c *gin.Context) {
	var categories []models.RoomCategory
	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// AddRoomCategory adds a new room category.
func AddRoomCategory(c *gin.Context) {
	var category models.RoomCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
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
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category created successfully",
	})
}

// DeleteRoomCategory deletes a room category by ID.
func DeleteRoomCategory(c *gin.Context) {
	categoryIDStr := c.Query("id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID query parameter is missing"})
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversion error"})
		return
	}

	if err := Init.DB.Where("id = ?", uint(categoryID)).Delete(&models.RoomCategory{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Delete error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Delete success"})
}
