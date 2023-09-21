package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewHotelCatagories(c *gin.Context) {
	var catagories []models.HotelCategory
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

func AddHotlecatagory(c *gin.Context) {
	var catagory models.HotelCategory

	if err := c.ShouldBindJSON(&catagory); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	validationErr := validate.Struct(catagory)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	record := Init.DB.Create(&catagory)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating catagory",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "catagory create Success",
	})
}

func DeleteHotelcatagory(c *gin.Context) {
	catagoryIDStr := c.DefaultQuery("catagoryid", "")
	if catagoryIDStr == "" {
		c.JSON(400, gin.H{"error": "catagoryid query parameter is missing"})
		return
	}
	catagoryID, err := strconv.Atoi(catagoryIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	if err := Init.DB.Where("catagory_id = ?", uint(catagoryID)).Delete(&models.HotelCategory{}).Error;err != nil {
		c.JSON(400, gin.H{"Error": "delete error"})
		return
	}
	c.JSON(200, gin.H{"status": "delete success"})
}
