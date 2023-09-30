package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> User HomePage <<<<<<<<<<<<<<<<<<<<<<<<<<
func UserHome(c *gin.Context) {
	var banner []models.Banner

	if err := Init.DB.Preload("Hotels").Where("available = ? AND active = ?", true, true).Find(&banner).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching banner"})
		return
	}
	city := c.DefaultQuery("location", "")
	if city == "" {
		c.JSON(400, gin.H{"error": "location query parameter is missing"})
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := 10
	pageInt := 1
	if p, err := strconv.Atoi(page); err == nil {
		pageInt = p
	}

	skip := (pageInt - 1) * limit

	var hotels []models.Hotels
	var rooms []models.Rooms

	// Retrieve hotels based on the location and pagination
	if err := Init.DB.Preload("HotelCategory").Offset(skip).Limit(limit).Where("city = ? AND is_available = ? AND isblock = ? AND adminapproval = ?", city, true, false, true).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching hotels"})
		return
	}

	// Retrieve rooms for each hotel
	for i := range hotels {
		var hotelRooms []models.Rooms

		// Retrieve rooms for the current hotel
		if err := Init.DB.Where("hotel_id = ? AND is_available = ? AND isblocked = ? adminapproval = ?", hotels[i].ID, true, false, true).Find(&hotelRooms).Error; err != nil {
			c.JSON(400, gin.H{"error": "error while fetching rooms"})
			return
		}

		// Append the rooms of the current hotel to the rooms slice
		rooms = append(rooms, hotelRooms...)
	}

	c.JSON(200, gin.H{"Hotels": hotels, "Rooms": rooms, "banners": banner})
}

// >>>>>>>>>>>>>> Banner Showing <<<<<<<<<<<<<<<<<<<<<<<<<<

func BannerShowing(c *gin.Context) {
	var banner []models.Banner

	if err := Init.DB.Preload("Hotels").Where("available = ? AND active = ?", true, true).Find(&banner).Error; err != nil {
		c.JSON(400, gin.H{"error": "banner retrireval error"})
		return
	}

	c.JSON(200, gin.H{"status": banner})
}
