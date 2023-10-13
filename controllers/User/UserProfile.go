package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> User Profile <<<<<<<<<<<<<<<<<<<<<<<<<<

func Profile(c *gin.Context) {
	var user models.User
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	if err := Init.DB.Preload("Wallet").Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{"success": user})
}

// >>>>>>>>>>>>>> User Profile Edit <<<<<<<<<<<<<<<<<<<<<<<<<<

func ProfileEdit(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
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
			"Message": "email already Exist",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", updateuser.Phone).First(&user)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone nuber already exist already Exist",
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

// >>>>>>>>>>>>>> User Password Change <<<<<<<<<<<<<<<<<<<<<<<<<<

func PasswordChange(c *gin.Context) {
	var pswrd struct {
		Old_password string `json:"oldpassword"`
		New_password string `json:"newpassword"`
	}

	if err := c.BindJSON(&pswrd); err != nil {
		c.JSON(400, gin.H{"error": "Binding error"})
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
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

	if err := user.CheckPassword(pswrd.Old_password); err != nil {
		c.JSON(400, gin.H{
			"msg": "password incorrect",
		})
		return
	}


	if err := user.HashPassword(pswrd.New_password); err != nil {
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

//you can add history booking here

func History(c *gin.Context) {
	userId := c.GetUint("userid")
	var booking models.Booking

	if err := Init.DB.Where("user_id = ?", userId).Find(&booking).Error; err != nil {
		c.JSON(400, gin.H{"Error": "error while etching booking"})
		return
	}

	c.JSON(200, gin.H{"history": booking})
}