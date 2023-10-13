package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func GetMessages(c *gin.Context) {
	var messages []models.Contact
	if err := Init.DB.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages})
}

// deleteMessage deletes a message by ID
func DeleteMessage(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message id query missing"})
		return
	}
	var message models.Contact
	if err := Init.DB.First(&message, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error while fetching message"})
		return
	}
	if err := Init.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deleted message error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Error": "message deleted"})
}
