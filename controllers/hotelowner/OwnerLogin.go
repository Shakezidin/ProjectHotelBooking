package hotelowner

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	auth "github.com/shaikhzidhin/Auth"
	controllers "github.com/shaikhzidhin/controllers/Otp"

	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

// OwnerSignUp handles owner registration.
func OwnerSignUp(c *gin.Context) {
	var owner models.Owner

	time.Sleep(time.Second * 10)

	if err := c.ShouldBindJSON(&owner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "binding error",
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(owner)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error"})
		return
	}

	result := Init.DB.Where("user_name = ?", owner.UserName).First(&owner)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username already exists",
		})
		return
	}

	email := Init.DB.Where("email = ?", owner.Email).First(&owner)
	if email.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already exists",
		})
		return
	}

	phone := Init.DB.Where("phone = ?", owner.Phone).First(&owner)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone number already exists",
		})
		return
	}

	if err := owner.HashPassword(owner.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "hashing error",
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

	// Inserting the OTP into Redis
	err = initializer.ReddisClient.Set(context.Background(), "signUpOTP"+owner.Email, Otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting OTP in Redis client"})
		return
	}

	// Inserting the data into Redis
	err = initializer.ReddisClient.Set(context.Background(), "userData"+owner.Email, jsonData, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting user data in Redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Go to owner/signup-verification"})
}

// OwnerSignUpVerification handles owner registration OTP verification.
func OwnerSignUpVerification(c *gin.Context) {
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
		jsonData, err := initializer.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting owner data from Redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &ownerData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error binding Redis JSON data to owner variable"})
			return
		}
		result := initializer.DB.Create(&ownerData)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "OTP verification success. Owner creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
	}
}

// OwnerLogin handles owner login.
func OwnerLogin(c *gin.Context) {
	var ownerLogin models.Login
	var owner models.Owner
	if err := c.BindJSON(&ownerLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	// Checking whether the username exists or not
	result := Init.DB.Where("user_name = ?", ownerLogin.Username).First(&owner)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}
	// Checks whether the user is blocked or not
	if owner.IsBlocked {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User has been blocked",
		})
		c.Abort()
		return
	}

	// Password verification
	if err := owner.CheckPassword(ownerLogin.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateJWT(owner.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"loginStatus": "Success", "username": ownerLogin.Username, "token": tokenString})
}
