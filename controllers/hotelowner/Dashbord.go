package hotelowner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// GetOwnerDashboard provides information for the owner's dashboard.
func GetOwnerDashboard(c *gin.Context) {
	var hotelsCount, roomsCount int64
	var bookings []models.Booking
	var owner models.Owner

	db := Init.DB

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get username"})
		return
	}

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "error while fetching owner"})
		return
	}

	// Retrieve the count of hotels and rooms owned by the owner
	if err := db.Model(&models.Hotels{}).Where("owner_username = ?", owner.UserName).Count(&hotelsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while fetching hotels count"})
		return
	}

	if err := db.Model(&models.Rooms{}).Where("owner_username = ?", owner.UserName).Count(&roomsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while fetching rooms count"})
		return
	}

	// Retrieve bookings for the owner
	if err := db.Preload("Hotels").Preload("User").Preload("Rooms").Where("owner_id = ?", owner.ID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error while fetching bookings"})
		return
	}

	// Provide owner dashboard data
	c.JSON(http.StatusOK, gin.H{
		"revenue":  owner.Revenue,
		"hotels":   hotelsCount,
		"rooms":    roomsCount,
		"bookings": bookings,
	})
}

// OwnerProfile retrieves the owner's profile.
func OwnerProfile(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get username"})
		return
	}
	var owner models.Owner

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": owner})
}

// ProfileEdit updates the owner's profile.
func ProfileEdit(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get username"})
		return
	}
	var owner models.Owner
	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	var updatedOwner struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.BindJSON(&updatedOwner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "binding error"})
		return
	}
	result := Init.DB.Where("email = ?", updatedOwner.Email).First(&owner)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "email already exists",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", updatedOwner.Phone).First(&owner)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "phone number already exists",
		})
		return
	}

	if updatedOwner.Email == "" {
		updatedOwner.Email = owner.Email
	}

	if updatedOwner.Phone == "" {
		updatedOwner.Phone = owner.Phone
	}

	if updatedOwner.Name == "" {
		updatedOwner.Name = owner.Name
	}

	owner.Name = updatedOwner.Name
	owner.Email = updatedOwner.Email
	owner.Phone = updatedOwner.Phone

	save := Init.DB.Save(&owner)
	if save.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": save.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
