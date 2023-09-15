package HotelOwner

import (
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
	var catagory []models.Cancellation
	if err := Init.DB.Find(&catagory).Error; err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"catagories": catagory,
	})
}

// >>>>>>>>>>>>>> room details needed to fill <<<<<<<<<<<<<<<<<<<<<<<<<<

func AddRoom(c *gin.Context) {
	var room struct {
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Adults      int     `json:"adults"`
		Children    int     `json:"children"`
		Bed         string  `json:"bed"`
		Images      string  `json:"images"`
		NoOfRooms   int     `json:"number_of_rooms"`
		IsAvailable bool    `json:"is_available"`
		Discount    float64 `json:"discount"`
	}
	c.JSON(200, gin.H{
		"room": room,
	})
}

// >>>>>>>>>>>>>> Add Room Details <<<<<<<<<<<<<<<<<<<<<<<<<<

func AddingRoom(c *gin.Context) {
	cancellationIDStr := c.DefaultQuery("cancellation_id", "")
	hotelIDStr := c.DefaultQuery("hotel_id", "")
	roomCategoryIDStr := c.DefaultQuery("room_category_id", "")

	// Convert string IDs to uint
	cancellationID, err := strconv.ParseUint(cancellationIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid cancellation_id"})
		return
	}

	hotelID, err := strconv.ParseUint(hotelIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid hotel_id"})
		return
	}

	roomCategoryID, err := strconv.ParseUint(roomCategoryIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid room_category_id"})
		return
	}
	var room models.Room

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	validationErr := validate.Struct(room)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"msg": validationErr,
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

	room.Cancellation_Id = uint(cancellationID)
	room.ID = uint(hotelID)
	room.Category_Id = uint(roomCategoryID)
	room.OwnerUsername = username
	room.DiscountPrice = room.Price - (room.Price * room.Discount / 100)
	if err := Init.DB.Create(&room).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}

// >>>>>>>>>>>>>> Edit Room <<<<<<<<<<<<<<<<<<<<<<<<<<

func EditRoom(c *gin.Context) {
	roomIDStr := c.DefaultQuery("roomid", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	hotelIdStr := c.DefaultQuery("hotelid", "")
	if hotelIdStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Room
	if err := Init.DB.Where("hotel_id = ? AND room_id = ?", uint(hotelId), uint(roomId)).First(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

	var updatedRoom struct {
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Adults      uint     `json:"adults"`
		Children    uint     `json:"children"`
		Bed         string   `json:"bed"`
		Images      string   `json:"images"`
		NoOfRooms   uint     `json:"number_of_rooms"`
		IsAvailable bool     `json:"is_available"`
		Discount    float64  `json:"discount"`
		Fecilities  []string `json:"facilities" gorm:"type:jsonb"`
	}
	if err := c.BindJSON(&updatedRoom); err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}

	room.Description = updatedRoom.Description
	room.Price = updatedRoom.Price
	room.Adults = updatedRoom.Adults
	room.Children = updatedRoom.Children
	room.Bed = updatedRoom.Bed
	room.Images = updatedRoom.Images
	room.NoOfRooms = updatedRoom.NoOfRooms
	room.IsAvailable = updatedRoom.IsAvailable
	room.Discount = updatedRoom.Discount
	room.Fecilities = updatedRoom.Fecilities
	room.DiscountPrice = room.Price - (room.Price * room.Discount / 100)

	result := Init.DB.Save(&room)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

// >>>>>>>>>>>>>> view Rooms <<<<<<<<<<<<<<<<<<<<<<<<<<

func ViewRooms(c *gin.Context) {
	hotelIdStr := c.DefaultQuery("hotelid", "")
	if hotelIdStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIdStr)
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
	var rooms []models.Room
	if err := Init.DB.Where("hotel_id = ? AND ownerusername = ?", uint(hotelId), username).Find(&rooms).Error; err != nil {
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
	roomIdStr := c.DefaultQuery("roomid", "")
	if roomIdStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	hotelIdStr := c.DefaultQuery("hotelid", "")
	if hotelIdStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	hotelId, err := strconv.Atoi(hotelIdStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	var room models.Room
	if err := Init.DB.Where("hotel_id = ? AND room_id = ?", uint(hotelId), uint(roomId)).First(&room).Error; err != nil {
		c.JSON(500, gin.H{
			"msg": err.Error(),
		})
		c.Abort()
		return
	}

}

func DeleteRoom(c *gin.Context) {
	var room models.Room
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username didnt get"})
		return
	}
	roomIdStr := c.DefaultQuery("roomid", "")
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

func RoomAvailability(c *gin.Context) {
	roomIDStr := c.DefaultQuery("hotelid", "")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	roomId, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var room models.Room
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

	c.Status(200)
}
