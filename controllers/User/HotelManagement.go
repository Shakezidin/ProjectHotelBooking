package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func SearchHotel(c *gin.Context) {
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
	result := Init.DB.Where("name ILIKE ?", "%"+req.LocOrPlace+"%").Find(&hotels)
	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": result.Error})
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

// func SearchedHotelFilter(c *gin.Context) {
// 	var filter struct {
// 		Hotelname string    `json:"hotelname" validate:"required"`
// 		MinPrice  float64   `json:"minprice" validate:"default:0"`
// 		MaxPrice  float64   `json:"maxprice" validate:"default:100000"`
// 		Orderby   string    `json:"orderby"`
// 		FromDate  time.Time `json:"fromdate" validate:"required"`
// 		ToDate    time.Time `json:"todate" validate:"required"`
// 		Children  uint      `json:"children" validate:"required"`
// 		Adults    uint      `json:"adults" validate:"required"`
// 	}
// 	if err := c.BindJSON(&filter); err != nil {
// 		c.JSON(400, gin.H{
// 			"msg":   "binding error",
// 			"error": err,
// 		})
// 		c.Abort()
// 		return
// 	}
// 	validationErr := validate.Struct(filter)
// 	if validationErr != nil {
// 		c.JSON(400, gin.H{
// 			"error": validationErr.Error(),
// 		})
// 		c.Abort()
// 		return
// 	}

// 	if filter.Orderby == "" {
// 		filter.Orderby = "ASC"
// 	}

// 	var hotels []models.Hotels
// 	if err := Init.DB.Where("city = ?", filter.Hotelname).Where("name ILIKE ?", "%"+filter.Hotelname+"%").Find(&hotels).Error; err != nil {
// 		c.JSON(400, gin.H{"error": "fetching hotels"})
// 		return
// 	}

// 	var roomIDs []uint
// 	for _, hotel := range hotels {
// 		var tempRoom []models.Rooms
// 		if err := Init.DB.Preload("Cancellation", "Hotels", "RoomCategory").Where("children >= ? AND adults >= ? AND hotel_id = ?", filter.Children, filter.Adults, hotel.ID).
// 			Where("isavailable = ? AND isblocked = ? AND adminapproval = ?", true, false, true).
// 			Where("price >= ? AND price <= ?", filter.MinPrice, filter.MaxPrice).Order(filter.Orderby).Find(&tempRoom).Error; err != nil {
// 			c.JSON(400, gin.H{"error": "Error while fetching rooms"})
// 			return
// 		}
// 		for _, room := range tempRoom {
// 			roomIDs = append(roomIDs, room.ID)
// 		}
// 	}
// 	roomids, err := helper.FindAvailableRoomIDs(filter.FromDate, filter.ToDate, roomIDs)
// 	if err != nil {
// 		c.JSON(400, gin.H{"Error": "Error while fetching available rooms"})
// 		return
// 	}

// 	var availableRooms []models.Rooms
// 	if err := Init.DB.Where("id IN (?)", roomids).Find(&availableRooms).Error; err != nil {
// 		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
// 		return
// 	}

// 	c.JSON(200, gin.H{"hotels": hotels, "rooms": availableRooms})
// }
