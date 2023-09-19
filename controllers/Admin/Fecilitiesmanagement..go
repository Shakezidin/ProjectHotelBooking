package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewHotelFecilities(c *gin.Context) {
	var fecilities []models.HotelAmenities
	if err := Init.DB.Find(&fecilities).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	// var hotel models.Hotel
	c.JSON(200, gin.H{
		"fecilities": fecilities,
		// "hotel": hotel,
	})
}

func AddHotelfecilility(c *gin.Context) {
	var fecility models.HotelAmenities

	if err := c.ShouldBindJSON(&fecility); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	validationErr := validate.Struct(fecility)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	record := Init.DB.Create(&fecility)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating fecility",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "fecility create Success",
	})
}

func DeleteHotelfecility(c *gin.Context) {
	fecilityIDStr := c.DefaultQuery("hotelid", "")
	if fecilityIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	fecilityID, err := strconv.Atoi(fecilityIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	if err := Init.DB.Where("fecilityid = ?", uint(fecilityID)).Delete(&models.HotelAmenities{}); err != nil {
		c.JSON(400, gin.H{"Error": "delete error"})
		return
	}
	c.JSON(200, gin.H{"status": "delete success"})
}
