package hotelowner

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// AddingRoom adds room details to a hotel.
func AddingRoom(c *gin.Context) {
	numberOfRooms, _ := strconv.Atoi(c.Query("numberofrooms"))
	floorNumber, _ := strconv.Atoi(c.Query("floornumber"))
	hotelID, _ := strconv.Atoi(c.Query("id"))

	var roomToBind models.Rooms
	if err := c.ShouldBindJSON(&roomToBind); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "binding error",
			"error":   err.Error(),
		})
		c.Abort()
		return
	}

	roomNumber := 1
	for i := 1; i <= numberOfRooms; i++ {
		var room models.Rooms

		room.Adults = roomToBind.Adults
		room.Bed = roomToBind.Bed
		room.Children = roomToBind.Children
		room.Description = roomToBind.Description
		room.Discount = roomToBind.Discount
		room.Facility = roomToBind.Facility
		room.HotelsID = uint(hotelID)
		room.Images = roomToBind.Images
		room.CancellationID = roomToBind.CancellationID
		room.IsAvailable = roomToBind.IsAvailable
		room.Price = roomToBind.Price
		room.RoomNo = floorNumber*100 + roomNumber
		room.DiscountPrice = roomToBind.Price - (roomToBind.Price * roomToBind.Discount / 100)
		room.RoomCategoryID = roomToBind.RoomCategoryID

		header := c.Request.Header.Get("Authorization")
		username, err := Auth.Trim(header)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Username not found"})
			return
		}
		room.OwnerUsername = username

		if err := Init.DB.Create(&room).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		roomNumber++
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// EditRoom updates room details.
func EditRoom(c *gin.Context) {
	hotelIDStr := c.DefaultQuery("id", "")
	roomCategoryIDStr := c.DefaultQuery("room_category_id", "")

	// Parse room_category_id to uint64
	roomCategoryID, err := strconv.ParseUint(roomCategoryIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room_category_id"})
		return
	}

	// Check if hotel_id is provided
	if hotelIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hotel_id query parameter is missing"})
		return
	}

	// Parse hotel_id to int
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel_id"})
		return
	}

	// Find all rooms based on hotel_id and room_category_id
	var rooms []models.Rooms
	if err := Init.DB.Where("hotels_id = ? AND room_category_id = ?", uint(hotelID), uint(roomCategoryID)).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Bind the updated room data from JSON request
	var updatedRoom struct {
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Adults      int     `json:"adults"`
		Children    int     `json:"children"`
		Bed         string  `json:"bed"`
		Images      string  `json:"images"`
		IsAvailable bool    `json:"is_available"`
		Discount    float64 `json:"discount"`
		// Facilities  []string `json:"facilities" gorm:"type:jsonb"`
	}
	if err := c.ShouldBindJSON(&updatedRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update each room individually
	for i := range rooms {
		rooms[i].Description = updatedRoom.Description
		rooms[i].Price = updatedRoom.Price
		rooms[i].Adults = updatedRoom.Adults
		rooms[i].Children = updatedRoom.Children
		rooms[i].Bed = updatedRoom.Bed
		rooms[i].Images = updatedRoom.Images
		rooms[i].IsAvailable = updatedRoom.IsAvailable
		rooms[i].Discount = updatedRoom.Discount
		// rooms[i].Fecilities = updatedRoom.Facilities
		rooms[i].DiscountPrice = rooms[i].Price - (rooms[i].Price * rooms[i].Discount / 100)

		// Save the updated room
		result := Init.DB.Save(&rooms[i])
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// ViewRooms returns a list of rooms for a hotel.
func ViewRooms(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}
	var rooms []models.Rooms
	if err := Init.DB.Preload("RoomCategory").Preload("Cancellation").Preload("Hotels").Where("owner_username = ?",username).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(rooms) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no rooms are available"})
	} else {
		c.JSON(http.StatusOK, gin.H{"rooms": rooms})
	}
}

// ViewSpecificRoom returns details of a specific room.
func ViewSpecificRoom(c *gin.Context) {
	roomIDStr := c.DefaultQuery("id", "")
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
	if err := Init.DB.Preload("Cancellation").Preload("RoomCategory").Where("id = ?", uint(roomID)).First(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}

// DeleteRoom deletes a room.
func DeleteRoom(c *gin.Context) {
	var room models.Rooms
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}
	roomIDStr := c.DefaultQuery("id", "")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "convert error"})
		return
	}
	if err := Init.DB.Where("ownerusername = ? AND room_id = ?", username, roomID).Delete(&room).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "delete error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "delete success"})
}

// RoomAvailability toggles room availability.
func RoomAvailability(c *gin.Context) {
	roomIDStr := c.DefaultQuery("id", "")
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
	if err := Init.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Room not found",
		})
		return
	}

	room.IsAvailable = !room.IsAvailable

	if err := Init.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save room availability",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "availability updated"})
}
