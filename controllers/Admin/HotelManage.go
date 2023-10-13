package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func BlockedHotels(c *gin.Context) {
	var hotels []models.Hotels

	if err := Init.DB.Where("is_block", true).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching blocked hotels error"})
		return
	}

	c.JSON(200, gin.H{"blocked hotels": hotels})
}

func OwnerHotels(c *gin.Context) {
	username := c.DefaultQuery("owner_username", "")
	if username == "" {
		c.JSON(400, gin.H{"error": "owner username query parameter is missing"})
		return
	}
	var hotels []models.Hotels

	if err := Init.DB.Where("owner_username = ?", username).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels})
}

func BlockandUnblockhotel(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var hotel models.Hotels

	if err := Init.DB.First(&hotel, uint(hotelID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "hotel not found",
		})
		return
	}

	hotel.IsBlock = !hotel.IsBlock

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save hotel availability",
		})
		return
	}
	c.Status(200)
}

func HotelforApproval(c *gin.Context) {
	var hotel models.Hotels

	if err := Init.DB.Preload("HotelCategory").Where("admin_approval = ?", false).Find(&hotel).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"approval pending hotels": hotel})
}

func HotelsApproval(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var hotel models.Hotels

	if err := Init.DB.First(&hotel, uint(hotelID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "hotel not found",
		})
		return
	}

	hotel.AdminApproval = !hotel.AdminApproval

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save hotel availability",
		})
		return
	}
	c.Status(200)
}
