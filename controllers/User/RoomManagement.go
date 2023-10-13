package user

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> User Searched Result <<<<<<<<<<<<<<<<<<<<<<<<<<

func Searching(c *gin.Context) {
	var req models.SearchRequest

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	// Convert the FromDate and ToDate to time.Time with default values
	fromdate, err := time.Parse(layout, req.FromDate)
	todate, err := time.Parse(layout, req.ToDate)

	fromdateStr := fromdate.Format(layout)
	todateStr := todate.Format(layout)

	err = initiializer.ReddisClient.Set(context.Background(), "fromdate", fromdateStr, 1*time.Hour).Err()
	err = initiializer.ReddisClient.Set(context.Background(), "todate", todateStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}
	// Call the GetRoomCountsByCategory function to get room counts by category
	roomCounts, err := helper.GetRoomCountsByCategory()
	if err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching room counts by category"})
		return
	}

	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", req.LocOrPlace).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching hotels"})
		return
	}

	var roomIDs []uint
	for _, hotel := range hotels {
		var tempRoom []models.Rooms
		if err := Init.DB.Where("hotels_id = ? AND adults >= ? AND children >= ? AND is_blocked = ? AND admin_approval = ?", hotel.ID, req.NumberOfAdults, req.NumberOfChildren, false, true).Find(&tempRoom).Error; err != nil {
			c.JSON(400, gin.H{"error": "Error while fetching rooms"})
			return
		}
		for _, room := range tempRoom {
			roomIDs = append(roomIDs, room.ID)
		}
	}

	// Use the modified FindAvailableRoomIDs function to get available room IDs
	roomids, err := helper.FindAvailableRoomIDs(fromdate, todate, roomIDs)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
		return
	}
	var availableRooms []models.Rooms
	if err := Init.DB.Where("id IN (?)", roomids).Find(&availableRooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels, "rooms": availableRooms, "room_counts": roomCounts})
}

// >>>>>>>>>>>>>> Display Every rooms <<<<<<<<<<<<<<<<<<<<<<<<<<

func RoomsView(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := 10
	pageInt := 1
	if p, err := strconv.Atoi(page); err == nil {
		pageInt = p
	}

	// Calculate the skip value
	skip := (pageInt - 1) * limit

	var rooms []models.Rooms
	var categories []models.RoomCategory

	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Offset(skip).Limit(limit).Where("is_available = ? AND is_blocked = ? AND admin_approval = ?", true, false, true).Find(&rooms).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"rooms":      rooms,
		"categories": categories,
	})
}

// >>>>>>>>>>>>>> See a specific room details <<<<<<<<<<<<<<<<<<<<<<<<<<

func RoomDetails(c *gin.Context) {
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

	if err := Init.DB.Preload("Hotels").Preload("Cancellation").Preload("RoomCategory").First(&room, uint(roomID)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"room": room,
	})
}
