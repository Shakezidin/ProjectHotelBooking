package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Dashboard returns the admin dashboard statistics.
func Dashboard(c *gin.Context) {
	var (
		noOfUsers  int64
		noOfOwners int64
		noOfHotels int64
		noOfRooms  int64
	)

	// Count the number of users
	if err := initializer.DB.Model(&models.User{}).Count(&noOfUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Count the number of owners
	if err := initializer.DB.Model(&models.Owner{}).Count(&noOfOwners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Count the number of hotels
	if err := initializer.DB.Model(&models.Hotels{}).Count(&noOfHotels).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Count the number of rooms
	if err := initializer.DB.Model(&models.Rooms{}).Count(&noOfRooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Fetch admin revenue with owner details
	var adminRevenue models.Revenue
	if err := initializer.DB.Preload("Owner").Find(&adminRevenue).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fetching admin revenue error"})
		return
	}

	// Respond with the dashboard statistics
	c.JSON(http.StatusOK, gin.H{
		"users":  noOfUsers,
		"owners": noOfOwners,
		"hotels": noOfHotels,
		"rooms":  noOfRooms,
		"income": adminRevenue,
	})
}
