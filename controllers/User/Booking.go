package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func OfflinePayment(c *gin.Context) {
	roomIDStr := c.Query("roomid")
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
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching user"})
		return
	}

	fromdateStr, err := initiializer.ReddisClient.Get(context.Background(), "fromdate").Result()
	fmt.Println(fromdateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}

	todateStr, err := initiializer.ReddisClient.Get(context.Background(), "todate").Result()
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
	code, _ := initiializer.ReddisClient.Get(context.Background(), "couponcode").Result()
	var coupon models.Coupon
	if err := Init.DB.Where("coupen_code = ?", code).First(&coupon).Error; err != nil {
		c.JSON(400, gin.H{"msg": "No coupons Applied"})
	}

	amountStr, err := initiializer.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
	}
	var room models.Rooms

	if err := Init.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(400, gin.H{"error": "error while fetching room"})
		return
	}

	err = initiializer.ReddisClient.Set(context.Background(), "roomid", room.ID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}
	var owner models.Owner
	owenrUsername := room.OwnerUsername

	if err := Init.DB.Where("user_name = ?", owenrUsername).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching owner"})
		return
	}

	duration := toDate.Sub(fromDate)
	days := duration.Hours() / 24

	totalPrice := room.DiscountPrice * days
	owneramount := totalPrice - ((30 / 100) * totalPrice)

	booking.UserID = user.User_Id
	booking.HotelID = room.HotelsId
	booking.RoomID = uint(roomID)
	booking.OwnerID = owner.ID
	booking.RoomNo = uint(room.RoomNo)
	booking.CheckInDate = fromDate
	booking.CheckOutDate = toDate
	booking.PaymentMethod = "PAY AT HOTEL"
	booking.PaymentAmount = float64(amount)
	booking.TotalDays = uint(days)
	booking.AdminAmount = float64(amount)
	booking.OwnerAmount = totalPrice - ((30 / 100) * totalPrice)
	booking.RoomCategoryID = room.RoomCategoryId
	booking.BookedAt = time.Now()

	if err := Init.DB.Create(&booking).Error; err != nil {
		c.JSON(400, gin.H{"error": "booking creation error"})
		return
	}

	var availablerooms models.AvailableRoom
	availablerooms.RoomID = room.ID
	availablerooms.CheckIn = fromDate
	availablerooms.Checkout = toDate
	availablerooms.IsAvailable = false

	Init.DB.Create(&availablerooms)

	owner.Revenue += int(owneramount)
	Init.DB.Save(&owner)

	adminRevenue := models.Revenue{}
	Init.DB.First(&adminRevenue, "owner_id= ?", owner.ID)
	if adminRevenue.OwnerId == 0 {
		newAdminRevenue := models.Revenue{
			OwnerId:      owner.ID,
			AdminRevenue: amount,
		}
		Init.DB.Create(&newAdminRevenue)
	} else {
		adminRevenue.AdminRevenue += amount
		Init.DB.Save(&adminRevenue)
	}

	var usedcoupen models.UsedCoupen

	usedcoupen.CouponId = coupon.ID
	usedcoupen.Username = username

	if err := Init.DB.Create(&usedcoupen).Error; err != nil {
		c.JSON(400, gin.H{"error": "usedcoupon creation error"})
		return
	}

	c.JSON(200, gin.H{"status": "booking success"})
}
