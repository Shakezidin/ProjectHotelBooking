package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
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

	// Convert the FromDate and ToDate to time.Time with default values
	fromDate, err := parseDateWithDefault(req.FromDate, time.Now())
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	toDate, err := parseDateWithDefault(req.ToDate, time.Now().AddDate(0, 0, 1)) // Default to tomorrow
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
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
	roomids, err := helper.FindAvailableRoomIDs(fromDate, toDate, roomIDs)
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

func parseDateWithDefault(dateStr string, defaultValue time.Time) (time.Time, error) {
	if dateStr == "" {
		return defaultValue, nil
	}
	return time.Parse("2006-01-02", dateStr)
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

	if err := Init.DB.Preload("cancellation,hotels,room_catagory").Offset(skip).Limit(limit).Where("is_available = ? AND is_blocked = ? AND admin_approval = ?",true,false,true).Find(&rooms).Error; err != nil {
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

	if err := Init.DB.Preload("Hotels").First(&room, uint(roomID)).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"room": room,
	})
}

//<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

func Searchingtwo(c *gin.Context) {
	city := c.GetString("location")

	fromdatestr := c.GetString("from_date")
	fromDate, err := time.Parse("2006-01-02", fromdatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	todatestr := c.GetString("todate")
	toDate, err := time.Parse("2006-01-02", todatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	childrenStr := c.DefaultQuery("number_of_children", "")
	childrenNo, err := strconv.Atoi(childrenStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	adultStr := c.DefaultQuery("number_of_adults", "")
	adultNo, err := strconv.Atoi(adultStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	// Get a list of hotels in the specified location
	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", city).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching hotels"})
		return
	}

	// Create a map to store room category counts for each hotel
	roomCounts := make(map[uint]map[uint]int)

	// Iterate through each hotel
	for _, hotel := range hotels {
		var tempRoom []models.Rooms

		// Get rooms for the hotel that meet the criteria
		if err := Init.DB.Where("hotels_id = ? AND adults >= ? AND children >= ? AND is_available = ?", hotel.ID, adultNo, childrenNo, true).Find(&tempRoom).Error; err != nil {
			c.JSON(400, gin.H{"error": "Error while fetching rooms"})
			return
		}

		// Create a map to store room category counts for this hotel
		hotelRoomCounts := make(map[uint]int)

		// Iterate through the rooms and count them by room category
		for _, room := range tempRoom {
			roomCategoryID := room.RoomCategoryId
			_, exists := hotelRoomCounts[roomCategoryID]
			if !exists {
				hotelRoomCounts[roomCategoryID] = 0
			}

			// Check room availability for the specified date range
			if isRoomAvailable(room.ID, fromDate, toDate) {
				hotelRoomCounts[roomCategoryID]++
			}
		}

		// Store the room category counts for this hotel
		roomCounts[hotel.ID] = hotelRoomCounts
	}

	c.JSON(200, gin.H{"hotels": hotels, "room_counts": roomCounts})
}

func isRoomAvailable(roomID uint, fromDate, toDate time.Time) bool {
	var availableRoom models.AvailableRoom
	if err := Init.DB.Where("room_id = ? AND is_available = ? AND ? NOT BETWEEN ANY(check_in) AND ? NOT BETWEEN ANY(checkout)", roomID, true, fromDate, toDate).First(&availableRoom).Error; err != nil {
		return false
	}
	return true
}
