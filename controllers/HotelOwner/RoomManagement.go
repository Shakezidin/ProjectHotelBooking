package HotelOwner

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

// >>>>>>>>>>>>>> view Room Fecility <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewRoomfecilities(c *gin.Context) {
	var fecilities []models.RoomFecilities
	if err := Init.DB.Find(&fecilities).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"fecilities": fecilities,
	})
}

// >>>>>>>>>>>>>> view Room Cancellation <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewCancellation(c *gin.Context) {
	var cancellation []models.Cancellation
	if err := Init.DB.Find(&cancellation).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"fecilities": cancellation,
	})
}

// >>>>>>>>>>>>>> view Room catagory <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewRoomCatagory(c *gin.Context) {
	var catagory []models.RoomCategory
	if err := Init.DB.Find(&catagory).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"catagories": catagory,
	})
}

// >>>>>>>>>>>>>> Add Room Details <<<<<<<<<<<<<<<<<<<<<<<<<<
func AddingRoom(c *gin.Context) {
	numberOfRooms, _ := strconv.Atoi(c.Query("numberofrooms"))
	floorNumber, _ := strconv.Atoi(c.Query("floornumber"))
	cancellationId, _ := strconv.Atoi(c.Query("cancellationid"))
	hotelId, _ := strconv.Atoi(c.Query("id"))
	roomCatagoryId, _ := strconv.Atoi(c.Query("room_category_id"))

	var roombind models.Rooms
	if err := c.ShouldBindJSON(&roombind); err != nil {
		c.JSON(400, gin.H{
			"msg":   "binding error1",
			"error": err,
		})
		c.Abort()
		return
	}

	roomnum := 1
	for i := 1; i <= numberOfRooms; i++ {
		var room models.Rooms

		room.Adults = roombind.Adults
		room.Bed = roombind.Bed
		room.CancellationId = uint(cancellationId)
		room.Children = roombind.Children
		room.Description = roombind.Description
		room.Discount = roombind.Discount
		room.Fecility = roombind.Fecility
		room.HotelsId = uint(hotelId)
		room.Images = roombind.Images
		room.IsAvailable = roombind.IsAvailable
		room.Price = roombind.Price
		room.RoomCategoryId = uint(roomCatagoryId)
		room.RoomNo = floorNumber*100 + roomnum
		room.DiscountPrice = room.Price - (room.Price * room.Discount / 100)

		header := c.Request.Header.Get("Authorization")
		username, err := Auth.Trim(header)
		if err != nil {
			c.JSON(404, gin.H{"error": "Username not found"})
			return
		}
		room.OwnerUsername = username

		if err := Init.DB.Create(&room).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		roomnum++
	}

	c.JSON(200, gin.H{"status": "success"})
}

// >>>>>>>>>>>>>> Edit Room <<<<<<<<<<<<<<<<<<<<<<<<<<
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

// >>>>>>>>>>>>>> view Rooms <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewRooms(c *gin.Context) {
	hotelIdStr := c.DefaultQuery("id", "")
	if hotelIdStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var rooms []models.Rooms
	if err := Init.DB.Preload("RoomCategory").Preload("Cancellation").Preload("Hotels").Where("hotels_id = ?", uint(hotelId)).Find(&rooms).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(rooms) == 0 {
		c.JSON(400, gin.H{"msg": "no rooms are there"})
	} else {
		c.JSON(200, gin.H{"rooms": rooms})
	}
}

// >>>>>>>>>>>>>> view Specific Room <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewspecificRoom(c *gin.Context) {
	roomIdStr := c.DefaultQuery("id", "")
	if roomIdStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	var room models.Rooms
	if err := Init.DB.Preload("Cancellation").Preload("RoomCategory").Where("id = ?", uint(roomId)).First(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"Room": room})
}

// >>>>>>>>>>>>>>>> Delete a Room <<<<<<<<<<<<<<<<<<<<<<<<<<

func DeleteRoom(c *gin.Context) {
	var room models.Rooms
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	roomIdStr := c.DefaultQuery("id", "")
	if roomIdStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	if err := Init.DB.Where("ownerusername = ? AND room_id = ?", username, roomId).Delete(&room); err != nil {
		c.JSON(401, gin.H{"error": "delete error"})
		return
	}
	c.JSON(200, gin.H{"status": "delete success"})
}

// >>>>>>>>>>>>>> Switching Room Availability <<<<<<<<<<<<<<<<<<<<<<<<<<

func RoomAvailability(c *gin.Context) {
	roomIDStr := c.DefaultQuery("id", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Rooms
	if err := Init.DB.First(&room, roomId).Error; err != nil {
		c.JSON(404, gin.H{
			"hello": "Hot",
			"error": "Room not found",
		})
		return
	}

	room.IsAvailable = !room.IsAvailable

	if err := Init.DB.Save(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save Room availability",
		})
		return
	}

	c.JSON(200, gin.H{"status": "availabilty updated"})
}
