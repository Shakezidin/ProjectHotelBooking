package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	controllers "github.com/shaikhzidhin/controllers/Otp"
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

// >>>>>>>>>>>>>> User Signup <<<<<<<<<<<<<<<<<<<<<<<<<<

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

	phone := Init.DB.Where("phone = ?", user.Phone).First(&user)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone number already exist",
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

	Otp := controllers.GetOTP(user.Name, user.Email)

	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error encoding JSON"})
		return
	}

	//inserting the otp into reddis
	err = initiializer.ReddisClient.Set(context.Background(), "signUpOTP"+user.Email, Otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting otp in redis client"})
		return
	}

	//inserting the data into reddis
	err = initiializer.ReddisClient.Set(context.Background(), "userData"+user.Email, jsonData, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting user data in redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "messsage": "Go to user/signup-verification"})
}

// >>>>>>>>>>>>>> User OTP verification <<<<<<<<<<<<<<<<<<<<<<<<<<

func SingupVerification(c *gin.Context) {
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
		var userData models.User
		superKey := "userData" + otpCred.Email
		jsonData, err := initiializer.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting user data from redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &userData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error binding reddis json data to user variable"})
			return
		} else {
			result := initiializer.DB.Create(&userData)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "falsee", "Error": result.Error})
				return
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Otp Verification success. User creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
	}
}
