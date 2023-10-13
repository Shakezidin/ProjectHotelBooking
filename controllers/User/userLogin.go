package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	Auth "github.com/shaikhzidhin/Auth"
	auth "github.com/shaikhzidhin/Auth"
	controllers "github.com/shaikhzidhin/controllers/Otp"
	"github.com/shaikhzidhin/initiializer"

	// "github.com/shaikhzidhin/controllers"

	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> User Login <<<<<<<<<<<<<<<<<<<<<<<<<<

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
	token := tokenString
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	err = initiializer.ReddisClient.Set(context.Background(), "userId", user.User_Id, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{"loginStatus": "Success", "username": userLogin.Username, "token": tokenString})
}

func ForgetPassword(c *gin.Context) {
	var Email struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := c.BindJSON(&Email); err != nil {
		c.JSON(400, gin.H{"Error": err})
		return
	}
	var user models.User
	if err := Init.DB.Where("email = ?", Email.Email).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	if user.Phone != Email.Phone {
		c.JSON(400, gin.H{"error": "email incorrect"})
		return
	}

	Otp := controllers.GetOTP(user.Name, user.Email)

	session := sessions.Default(c)
	session.Set(user.Email, Otp)
	session.Save()

	c.JSON(200, gin.H{"status": "recovery otp sent to email address"})
}

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
	session := sessions.Default(c)
	if session.Get(otpCred.Email) != otpCred.Otp {
		c.JSON(400, gin.H{"error": "otp validation error"})
		return
	}
	tokenString, err := auth.GenerateJWT(otpCred.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"status": "success,go to user/newpassword", "token": tokenString})
}

func Newpassword(c *gin.Context) {
	var password struct {
		Newpassword     string `json:"newpassword"`
		ConfirmPassword string `json:"confirmpassword"`
	}

	if err := c.BindJSON(&password); err != nil {
		c.JSON(400, gin.H{"error": "binding new passwod error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	email, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	if password.Newpassword != password.ConfirmPassword {
		c.JSON(400, gin.H{"error": "Password mismach"})
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

	if err := user.HashPassword(password.Newpassword); err != nil {
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
