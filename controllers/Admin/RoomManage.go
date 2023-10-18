package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

//ViewRooms returns rooms
func ViewRooms(c *gin.Context) {
	var rooms []models.Rooms
	if err := Init.DB.Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching rooms"})
		return
	}
	c.JSON(200, gin.H{"rooms": rooms})
}

// BlockedRooms returns a list of all blocked rooms.
func BlockedRooms(c *gin.Context) {
	var rooms []models.Rooms

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("is_blocked", true).Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching blocked rooms"})
		return
	}

	c.JSON(200, gin.H{"blocked rooms": rooms})
}

// BlockAndUnblockRooms toggles the 'IsBlocked' field of a room.
func BlockAndUnblockRooms(c *gin.Context) {
	roomIDStr := c.DefaultQuery("id", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "room ID query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
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

// RoomsForApproval returns a list of rooms pending admin approval.
func RoomsForApproval(c *gin.Context) {
	var rooms []models.Rooms

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("admin_approval = ?", false).Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching rooms pending approval"})
		return
	}

	c.JSON(200, gin.H{"approval pending rooms": rooms})
}

// RoomsApproval toggles the 'AdminApproval' field of a room.
func RoomsApproval(c *gin.Context) {
	roomIDStr := c.DefaultQuery("id", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "room ID query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
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
	c.JSON(200, gin.H{"status": "Room approval updated"})
}
