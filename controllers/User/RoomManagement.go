package user

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Searching helps to find available rooms and hotels in a loacation
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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting 'fromdate' in Redis client"})
		return
	}
	err = initializer.ReddisClient.Set(context.Background(), "todate", req.ToDate, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting 'todate' in Redis client"})
		return
	}

	// Fetch hotels that match the location or place
	var hotels []models.Hotels
	if err := initializer.DB.Where("city = ?", req.LocOrPlace).Find(&hotels).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching hotels"})
		return
	}

	// Create a list to store hotel details including room category details
	hotelDetails := make([]map[string]interface{}, 0)

	for _, hotel := range hotels {
		// Fetch available rooms for the hotel
		var availableRooms []models.Rooms
		err := initializer.DB.Where("hotels_id = ? AND (available_rooms.room_id IS NULL OR ? > available_rooms.check_out OR ? < available_rooms.check_in)",
			hotel.ID, fromDate, toDate).Where("adults >= ? AND children >= ? AND is_blocked = ? AND admin_approval = ?", req.NumberOfAdults, req.NumberOfChildren, false, true).
			Joins("LEFT JOIN available_rooms ON rooms.id = available_rooms.room_id").
			Joins("LEFT JOIN room_categories ON rooms.room_category_id = room_categories.id").
			Preload("RoomCategory").
			Find(&availableRooms).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching rooms for the hotel"})
			return
		}

		// Create a list to store room details for each category
		categoryDetails := make(map[string][]map[string]interface{})
		addedCategories := make(map[string]bool)

		for _, room := range availableRooms {
			category := room.RoomCategory.Name
			if categoryDetails[category] == nil {
				categoryDetails[category] = make([]map[string]interface{}, 0)
			}

			// Add room details to the category only if it hasn't been added before
			if !addedCategories[category] {
				roomDetails := map[string]interface{}{
					"room_id":     room.ID,
					"description": room.Description,
					"price":       room.Price,
					"adults":      room.Adults,
					"children":    room.Children,
					"bed":         room.Bed,
					"images":      room.Images,
				}
				categoryDetails[category] = append(categoryDetails[category], roomDetails)
				addedCategories[category] = true
			}
		}

		// Calculate the available room count for each category
		categoryCounts := make(map[string]int)
		for _, room := range availableRooms {
			categoryCounts[room.RoomCategory.Name]++
		}

		// Add hotel details including room category details and room counts to the list
		hotelDetails = append(hotelDetails, map[string]interface{}{
			"hotel_name":           hotel.Name,
			"place":                hotel.City,
			"facilities":           hotel.Facility,
			"category_details":     categoryDetails,
			"available_room_count": categoryCounts, // Add available room counts
		})
	}

	c.JSON(http.StatusOK, gin.H{"hotels": hotels, "available rooms and counts": hotelDetails})
}

// RoomsView returns a list of rooms for viewing.
func RoomsView(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
		return
	}
	limit := 10
	skip := (pageInt - 1) * limit

	var rooms []models.Rooms
	var categories []models.RoomCategory

	if err := initializer.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").
		Offset(skip).Limit(limit).
		Where("is_available = ? AND is_blocked = ? AND admin_approval = ?", true, false, true).
		Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := initializer.DB.Find(&categories).Error; err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	var room models.Rooms

	if err := initializer.DB.Preload("Hotels").Preload("Cancellation").Preload("RoomCategory").
		First(&room, uint(roomID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}
