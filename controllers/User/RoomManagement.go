package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func rooms(c *gin.Context) {
	var rooms []models.Rooms
	var categories []models.RoomCategory

	if err := Init.DB.Find(&rooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.HTML(200, "rooms.tmpl", gin.H{
		"rooms":      rooms,
		"categories": categories,
	})
}

func roomsFilter(c *gin.Context) {
	categoryIDStr := c.DefaultQuery("categoryid", "")
	if categoryIDStr == "" {
		c.JSON(400, gin.H{"error": "category query parameter is missing"})
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var rooms []models.Rooms
	var categories []models.RoomCategory

	if err := Init.DB.Where("categoryid = ?", uint(categoryID)).Find(&rooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.HTML(200, "rooms.tmpl", gin.H{
		"rooms":      rooms,
		"categories": categories,
	})
}

func roomDetails(c *gin.Context) {
	roomIDStr := c.DefaultQuery("roomid", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Rooms

	if err := Init.DB.First(&room, uint(roomID)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	childrents := room.Children

	c.JSON(200, gin.H{
		"room":       room,
		"childrents": childrents,
	})
}
