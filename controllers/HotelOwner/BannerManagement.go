package HotelOwner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewBanners(c *gin.Context) {
	ownerId := c.MustGet("userID").(uint) // Assuming you store the owner's ID in the Gin context

	var banners []models.Banner

	db := Init.DB

	// Retrieve banners owned by the owner
	if err := db.Where("owner_id = ?", ownerId).Preload("Hotel").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching bannners"})
		return
	}

	// Render the viewBanners template
	c.JSON(http.StatusOK, gin.H{
		"banner": banners,
	})
}
