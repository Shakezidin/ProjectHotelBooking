package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Home is a handler for the user's homepage.
func Home(c *gin.Context) {
	var banners []models.Banner

	if err := Init.DB.Preload("Hotels").Where("available = ? AND active = ?", true, true).Find(&banners).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching banners"})
		return
	}
	city := c.Query("loc")
	if city == "" {
		c.JSON(400, gin.H{"error": "location query parameter is missing"})
		return
	}

	page := c.Query("page")
	limit := 10
	pageInt, _ := strconv.Atoi(page)

	skip := (pageInt - 1) * limit

	var hotels []models.Hotels
	var rooms []models.Rooms

	// Retrieve hotels based on the location and pagination
	if err := Init.DB.Preload("HotelCategory").Offset(skip).Limit(limit).Where("city = ? AND is_available = ? AND is_block = ? AND admin_approval = ?", city, true, false, true).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching hotels"})
		return
	}

	// Retrieve rooms for each hotel
	for i := range hotels {
		var hotelRooms []models.Rooms

		// Retrieve rooms for the current hotel
		if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("hotels_id = ? AND is_available = ? AND is_blocked = ? AND admin_approval = ?", hotels[i].ID, true, false, true).Find(&hotelRooms).Error; err != nil {
			c.JSON(400, gin.H{"error": "error while fetching rooms"})
			return
		}

		// Append the rooms of the current hotel to the rooms slice
		rooms = append(rooms, hotelRooms...)
	}

	c.JSON(200, gin.H{"Hotels": hotels, "Rooms": rooms, "Banners": banners})
}

// BannerShowing is a handler for displaying all banners.
func BannerShowing(c *gin.Context) {
	var banners []models.Banner

	if err := Init.DB.Preload("Hotels").Where("available = ? AND active = ?", true, true).Find(&banners).Error; err != nil {
		c.JSON(400, gin.H{"error": "banner retrieval error"})
		return
	}

	c.JSON(200, gin.H{"status": banners})
}
