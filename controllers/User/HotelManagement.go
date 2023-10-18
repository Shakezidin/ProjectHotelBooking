package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// SearchHotelByName searches for hotels by name based on the provided request.
func SearchHotelByName(c *gin.Context) {
	var req models.SearchRequest

	// Bind the request data to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	// Convert the FromDate and ToDate to time.Time with default values
	fromdate, err := time.Parse(layout, req.FromDate)
	todate, err := time.Parse(layout, req.ToDate)

	fromdateStr := fromdate.Format(layout)
	todateStr := todate.Format(layout)

	err = initializer.ReddisClient.Set(context.Background(), "fromdate", fromdateStr, 1*time.Hour).Err()
	err = initializer.ReddisClient.Set(context.Background(), "todate", todateStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	// Call the GetRoomCountsByCategory function to get room counts by category
	roomCounts, err := helper.GetRoomCountsByCategory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching room counts by category"})
		return
	}

	var hotels []models.Hotels
	result := Init.DB.Where("name ILIKE ?", "%"+req.LocOrPlace+"%").Find(&hotels)
	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": result.Error})
		return
	}

	var roomIDs []uint
	for _, hotel := range hotels {
		var tempRoom []models.Rooms
		if err := Init.DB.Where("hotels_id = ? AND adults >= ? AND children >= ? AND is_blocked = ? AND admin_approval = ?", hotel.ID, req.NumberOfAdults, req.NumberOfChildren, false, true).Find(&tempRoom).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching rooms"})
			return
		}
		for _, room := range tempRoom {
			roomIDs = append(roomIDs, room.ID)
		}
	}

	// Use the modified FindAvailableRoomIDs function to get available room IDs
	roomIDs, errr := helper.FindAvailableRoomIDs(fromdate, todate, roomIDs)
	if errr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	var availableRooms []models.Rooms
	if err := Init.DB.Where("id IN (?)", roomIDs).Find(&availableRooms).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotels": hotels, "rooms": availableRooms, "room_counts": roomCounts})
}
