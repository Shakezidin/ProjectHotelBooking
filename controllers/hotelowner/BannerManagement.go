package hotelowner

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewBanners returns all banners owned by the current owner.
func ViewBanners(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get username"})
		return
	}

	var owner models.Owner

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching owner"})
		return
	}

	var banners []models.Banner

	db := Init.DB

	// Retrieve banners owned by the owner
	if err := db.Where("owner_id = ?", owner.ID).Preload("Hotels").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching banners"})
		return
	}

	// Return the banners
	c.JSON(http.StatusOK, gin.H{
		"banners": banners,
	})
}

// AddBanner adds a new banner for the current owner.
func AddBanner(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get username"})
		return
	}

	var owner models.Owner

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching owner"})
		return
	}

	var banner models.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "banner binding error"})
		return
	}

	var existingBannerCount int64
	if err := Init.DB.Where("hotels_id = ?", banner.HotelsID).Find(&models.Banner{}).Count(&existingBannerCount).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching count"})
		return
	}

	if existingBannerCount >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "exceeded the limit"})
		return
	}
	banner.OwnerID = owner.ID
	banner.LinkTo = "/user/home/banner/hotel?id=" + strconv.Itoa(int(banner.HotelsID))

	if err := Init.DB.Create(&banner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while saving the banner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Banner added"})
}

// UpdateBanner updates an existing banner.
func UpdateBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid banner ID"})
		return
	}

	var requestBody struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		ImageURL string `json:"imageURL"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var banner models.Banner

	if err := Init.DB.Where("id = ?", uint(bannerID)).First(&banner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching banner"})
		return
	}

	if requestBody.ImageURL == "" {
		requestBody.ImageURL = banner.ImageURL
	}

	if requestBody.Subtitle == "" {
		requestBody.Subtitle = banner.Subtitle
	}

	if requestBody.Title == "" {
		requestBody.Title = banner.Title
	}

	banner.Title = requestBody.Title
	banner.Subtitle = requestBody.Subtitle
	banner.ImageURL = requestBody.ImageURL

	result := Init.DB.Save(&banner)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while saving banner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "Banner updated successfully"})
}

// AvailableBanner toggles the availability of a banner.
func AvailableBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid banner ID"})
		return
	}

	// Get the banner by ID
	var banner models.Banner
	if err := Init.DB.First(&banner, bannerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "banner not found"})
		return
	}

	// Toggle the availability of the banner
	banner.Available = !banner.Available

	// Save the updated banner
	if err := Init.DB.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while updating banner availability"})
		return
	}

	// Respond with the updated availability status
	c.JSON(http.StatusOK, gin.H{"available": "availability updated"})
}

// DeleteBanner deletes a banner.
func DeleteBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid banner ID"})
		return
	}

	// Delete the banner by ID
	if err := Init.DB.Delete(&models.Banner{}, bannerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while deleting banner"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"Status": "Success",
	})
}
