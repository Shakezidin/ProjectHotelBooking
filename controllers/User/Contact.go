package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// SubmitContact handles submitting a contact message to the admin.
func SubmitContact(c *gin.Context) {
	var message struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching user"})
		return
	}

	contact := models.Contact{
		Message: message.Message,
		UserID:  user.ID,
	}

	if err := Init.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
