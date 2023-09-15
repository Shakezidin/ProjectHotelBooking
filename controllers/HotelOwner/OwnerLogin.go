package HotelOwner

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

func OwnerSignUp(c *gin.Context) {
	var owner models.Owner

	if err := c.ShouldBindJSON(&owner); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Message": "binding error",
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(owner)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	result := Init.DB.Where("user_name = ?", owner.UserName).First(&owner)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "User_Name already Exist",
		})
		return
	}

	email := Init.DB.Where("email = ?", owner.Email).First(&owner)
	if email.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Email already exist",
		})
		return
	}

	if err := owner.HashPassword(owner.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}

	record := Init.DB.Create(&owner)
	if record.Error != nil {
		// Log the error for debugging purposes.
		fmt.Println("Database error:", record.Error)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating owner",
			"error":   record.Error.Error(), // Include the specific database error message.
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Owner create Success",
	})
}

func OwnerLogin(c *gin.Context) {
	var ownerLogin models.Login
	var owner models.Owner
	if err := c.BindJSON(&ownerLogin); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//checking weather the username exist or not
	result := Init.DB.Where("user_name = ?", ownerLogin.Username).First(&owner)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"msg": "user not found",
		})
		return
	}
	//checks weather user is blocked or not
	if owner.Is_Block {
		c.JSON(404, gin.H{
			"msg": "user has been blocked",
		})
		c.Abort()
		return
	}

	//password verification
	if err := owner.CheckPassword(ownerLogin.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateJWT(owner.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"loginStatus": "Success", "username": ownerLogin.Username, "token": tokenString})
}