package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewUser(c *gin.Context) {
	var users []models.User

	// Retrieve all users
	if err := Init.DB.Preload("Wallet").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetcing user Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func ViewBlockedUser(c *gin.Context) {
	var users []models.User

	if err := Init.DB.Where("is_block = ?", true).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Blocked user fetching error"})
		return
	}

	c.JSON(200, gin.H{
		"blocked users": users,
	})

}

func ViewUnblockedUsers(c *gin.Context) {
	var users []models.User

	if err := Init.DB.Where("is_block = ?", false).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "un-Blocked user fetching error"})
		return
	}

	c.JSON(200, gin.H{
		"un-blocked users": users,
	})
}

func BlockandUnblockUser(c *gin.Context) {
	userIDStr := c.DefaultQuery("id", "")
	if userIDStr == "" {
		c.JSON(400, gin.H{"error": "userid query parameter is missing"})
		return
	}
	userId, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var user models.User
	if err := Init.DB.First(&user, uint(userId)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	user.Is_Block = !user.Is_Block

	if err := Init.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save user availability",
		})
		return
	}
	c.Status(200)
}
