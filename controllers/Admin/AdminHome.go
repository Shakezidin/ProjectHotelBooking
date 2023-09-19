package admin

import (
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func Dashboard(c *gin.Context) {
	var (
		noOfUsers  int64
		noOfOwners int64
		noOfHotels int64
		noOfRooms  int64
	)

	if err := Init.DB.Model(&models.User{}).Count(&noOfUsers).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Model(&models.Owner{}).Count(&noOfOwners).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Model(&models.Hotels{}).Count(&noOfHotels).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Model(&models.Rooms{}).Count(&noOfRooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"users":  noOfUsers,
		"owners": noOfOwners,
		"hotels": noOfHotels,
		"rooms":  noOfRooms,
	})
}
