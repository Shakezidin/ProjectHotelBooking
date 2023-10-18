package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	controllers "github.com/shaikhzidhin/controllers/Otp"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Login handles user login.
func Login(c *gin.Context) {
	var userLogin models.Login
	var user models.User
	if err := c.BindJSON(&userLogin); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// Checking whether the username exists or not
	result := Init.DB.Where("user_name = ?", userLogin.Username).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"msg": "user not found",
		})
		return
	}
	// Checks whether the user is blocked or not
	if user.IsBlocked {
		c.JSON(404, gin.H{
			"msg": "user has been blocked",
		})
		c.Abort()
		return
	}
	// Password verification
	if err := user.CheckPassword(userLogin.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	tokenString, err := Auth.GenerateJWT(user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	token := tokenString
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err = initializer.ReddisClient.Set(context.Background(), "userId", user.ID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{"loginStatus": "Success", "username": userLogin.Username, "token": tokenString})
}

// ForgetPassword handles the user's password recovery request.
func ForgetPassword(c *gin.Context) {
	var emailData struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := c.BindJSON(&emailData); err != nil {
		c.JSON(400, gin.H{"Error": err})
		return
	}
	var user models.User
	if err := Init.DB.Where("email = ?", emailData.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	if user.Phone != emailData.Phone {
		c.JSON(400, gin.H{"error": "email incorrect"})
		return
	}

	Otp := controllers.GetOTP(user.Name, user.Email)

	err := initializer.ReddisClient.Set(context.Background(), user.Email, Otp, 10*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{"status": "recovery OTP sent to the email address"})
}

// VerifyOTP handles OTP verification for password recovery.
func VerifyOTP(c *gin.Context) {
	type otpCredentials struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err})
		return
	}
	otp, err := initializer.ReddisClient.Get(context.Background(), otpCred.Email).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'todate' from Redis client"})
		return
	}

	if otp != otpCred.Otp {
		c.JSON(400, gin.H{"Error": "Invalid OTP"})
		return
	}

	tokenString, err := Auth.GenerateJWT(otpCred.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"status": "success, go to user/newpassword", "token": tokenString})
}

// NewPassword handles setting a new password for the user.
func NewPassword(c *gin.Context) {
	var password struct {
		NewPassword     string `json:"newpassword"`
		ConfirmPassword string `json:"confirmpassword"`
	}

	if err := c.BindJSON(&password); err != nil {
		c.JSON(400, gin.H{"error": "binding new password error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	email, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "email didn't get"})
		return
	}

	if password.NewPassword != password.ConfirmPassword {
		c.JSON(400, gin.H{"error": "Password mismatch"})
		return
	}
	var user models.User
	if err := Init.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	if err := user.HashPassword(password.NewPassword); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}
	result := Init.DB.Save(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "saving error"})
		return
	}

	c.JSON(200, gin.H{"status": "password updated"})
}
