package hotelowner

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// AddHotel adds a new hotel.
func AddHotel(c *gin.Context) {
	var hotel models.Hotels
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(400, gin.H{
			"message": "binding error",
			"error":   err,
		})
		c.Abort()
		return
	}

	fmt.Println(hotel)

	validationErr := validate.Struct(hotel)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"message": "validation error",
			"error":   validationErr.Error(),
		})
		c.Abort()
		return
	}

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	hotel.OwnerUsername = username
	result := Init.DB.Create(&hotel)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"status": "hotel added"})
}

// ViewHotels retrieves all hotels owned by the user.
func ViewHotels(c *gin.Context) {
	var hotels []models.Hotels
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}
	if err := Init.DB.Preload("HotelCategory").Where("owner_username = ?", username).Find(&hotels).Error; err != nil {
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(200, gin.H{"hotels": hotels})
}

// ViewSpecificHotel retrieves a specific hotel by its ID.
func ViewSpecificHotel(c *gin.Context) {
	hotelIDStr := c.Query("id")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotel ID query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "conversion error"})
		return
	}
	var hotel models.Hotels

	if err := Init.DB.Where("id = ?", uint(hotelID)).First(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"hotel": hotel})
}

// Hoteledit updates a hotel's details.
func Hoteledit(c *gin.Context) {
	hotelIDStr := c.Query("id")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotel ID query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "conversion error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	var hotel models.Hotels
	if err := Init.DB.Where("owner_username = ? AND id = ?", username, uint(hotelID)).First(&hotel).Error; err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
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
		CategoryID    uint    `json:"category_id" validate:"required"`
		IsBlock       bool
	}

	if err := c.BindJSON(&updatedHotel); err != nil {
		c.JSON(400, gin.H{
			"message": "binding error",
			"error":   err,
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
	hotel.HotelCategoryID = updatedHotel.CategoryID
	hotel.IsBlock = true

	// Save the updated hotel record
	result := Init.DB.Save(&hotel)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "hotel updated"})
}

// HotelAvailability updates the availability status of a hotel.
func HotelAvailability(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotel ID query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "conversion error"})
		return
	}

	var hotel models.Hotels
	if err := Init.DB.First(&hotel, hotelID).Error; err != nil {
		c.JSON(404, gin.H{"error": "hotel not found"})
		return
	}

	hotel.IsAvailable = !hotel.IsAvailable

	if err := Init.DB.Save(&hotel).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to save hotel availability"})
		return
	}

	c.JSON(200, gin.H{"status": "hotel availability updated"})
}

// DeleteHotel deletes a hotel.
func DeleteHotel(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	if hotelIDStr == "" {
		c.JSON(400, gin.H{"error": "hotel ID query parameter is missing"})
		return
	}
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "conversion error"})
		return
	}
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didn't get"})
		return
	}

	result := Init.DB.Where("owner_username = ? AND id = ?", username, uint(hotelID)).Delete(&models.Hotels{})
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "hotel not found"})
		return
	}
	c.JSON(200, gin.H{"message": "hotel deleted"})
}
