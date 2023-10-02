package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func BannerView(c *gin.Context) {
	var banners []models.Banner
	if err := Init.DB.Preload("Owner").Preload("Hotels").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while fetching banner",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"banners": banners,
	})
}

func BannerSetActive(c *gin.Context) {
	id := c.Query("bannerid")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	var banner models.Banner
	if err := Init.DB.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Banner not found"})
		return
	}

	banner.Active = !banner.Active
	if err := Init.DB.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update banner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func BannerDetails(c *gin.Context) {
	bannerID := c.Query("bannerid")

	if bannerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
		return
	}

	var banner models.Banner
	if err := Init.DB.First(&banner, bannerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Banner not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"banner": banner,
	})
}

func DeleteBanner(c *gin.Context) {
	bannerID := c.Query("bannerid")

	if bannerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter"})
		return
	}

	var banner models.Banner
	if err := Init.DB.Delete(&banner, bannerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Banner not found"})
		return
	}

	c.JSON(200, gin.H{"status": "Deleted"})
}
