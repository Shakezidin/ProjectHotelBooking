package user

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"github.com/shaikhzidhin/initializer"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

type pageVariables struct {
	OrderID string
}

// RazorpayPaymentGateway handles RazorPay payments.
func RazorpayPaymentGateway(c *gin.Context) {
	UserIDStr := c.Query("id")
	UserID, err := strconv.Atoi(UserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
	}

	amountStr, err := initializer.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
	}

	var user models.User
	if err := Init.DB.Where("id = ?", uint(UserID)).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))

	amountInPaisa := int(amount) * 100
	data := map[string]interface{}{
		"amount":   amountInPaisa,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Printf("Problem getting repository information: %v\n", err)
		os.Exit(1)
	}

	value := body["id"]
	str := value.(string)

	homepageVariables := pageVariables{
		OrderID: str,
	}

	c.HTML(http.StatusOK, "app.html", gin.H{
		"userID":      user.ID,
		"totalPrice":  amountInPaisa / 100,
		"total":       amountInPaisa,
		"orderID":     homepageVariables.OrderID,
		"email":       user.Email,
		"phoneNumber": user.Phone,
	})
}

// RazorpaySuccess handles successful RazorPay payments.
func RazorpaySuccess(c *gin.Context) {
	userID := c.Query("user_id")
	UserID, _ := strconv.Atoi(userID)
	orderID := c.Query("order_id")
	paymentID := c.Query("payment_id")
	signature := c.Query("signature")
	paymentAmount := c.Query("total")
	amount, _ := strconv.Atoi(paymentAmount)
	rPay := models.RazorPay{
		UserID:          uint(UserID),
		RazorPaymentID:  paymentID,
		Signature:       signature,
		RazorPayOrderID: orderID,
		AmountPaid:      float64(amount),
	}

	result := Init.DB.Create(&rPay)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	fromDateStr, err := initializer.ReddisClient.Get(context.Background(), "fromdate").Result()
	fmt.Println(fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}

	toDateStr, err := initializer.ReddisClient.Get(context.Background(), "todate").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'todate' from Redis client"})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromdate format"})
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toDate format"})
		return
	}
	roomIDStr, err := initializer.ReddisClient.Get(context.Background(), "roomid").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'roomid' from Redis client"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
	}
	var room models.Rooms
	if err := Init.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching room"})
		return
	}

	var owner models.Owner
	if err := Init.DB.Where("user_name = ?", room.OwnerUsername).First(&owner).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching owner"})
	}
	duration := toDate.Sub(fromDate)
	days := duration.Hours() / 24

	ownerAmount := (room.DiscountPrice * days) * 0.7

	var booking models.Booking
	booking.UserID = uint(UserID)
	booking.HotelID = room.HotelsID
	booking.RoomID = room.ID
	booking.OwnerID = owner.ID
	booking.RoomNo = uint(room.RoomNo)
	booking.CheckInDate = fromDate
	booking.CheckOutDate = toDate
	booking.PaymentMethod = "RAZOR PAY"
	booking.AdminAmount = float64(amount)
	booking.PaymentAmount = float64(amount)
	booking.OwnerAmount = ownerAmount
	booking.TotalDays = uint(days)
	booking.RoomCategoryID = room.RoomCategoryID
	booking.CancellationID = room.CancellationID
	booking.BookedAt = time.Now()

	if err := Init.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking error"})
		return
	}

	couponIDStr, _ := initializer.ReddisClient.Get(context.Background(), "couponID").Result()
	couponID, _ := strconv.Atoi(couponIDStr)

	var usedCoupon models.UsedCoupon
	usedCoupon.CouponID = uint(couponID)
	usedCoupon.UserID = uint(UserID)

	if err := Init.DB.Create(&usedCoupon).Error; err != nil {
		c.JSON(400, gin.H{"error": "usedcoupon creation error"})
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

	c.JSON(http.StatusOK, gin.H{"status": true})
}

// SuccessPage renders the success page.
func SuccessPage(c *gin.Context) {
	pID := c.Query("id")
	fmt.Println(pID)
	fmt.Println("Fully successful")

	c.HTML(http.StatusOK, "success.html", gin.H{
		"paymentID": pID,
	})
}
