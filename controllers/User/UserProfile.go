package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Profile handles user profile retrieval.
func Profile(c *gin.Context) {
	var user models.User
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	if err := Init.DB.Preload("Wallet").Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{"success": user})
}

// ProfileEdit handles editing user profile.
func ProfileEdit(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}
	var updateuser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.BindJSON(&updateuser); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}
	result := Init.DB.Where("email = ?", updateuser.Email).First(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "email already exists",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", updateuser.Phone).First(&user)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone number already exists",
		})
		return
	}

	if updateuser.Email == "" {
		updateuser.Email = user.Email
	}

	if updateuser.Phone == "" {
		updateuser.Phone = user.Phone
	}

	if updateuser.Name == "" {
		updateuser.Name = user.Name
	}

	user.Name = updateuser.Name
	user.Email = updateuser.Email
	user.Phone = updateuser.Phone

	save := Init.DB.Save(&user)
	if save.Error != nil {
		c.JSON(400, gin.H{"error": save.Error})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

// PasswordChange handles changing user password.
func PasswordChange(c *gin.Context) {
	var pswrd struct {
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}

	if err := c.BindJSON(&pswrd); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	if err := user.CheckPassword(pswrd.OldPassword); err != nil {
		c.JSON(400, gin.H{
			"msg": "password incorrect",
		})
		return
	}

	if err := user.HashPassword(pswrd.NewPassword); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}

	result := Init.DB.Save(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

// History handles user booking history retrieval.
func History(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}
	var booking []models.Booking

	if err := Init.DB.Where("user_id = ?", user.ID).Find(&booking).Error; err != nil {
		c.JSON(400, gin.H{"Error": "error while fetching booking"})
		return
	}

	c.JSON(200, gin.H{"history": booking})
}
