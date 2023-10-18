package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewHotelFacilities returns a list of all hotel facilities.
func ViewHotelFacilities(c *gin.Context) {
	var facilities []models.HotelAmenities
	if err := Init.DB.Find(&facilities).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"facilities": facilities,
	})
}

// AddHotelFacility adds a new hotel facility.
func AddHotelFacility(c *gin.Context) {
	var facility models.HotelAmenities

	if err := c.ShouldBindJSON(&facility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(facility)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error"})
		return
	}

	record := Init.DB.Create(&facility)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating facility",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "facility create success",
	})
}

// DeleteHotelFacility deletes a hotel facility by ID.
func DeleteHotelFacility(c *gin.Context) {
	facilityIDStr := c.DefaultQuery("id", "")
	if facilityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "facility ID is missing"})
		return
	}
	facilityID, err := strconv.Atoi(facilityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversion error"})
		return
	}

	if err := Init.DB.Where("facility_id = ?", uint(facilityID)).Delete(&models.HotelAmenities{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delete error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "delete success"})
}

// $2a$14$xYyRw2PmVhk/Ea/d9wAxguJa3NhxLNlFR2CriEMqzK4KZU6NBBbsm
