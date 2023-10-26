package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// OfflinePayment handles offline payment booking.
func OfflinePayment(c *gin.Context) {
	roomIDStr := c.Query("id")
	if roomIDStr == "" {
		c.JSON(400, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	var booking models.Booking
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "username does not exist"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching user"})
		return
	}

	fromdateStr, err := initializer.ReddisClient.Get(context.Background(), "fromdate").Result()
	fmt.Println(fromdateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}

	todateStr, err := initializer.ReddisClient.Get(context.Background(), "todate").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'todate' from Redis client"})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromdateStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid fromdate format"})
		return
	}

	toDate, err := time.Parse("2006-01-02", todateStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid toDate format"})
		return
	}

	couponIDStr, _ := initializer.ReddisClient.Get(context.Background(), "couponID").Result()
	couponID, _ := strconv.Atoi(couponIDStr)

	amountStr, err := initializer.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string conversion"})
	}

	var room models.Rooms
	if err := Init.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching room"})
		return
	}

	err = initializer.ReddisClient.Set(context.Background(), "roomid", room.ID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	var owner models.Owner
	ownerUsername := room.OwnerUsername
	if err := Init.DB.Where("user_name = ?", ownerUsername).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching owner"})
		return
	}

	duration := toDate.Sub(fromDate)
	days := duration.Hours() / 24
	totalPrice := room.DiscountPrice * days
	ownerAmount := totalPrice - ((30 / 100) * totalPrice)

	booking.UserID = user.ID
	booking.HotelID = room.HotelsID
	booking.RoomID = uint(roomID)
	booking.OwnerID = owner.ID
	booking.RoomNo = uint(room.RoomNo)
	booking.CheckInDate = fromDate
	booking.CheckOutDate = toDate
	booking.PaymentMethod = "PAY AT HOTEL"
	booking.PaymentAmount = float64(amount)
	booking.TotalDays = uint(days)
	booking.AdminAmount = float64(amount)
	booking.OwnerAmount = ownerAmount
	booking.RoomCategoryID = room.RoomCategoryID
	booking.CancellationID = room.CancellationID
	booking.BookedAt = time.Now()

	if err := Init.DB.Create(&booking).Error; err != nil {
		c.JSON(400, gin.H{"error": "booking creation error"})
		return
	}

	var availableRooms models.AvailableRoom
	availableRooms.BookingID = booking.ID
	availableRooms.RoomID = room.ID
	availableRooms.CheckIn = fromDate
	availableRooms.CheckOut = toDate

	Init.DB.Create(&availableRooms)

	owner.Revenue += uint(ownerAmount)
	Init.DB.Save(&owner)

	adminRevenue := models.Revenue{}
	Init.DB.First(&adminRevenue, "owner_id= ?", owner.ID)
	if adminRevenue.OwnerID == 0 {
		newAdminRevenue := models.Revenue{
			OwnerID:      owner.ID,
			AdminRevenue: uint(amount),
		}
		Init.DB.Create(&newAdminRevenue)
	} else {
		adminRevenue.AdminRevenue += uint(amount)
		Init.DB.Save(&adminRevenue)
	}

	var usedCoupon models.UsedCoupon
	usedCoupon.CouponID = uint(couponID)
	usedCoupon.UserID = user.ID

	if err := Init.DB.Create(&usedCoupon).Error; err != nil {
		c.JSON(400, gin.H{"error": "usedcoupon creation error"})
		return
	}

	c.JSON(200, gin.H{"status": "booking success"})
}
