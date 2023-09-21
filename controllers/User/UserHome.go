package user

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> User HomePage <<<<<<<<<<<<<<<<<<<<<<<<<<

func UserHome(c *gin.Context) {
	city := c.DefaultQuery("location", "")
	if city == "" {
		c.JSON(400, gin.H{"error": "location query parameter is missing"})
		return
	}
	var hotels []models.Hotels

	if err := Init.DB.Preload("HotelCategory").Where("city = ?", city).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error"})
		return
	}
	var rooms []models.Rooms

	for i := range hotels {

		// Retrieve rooms for the current hotel
		if err := Init.DB.Where("hotels_id = ?", hotels[i].ID).Find(&rooms).Error; err != nil {
			// Handle the error, e.g., return an error response
		}

		// Assign the rooms to the current hotel
		// hotels[i].Rooms = rooms
	}

	c.JSON(200, gin.H{"Hotels": hotels, "rooms": rooms})
}

// >>>>>>>>>>>>>> User Searched Result <<<<<<<<<<<<<<<<<<<<<<<<<<

func Searching(c *gin.Context) {
	city := c.GetString("location")

	fromdatestr := c.GetString("from date")
	if fromdatestr == "" {
		c.JSON(400, gin.H{"error": "from date query parameter is missing"})
		return
	}
	fromDate, err := time.Parse("2006-01-02", fromdatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	todatestr := c.GetString("todate")
	if todatestr == "" {
		c.JSON(400, gin.H{"error": "To date query parameter is missing"})
		return
	}
	toDate, err := time.Parse("2006-01-02", todatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	childrenStr := c.DefaultQuery("number_of_children", "")
	if childrenStr == "" {
		c.JSON(400, gin.H{"error": "number of children query parameter is missing"})
		return
	}
	childrenNo, err := strconv.Atoi(childrenStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	adultStr := c.DefaultQuery("number_of_adults", "")
	if adultStr == "" {
		c.JSON(400, gin.H{"error": "number of adults query parameter is missing"})
		return
	}
	adultNo, err := strconv.Atoi(adultStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	var hotels []models.Hotels
	if err := Init.DB.Where("city = ?", city).Find(&hotels).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching hotels"})
		return
	}

	var rooms []models.Rooms
	if err := Init.DB.Where("children >= ? AND adults >= ?", childrenNo, adultNo).
		Where("isavailable = ? AND isblocked = ? AND adminapproval = ?", true, false, true).
		Where("checkin <= ? AND checkout >= ?", toDate, fromDate).
		Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching rooms"})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels, "rooms": rooms})
}
