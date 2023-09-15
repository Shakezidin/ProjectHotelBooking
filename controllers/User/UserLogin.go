package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

func UserSignup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Message": "binding error",
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	result := Init.DB.Where("user_name = ?", user.UserName).First(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "User_Name already Exist",
		})
		return
	}

	email := Init.DB.Where("email = ?", user.Email).First(&user)
	if email.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Email already exist",
		})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}

	record := Init.DB.Create(&user)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating owner",
			"error":   record.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Owner create Success",
	})
}

func UserLogin(c *gin.Context) {
	var userLogin models.Login
	var user models.User
	if err := c.BindJSON(&userLogin); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//checking weather the username exist or not
	result := Init.DB.Where("user_name = ?", userLogin.Username).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"msg": "user not found",
		})
		return
	}
	//checks weather user is blocked or not
	if user.Is_Block {
		c.JSON(404, gin.H{
			"msg": "user has been blocked",
		})
		c.Abort()
		return
	}

	//password verification
	if err := user.CheckPassword(userLogin.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateJWT(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"loginStatus": "Success", "username": userLogin.Username, "token": tokenString})
}
