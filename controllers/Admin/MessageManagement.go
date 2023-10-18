package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// GetMessages returns a list of all messages.
func GetMessages(c *gin.Context) {
	var messages []models.Contact
	if err := Init.DB.Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// DeleteMessage deletes a message by its ID.
func DeleteMessage(c *gin.Context) {
	messageID := c.DefaultQuery("id", "")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID query parameter is missing"})
		return
	}

	var message models.Contact
	if err := Init.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error while fetching message"})
		return
	}

	if err := Init.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Deleted message error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
