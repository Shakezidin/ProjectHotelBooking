package user

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	helper "github.com/shaikhzidhin/helper"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// Filter represents filtering criteria for room search.
type Filter struct {
	City     string  `json:"city" validate:"required"`
	MinPrice float64 `json:"minprice" validate:"min=5"`
	MaxPrice float64 `json:"maxprice" validate:"max=10000"`
	Orderby  string  `json:"orderby"`
	FromDate string  `json:"fromdate" validate:"required"`
	ToDate   string  `json:"todate" validate:"required"`
	Children uint    `json:"children" validate:"required"`
	Adults   uint    `json:"adults" validate:"required"`
}

// RoomFilter filters rooms based on the given criteria.
func RoomFilter(c *gin.Context) {
	var filter Filter
	if err := c.BindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		c.Abort()
		return
	}

	layout := "2006-01-02"

	fromDate, err := time.Parse(layout, filter.FromDate)
	toDate, err := time.Parse(layout, filter.ToDate)

	fromDateStr := fromDate.Format(layout)
	toDateStr := toDate.Format(layout)

	err = initializer.ReddisClient.Set(context.Background(), "fromdate", fromDateStr, 1*time.Hour).Err()
	err = initializer.ReddisClient.Set(context.Background(), "todate", toDateStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	validate := validator.New()
	validationErr := validate.Struct(filter)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"error1": "Validation error",
			"error":  validationErr.Error(),
		})
		c.Abort()
		return
	}

	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", filter.City).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error fetching hotels"})
		return
	}

	var roomIDs []uint
	for _, hotel := range hotels {
		var tempRoom []models.Rooms
		if err := Init.DB.Where("children >= ? AND adults >= ? AND hotels_id = ?", filter.Children, filter.Adults, hotel.ID).
			Where("is_available = ? AND is_blocked = ? AND admin_approval = ?", true, false, true).
			Where("price >= ? AND price <= ?", filter.MinPrice, filter.MaxPrice).Find(&tempRoom).Error; err != nil {
			c.JSON(400, gin.H{"error": "Error while fetching rooms"})
			return
		}
		for _, room := range tempRoom {
			roomIDs = append(roomIDs, room.ID)
		}
	}
	roomIDs, errr := helper.FindAvailableRoomIDs(fromDate, toDate, roomIDs)
	if errr != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching available rooms"})
		return
	}

	var availableRooms []models.Rooms
	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("id IN (?)", roomIDs).Order("price " + filter.Orderby).Find(&availableRooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hotels": hotels, "rooms": availableRooms})
}
