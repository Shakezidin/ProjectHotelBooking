package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

// Login handles the login of admin users.
func Login(c *gin.Context) {
	var adminLogin struct {
		Username string `json:"username" validation:"required"`
		Password string `json:"password" validation:"required"`
	}
	var admin models.Admin

	// Bind the JSON request body to the adminLogin struct
	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the adminLogin struct
	if validationErr := validate.Struct(adminLogin); validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error"})
		return
	}

	// Find the admin user by username
	result := Init.DB.Where("user_name = ?", adminLogin.Username).First(&admin)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Verify the password
	if err := admin.CheckPassword(adminLogin.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate a JWT token
	tokenString, err := auth.GenerateJWT(admin.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		c.Abort()
		return
	}

	// Respond with success and the token
	c.JSON(http.StatusOK, gin.H{"status": "Success", "username": admin.UserName, "token": tokenString})
}
