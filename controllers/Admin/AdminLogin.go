package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

func AdminLogin(c *gin.Context) {
	var adminLogin struct {
		Username string `json:"username"  validation:"required"`
		Password string `json:"password" validation:"required"`
	}
	var admin models.Admin

	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		c.JSON(400, gin.H{"Error": "Admin Binding error"})
		return
	}

	validationErr := validate.Struct(admin)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	result := Init.DB.Where("username = ?", adminLogin.Username).First(&admin)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "username not found"})
		return
	}

	//password verification
	if err := admin.CheckPassword(adminLogin.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateJWT(admin.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"loginStatus": "Success", "username": admin.UserName, "token": tokenString})
}
