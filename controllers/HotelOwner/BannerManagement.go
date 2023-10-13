package HotelOwner

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewBanners(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	var owner models.Owner

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetcing owner"})
		return
	}

	var banners []models.Banner

	db := Init.DB

	// Retrieve banners owned by the owner
	if err := db.Where("owner_id = ?", owner.ID).Preload("Hotels").Find(&banners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching bannners"})
		return
	}

	// Render the viewBanners template
	c.JSON(http.StatusOK, gin.H{
		"banner": banners,
	})
}

func AddBanner(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	var owner models.Owner

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetcing owner"})
		return
	}

	var banner models.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(400, gin.H{"Error": "Banner binding error"})
		return
	}

	var existbannnercount int64
	if err := Init.DB.Where("hotels_id = ?", banner.HotelsId).Find(&models.Banner{}).Count(&existbannnercount).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error whiler fetching count"})
		return
	}

	if existbannnercount >= 3 {
		c.JSON(400, gin.H{"Error": "Exceeded the limit"})
		return
	}
	banner.OwnerID = owner.ID
	banner.LinkTo = "/hotel/home?id=" + strconv.Itoa(int(banner.HotelsId))

	if err := Init.DB.Create(&banner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while saving the banner"})
		return
	}

	c.JSON(200, gin.H{"Status": "Banner Added"})
}

func UpdateBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	var requestBody struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		ImageURL string `json:"Imageurl"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var banner models.Banner

	if err := Init.DB.Where("id = ?", uint(bannerID)).First(&banner).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching banner"})
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
		c.JSON(400, gin.H{"Error": "error while saving banner"})
		return
	}

	c.JSON(200, gin.H{"Status": "banner updated success"})
}

func AvailableBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	// Get the banner by ID
	var banner models.Banner
	if err := Init.DB.First(&banner, bannerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	// Toggle the availability of the banner
	banner.Available = !banner.Available

	// Save the updated banner
	if err := Init.DB.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating banner availability"})
		return
	}

	// Respond with the updated availability status
	c.JSON(http.StatusOK, gin.H{"available": "available updated"})
}

func DeleteBanner(c *gin.Context) {
	bannerID, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	// Delete the banner by ID
	if err := Init.DB.Delete(&models.Banner{}, bannerID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting banner"})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"Status": "Success",
	})
}
