package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewRoomCancellation(c *gin.Context) {
	var cancellation []models.Cancellation
	if err := Init.DB.Find(&cancellation).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	// var hotel models.Hotel
	c.JSON(200, gin.H{
		"fecilities": cancellation,
		// "hotel": hotel,
	})
}

func Addcancellation(c *gin.Context) {
	var cancellation models.Cancellation

	if err := c.ShouldBindJSON(&cancellation); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	validationErr := validate.Struct(cancellation)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	record := Init.DB.Create(&cancellation)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating cancellation",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cancellation create Success",
	})
}

func Deletecancellation(c *gin.Context) {
	cancellationIDStr := c.DefaultQuery("cancellationid", "")
	if cancellationIDStr == "" {
		c.JSON(400, gin.H{"error": "catagoryid query parameter is missing"})
		return
	}
	cancellationID, err := strconv.Atoi(cancellationIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	if err := Init.DB.Where("cancellation_id = ?", uint(cancellationID)).Delete(&models.Cancellation{}); err != nil {
		c.JSON(400, gin.H{"Error": "delete error"})
		return
	}
	c.JSON(200, gin.H{"status": "delete success"})
}
