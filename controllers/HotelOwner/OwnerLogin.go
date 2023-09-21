package HotelOwner

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	auth "github.com/shaikhzidhin/Auth"
	"github.com/shaikhzidhin/controllers"
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

// >>>>>>>>>>>>>> owner Signup <<<<<<<<<<<<<<<<<<<<<<<<<<

func OwnerSignUp(c *gin.Context) {
	var owner models.Owner

	time.Sleep(time.Second*10)

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
			"Message": "Email already exists",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", owner.Email).First(&owner)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone already exist",
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

	Otp := controllers.GetOTP(owner.UserName, owner.Email)

	jsonData, err := json.Marshal(owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error encoding JSON"})
		return
	}

	//inserting the otp into reddis
	err = initiializer.ReddisClient.Set(context.Background(), "signUpOTP"+owner.Email, Otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting otp in redis client"})
		return
	}

	//inserting the data into reddis
	err = initiializer.ReddisClient.Set(context.Background(), "userData"+owner.Email, jsonData, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting user data in redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "messsage": "Go to owner/signup-verification"})
}

// >>>>>>>>>>>>>> Owner OTP verification <<<<<<<<<<<<<<<<<<<<<<<<<<

func OwnerSingupVerification(c *gin.Context) {
	type otpCredentials struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err})
		return
	}

	if controllers.VerifyOTP("signUpOTP"+otpCred.Email, otpCred.Otp, c) {
		var ownerData models.Owner
		superKey := "userData" + otpCred.Email
		jsonData, err := initiializer.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting owner data from redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &ownerData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error binding reddis json data to owner variable"})
			return
		} else {
			result := initiializer.DB.Create(&ownerData)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
				return
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Otp Verification success. owner creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
	}
}

// >>>>>>>>>>>>>> Owner Login <<<<<<<<<<<<<<<<<<<<<<<<<<

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
