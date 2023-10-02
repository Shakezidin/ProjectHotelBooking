package HotelOwner

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func GetOwnerDashboard(c *gin.Context) {
	var hotelsCount, roomsCount int64
	var bookings []models.Booking

	db := Init.DB

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	var owner models.Owner

	if err := Init.DB.Where("username = ?", username).First(&owner); err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetcing owner"})
		return
	}

	// Retrieve the count of hotels and rooms owned by the owner
	if err := db.Model(&models.Hotels{}).Where("owner_username = ?", owner.UserName).Count(&hotelsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetchin hotels count"})
		return
	}

	if err := db.Model(&models.Rooms{}).Where("owner_username = ?", owner.UserName).Count(&roomsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetchin rooms count"})
		return
	}

	// Retrieve bookings for the owner
	if err := db.Preload("Hotel").Preload("User").Preload("Room").Where("owner_id = ?", owner.ID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetching bookings"})
		return
	}

	// Render the owner dashboard template
	c.HTML(http.StatusOK, "ownerDashboard.tmpl", gin.H{
		"revenue": owner.Revenue,
		"hotels":  hotelsCount,
		"rooms":   roomsCount,
		"booking": bookings,
	})
}

// >>>>>>>>>>>>>> owner Profile <<<<<<<<<<<<<<<<<<<<<<<<<<

func OwnerProfile(c *gin.Context) {
	var owner models.Owner

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"error": "owner not found"})
		return
	}

	c.JSON(200, gin.H{"success": owner})
}

// >>>>>>>>>>>>>> owner Profile Edit <<<<<<<<<<<<<<<<<<<<<<<<<<

func ProfileEdit(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	var owner models.Owner
	if err := Init.DB.Where("user_name = ?", username).First(&owner).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}
	var updatedowner struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.BindJSON(&updatedowner); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}
	result := Init.DB.Where("email = ?", updatedowner.Email).First(&owner)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "email already Exist",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", updatedowner.Phone).First(&owner)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone nuber already exist already Exist",
		})
		return
	}

	owner.Name = updatedowner.Name
	owner.Email = updatedowner.Email
	owner.Phone = updatedowner.Phone

	save := Init.DB.Save(&owner)
	if save.Error != nil {
		c.JSON(400, gin.H{"error": save.Error})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
