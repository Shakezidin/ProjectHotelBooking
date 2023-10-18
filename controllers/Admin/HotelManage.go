package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// BlockedHotels returns a list of blocked hotels.
func BlockedHotels(c *gin.Context) {
	var hotels []models.Hotels

	if err := Init.DB.Where("is_block = ?", true).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error fetching blocked hotels"})
		return
	}

	c.JSON(200, gin.H{"blockedHotels": hotels})
}

// OwnerHotels returns a list of hotels owned by a specific owner.
func OwnerHotels(c *gin.Context) {
	ownerUsername := c.Query("ownerUsername")
	if ownerUsername == "" {
		c.JSON(400, gin.H{"error": "Owner username query parameter is missing"})
		return
	}

	var hotels []models.Hotels

	if err := Init.DB.Where("owner_username = ?", ownerUsername).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels})
}

// BlockAndUnblockHotel toggles the 'isBlock' field of a hotel.
func BlockAndUnblockHotel(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "Hotel ID query parameter is missing"})
		return
	}

	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Conversion error"})
		return
	}

	var hotel models.Hotels

	if err := Init.DB.First(&hotel, uint(hotelID)).Error; err != nil {
		c.JSON(404, gin.H{"error": "Hotel not found"})
		return
	}

	hotel.IsBlock = !hotel.IsBlock

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save hotel availability"})
		return
	}

	c.Status(200)
}

// HotelsForApproval returns a list of hotels pending approval.
func HotelsForApproval(c *gin.Context) {
	var hotels []models.Hotels

	if err := Init.DB.Preload("HotelCategory").Where("admin_approval = ?", false).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"hotelsForApproval": hotels})
}

// HotelsApproval toggles the 'adminApproval' field of a hotel.
func HotelsApproval(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "Hotel ID query parameter is missing"})
		return
	}

	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Conversion error"})
		return
	}

	var hotel models.Hotels

	if err := Init.DB.First(&hotel, uint(hotelID)).Error; err != nil {
		c.JSON(404, gin.H{"error": "Hotel not found"})
		return
	}

	hotel.AdminApproval = !hotel.AdminApproval

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save hotel availability"})
		return
	}

	c.Status(200)
}
