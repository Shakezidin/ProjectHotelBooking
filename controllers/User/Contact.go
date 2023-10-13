package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

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
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching user"})
		return
	}

	contact := models.Contact{
		Message: message.Message,
		User_Id: user.User_Id,
	}

	if err := Init.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
