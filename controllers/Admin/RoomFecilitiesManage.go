package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewRoomFacilities returns a list of all room facilities.
func ViewRoomFacilities(c *gin.Context) {
	var facilities []models.RoomFacilities
	if err := Init.DB.Find(&facilities).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"facilities": facilities})
}

// AddRoomFacility adds a new room facility.
func AddRoomFacility(c *gin.Context) {
	var facility models.RoomFacilities

	if err := c.ShouldBindJSON(&facility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Binding error"})
		return
	}

	validationErr := validate.Struct(facility)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
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
		"message": "Facility created successfully",
	})
}

// DeleteRoomFacility deletes a room facility by ID.
func DeleteRoomFacility(c *gin.Context) {
	facilityIDStr := c.Query("id")
	if facilityIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Facility ID query parameter is missing"})
		return
	}
	facilityID, err := strconv.Atoi(facilityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Conversion error"})
		return
	}

	if err := Init.DB.Where("id = ?", uint(facilityID)).Delete(&models.RoomFacilities{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Delete error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Delete success"})
}
