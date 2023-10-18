package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
	"gorm.io/gorm"
)

// Searching finds available rooms within a specific date range and for a given location.
func Searching(c *gin.Context) {
	var req models.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	fromDate, err := time.Parse(layout, req.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from_date format"})
		return
	}
	toDate, err := time.Parse(layout, req.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to_date format"})
		return
	}

	err = initializer.ReddisClient.Set(context.Background(), "fromdate", req.FromDate, 1*time.Hour).Err()
	err = initializer.ReddisClient.Set(context.Background(), "todate", req.ToDate, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	// Fetch hotels that match the location or place
	var hotels []models.Hotels
	if err := initializer.DB.Where("city = ?", req.LocOrPlace).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching hotels"})
		return
	}

	// Store the available rooms
	availableRooms := []models.Rooms{}

	for _, hotel := range hotels {
		var filteredRooms []models.Rooms
		err := initializer.DB.
			Scopes(RoomFilters(&req)).
			Joins("LEFT JOIN available_rooms ON rooms.id = available_rooms.room_id").
			Where("rooms.hotels_id = ? AND rooms.is_available = ? "+
				"AND (available_rooms.room_id IS NULL OR "+
				"(? < available_rooms.check_out AND ? < available_rooms.check_in) AND "+
				"? > available_rooms.check_in AND ? > available_rooms.check_out)",
				hotel.ID, true, toDate, fromDate, fromDate, toDate).
			Find(&filteredRooms).Error

		if err != nil {
			fmt.Printf("Error fetching rooms: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching rooms"})
			return
		}

		// Filter rooms available within a specific date range
		for _, room := range filteredRooms {
			availableRoomIDs, err := helper.FindAvailableRoomIDs(fromDate, toDate, []uint{room.ID})
			if err != nil {
				fmt.Printf("Error finding available room IDs: %v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Error finding available room IDs"})
				return
			}

			if len(availableRoomIDs) > 0 {
				availableRooms = append(availableRooms, room)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"hotels": hotels, "available_rooms": availableRooms})
}

// RoomFilters defines filtering criteria for rooms.
func RoomFilters(req *models.SearchRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("adults >= ? AND children >= ? AND is_blocked = ? AND admin_approval = ?", req.NumberOfAdults, req.NumberOfChildren, false, true)
	}
}

// RoomsView returns a list of rooms for viewing.
func RoomsView(c *gin.Context) {
	page := c.Query("page")
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := Init.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms":      rooms,
		"categories": categories,
	})
}

// RoomDetails returns details of a specific room.
func RoomDetails(c *gin.Context) {
	roomIDStr := c.Query("id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "convert error"})
		return
	}
	var room models.Rooms

	if err := Init.DB.Preload("Hotels").Preload("Cancellation").Preload("RoomCategory").First(&room, uint(roomID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room": room,
	})
}
