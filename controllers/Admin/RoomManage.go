package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func BlockedRooms(c *gin.Context) {
	var room models.Rooms

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("is_blocked", true).Find(&room).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching blocked rooms error"})
		return
	}

	c.JSON(200, gin.H{"blocked hotels": room})
}

func OwnerRooms(c *gin.Context) {
	username := c.DefaultQuery("owner_username", "")
	if username == "" {
		c.JSON(400, gin.H{"error": "owner username query parameter is missing"})
		return
	}
	var rooms []models.Rooms

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("owner_username = ?", username).Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"hotels": rooms})
}

func BlockandUnblockRooms(c *gin.Context) {
	RoomIDStr := c.DefaultQuery("roomId", "")
	if RoomIDStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(RoomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Rooms

	if err := Init.DB.First(&room, uint(roomID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "room not found",
		})
		return
	}

	room.IsBlocked = !room.IsBlocked

	if err := Init.DB.Save(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save room availability",
		})
		return
	}
	c.JSON(200, gin.H{"status": "Room block updated"})
}

func RoomsforApproval(c *gin.Context) {
	var rooms models.Rooms

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("admin_approval = ?", false).Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}

	c.JSON(200, gin.H{"approval pending rooms": rooms})
}

func RoomsApproval(c *gin.Context) {
	roomIdStr := c.DefaultQuery("roomid", "")
	if roomIdStr == "" {
		c.JSON(400, gin.H{"error": "roomId query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Rooms

	if err := Init.DB.First(&room, uint(roomID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "room not found",
		})
		return
	}

	room.AdminApproval = !room.AdminApproval

	if err := Init.DB.Save(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save room availability",
		})
		return
	}
	c.JSON(200, gin.H{"Status": "Room approval updated"})
}
