package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	helper "github.com/shaikhzidhin/helper"
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

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

func RoomFilter(c *gin.Context) {
	var filter Filter
	if err := c.BindJSON(&filter); err != nil {
		c.JSON(400, gin.H{
			"msg": err,
		})
		c.Abort()
		return
	}

	layout := "2006-01-02"

	fromdate, err := time.Parse(layout, filter.FromDate)
	todate, err := time.Parse(layout, filter.ToDate)

	fromdateStr := fromdate.Format(layout)
	todateStr := todate.Format(layout)

	err = initiializer.ReddisClient.Set(context.Background(), "fromdate", fromdateStr, 1*time.Hour).Err()
	err = initiializer.ReddisClient.Set(context.Background(), "todate", todateStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
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

	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", filter.City).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error fetching hotels"})
		return
	}

	var roomIDs []uint
	for _, hotel := range hotels {
		var tempRoom []models.Rooms
		fmt.Println(filter.MinPrice, filter.MaxPrice)
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
	roomids, err := helper.FindAvailableRoomIDs(fromdate, todate, roomIDs)
	if err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching available rooms"})
		return
	}

	var availableRooms []models.Rooms
	if err := Init.DB.Preload("Cancellation").Preload("Hotels").Preload("RoomCategory").Where("id IN (?)", roomids).Order("price " + filter.Orderby).Find(&availableRooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching available rooms"})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels, "rooms": availableRooms})
}
