package HotelOwner

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> view fecilities <<<<<<<<<<<<<<<<<<<<<<<<<<
func ViewHotelFecilities(c *gin.Context) {
	var fecilities []models.HotelAmenities
	if err := Init.DB.Find(&fecilities).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	// var hotel models.Hotel
	c.JSON(200, gin.H{
		"fecilities": fecilities,
		// "hotel": hotel,
	})
}

// >>>>>>>>>>>>>> Add hotel <<<<<<<<<<<<<<<<<<<<<<<<<<
func AddHotel(c *gin.Context) {
	var hotel models.Hotel
	if err := c.BindJSON(&hotel); err != nil {
		c.JSON(400, gin.H{
			"msg":   "binding error1",
			"error": err,
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(hotel)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"msg":   "validate error2",
			"error": validationErr.Error(),
		})
		c.Abort()
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	hotel.OwnerUsername = username
	result := Init.DB.Create(&hotel)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

// >>>>>>>>>>>>>> view hotels <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewHotels(c *gin.Context) {
	var hotel []models.Hotel
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	if err := Init.DB.Where("owner_username = ?", username).Find(&hotel).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{
		"msg": hotel,
	})
}

// >>>>>>>>>>>>>> view specific hotel <<<<<<<<<<<<<<<<<<<<<<<<<<
func ViewSpecificHotel(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("hotelid", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var hotel models.Hotel

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	if err := Init.DB.Where("owner_username = ? AND id = ?", username, uint(hotelId)).First(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg": hotel,
	})
}

// >>>>>>>>>>>>>> edited hotel details saving <<<<<<<<<<<<<<<<<<<<<<<<<<

func Hoteledit(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("hotelid", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}

	var hotel models.Hotel
	if err := Init.DB.Where("owner_username = ? AND id = ?", username, uint(hotelId)).First(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	var updatedHotel struct {
		Name          string  `json:"name" validate:"required"`
		Title         string  `json:"title" validate:"required"`
		Description   string  `json:"description" validate:"required"`
		StartingPrice float64 `json:"startingprice" validate:"required"`
		City          string  `json:"city" validate:"required"`
		Pincode       string  `json:"pincode" validate:"required"`
		Address       string  `json:"address" validate:"required"`
		Images        string  `json:"images" validate:"required"`
		TypesOfRoom   int     `json:"typesofroom" validate:"required"`
		IsBlock       bool
	}

	if err := c.BindJSON(&updatedHotel); err != nil {
		c.JSON(400, gin.H{
			"msg":   "binding error1",
			"error": err,
		})
		c.Abort()
		return
	}

	// Update the fields of the existing hotel record
	hotel.Name = updatedHotel.Name
	hotel.Title = updatedHotel.Title
	hotel.Description = updatedHotel.Description
	hotel.StartingPrice = updatedHotel.StartingPrice
	hotel.City = updatedHotel.City
	hotel.Pincode = updatedHotel.Pincode
	hotel.Address = updatedHotel.Address
	hotel.Images = updatedHotel.Images
	hotel.TypesOfRoom = updatedHotel.TypesOfRoom
	hotel.IsBlock = true

	// Save the updated hotel record
	result := Init.DB.Save(&hotel)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

// >>>>>>>>>>>>>> delete hotel<<<<<<<<<<<<<<<<<<<<<<<<<<

func DeleteHotel(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("hotelid", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	var hotel models.Hotel
	result := Init.DB.Where("owner_username = ? AND id = ?", username, uint(hotelId)).Delete(&hotel)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "hotel not found"})
		return
	}
	c.JSON(200, gin.H{"msg": "deleted"})
}

func HotelAvailability(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("hotelid", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	var hotel models.Hotel
	if err := Init.DB.First(&hotel, hotelId).Error; err != nil {
		c.JSON(404, gin.H{
			"hello": "Hot",
			"error": "Room not found",
		})
		return
	}

	hotel.IsAvailable = !hotel.IsAvailable

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save Room availability",
		})
		return
	}

	c.Status(200)
}
