package user

import (
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func RoomFilter(c *gin.Context) {
	var filter struct {
		City     string    `json:"city" validate:"required"`
		MinPrice float64   `json:"minprice" validate:"default:0"`
		MaxPrice float64   `json:"maxprice" validate:"default:100000"`
		Orderby  string    `json:"orderby"`
		FromDate time.Time `json:"fromdate" validate:"required"`
		ToDate   time.Time `json:"todate" validate:"required"`
		Children uint      `json:"children" validate:"required"`
		Adults   uint      `json:"adults" validate:"required"`
	}

	if err := c.BindJSON(&filter); err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid JSON data",
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(filter)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"error1": "Validation error",
			"error":  validationErr.Error(),
		})
		c.Abort()
		return
	}

	if filter.Orderby == "" {
		filter.Orderby = "ASC"
	}

	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", filter.City).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error fetching hotels"})
		return
	}

	var roomIDs []uint
	for _, hotel := range hotels {
		var tempRoom []models.Rooms
		if err := Init.DB.Preload("Cancellation", "Hotels", "RoomCategory").Where("children >= ? AND adults >= ? AND hotel_id = ?", filter.Children, filter.Adults, hotel.ID).
			Where("isavailable = ? AND isblocked = ? AND adminapproval = ?", true, false, true).
			Where("price >= ? AND price <= ?", filter.MinPrice, filter.MaxPrice).Order(filter.Orderby).Find(&tempRoom).Error; err != nil {
			c.JSON(400, gin.H{"error": "Error while fetching rooms"})
			return
		}
		for _, room := range tempRoom {
			roomIDs = append(roomIDs, room.ID)
		}
	}
	roomids, err := helper.FindAvailableRoomIDs(filter.FromDate, filter.ToDate, roomIDs)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching available rooms"})
		return
	}

	var availableRooms []models.Rooms
	if err := Init.DB.Where("id IN (?)", roomids).Find(&availableRooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels, "rooms": availableRooms})
}
