package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewRoomCancellation returns a list of room cancellation options.
func ViewRoomCancellation(c *gin.Context) {
	var cancellation []models.Cancellation
	if err := Init.DB.Find(&cancellation).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"cancellation": cancellation})
}

// AddCancellation adds a new cancellation option.
func AddCancellation(c *gin.Context) {
	var cancellation models.Cancellation

	if err := c.ShouldBindJSON(&cancellation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	validationErr := validate.Struct(cancellation)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	if record := Init.DB.Create(&cancellation); record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating cancellation",
			"error":   record.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cancellation created successfully"})
}

// DeleteCancellation deletes a cancellation option by ID.
func DeleteCancellation(c *gin.Context) {
	cancellationIDStr := c.Query("id")
	if cancellationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: Missing 'id' parameter"})
		return
	}
	cancellationID, err := strconv.Atoi(cancellationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}

	if err := Init.DB.Where("id = ?", uint(cancellationID)).Delete(&models.Cancellation{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete cancellation"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Deletion successful"})
}
