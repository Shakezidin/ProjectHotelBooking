package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewRoomCatagory(c *gin.Context) {
	var catagories []models.RoomCategory
	if err := Init.DB.Find(&catagories).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	// var hotel models.Hotel
	c.JSON(200, gin.H{
		"fecilities": catagories,
		// "hotel": hotel,
	})
}

func AddRoomCatagory(c *gin.Context) {
	var catagories models.RoomCategory

	if err := c.ShouldBindJSON(&catagories); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	validationErr := validate.Struct(catagories)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	record := Init.DB.Create(&catagories)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating catagories",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "catagories create Success",
	})
}

func DeleteRoomCatagories(c *gin.Context) {
	catatagoryIDStr := c.DefaultQuery("id", "")
	if catatagoryIDStr == "" {
		c.JSON(400, gin.H{"error": "catagoryId query parameter is missing"})
		return
	}
	catagoryID, err := strconv.Atoi(catatagoryIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	if err := Init.DB.Where("id = ?", uint(catagoryID)).Delete(&models.RoomCategory{}).Error; err != nil {
		c.JSON(400, gin.H{"Error": "delete error"})
		return
	}
	c.JSON(200, gin.H{"status": "delete success"})
}
