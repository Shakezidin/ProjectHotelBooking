package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewUsers returns a list of all users.
func ViewUsers(c *gin.Context) {
	var users []models.User

	if err := Init.DB.Preload("Wallet").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ViewBlockedUsers returns a list of all blocked users.
func ViewBlockedUsers(c *gin.Context) {
	var users []models.User

	if err := Init.DB.Where("is_blocked = ?", true).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching blocked users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blocked users": users})
}

// ViewUnblockedUsers returns a list of all unblocked users.
func ViewUnblockedUsers(c *gin.Context) {
	var users []models.User

	if err := Init.DB.Where("is_blocked = ?", false).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching unblocked users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unblocked users": users})
}

// BlockAndUnblockUser toggles the 'IsBlocked' field of a user.
func BlockAndUnblockUser(c *gin.Context) {
	userIDStr := c.Query("id")
	if userIDStr == "" {
		c.JSON(400, gin.H{"error": "user ID query parameter is missing"})
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var user models.User

	if err := Init.DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	user.IsBlocked = !user.IsBlocked

	if err := Init.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save user availability",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "user block updation success"})
}
