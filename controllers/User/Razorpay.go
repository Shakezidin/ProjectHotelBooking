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
	"github.com/shaikhzidhin/initiializer"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

type pageVariables struct {
	OrderId string
}

func Razorpay(c *gin.Context) {
	UserIdstr, err := initiializer.ReddisClient.Get(context.Background(), "userId").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'userid' from Redis client"})
		return
	}
	userid, err := strconv.Atoi(UserIdstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
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

	var user models.User
	if err := Init.DB.Where("user_id = ?", uint(userid)).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "User not found"})
		return
	}
	client := razorpay.NewClient("rzp_test_loKZTxH4NevSeO", "C41yIa655LTCzsbNS2Cietro")

	amountt := int(amount) * 100
	data := map[string]interface{}{
		"amount":   (amountt),
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)

	if err != nil {
		fmt.Printf("Problem is getting repository information %v\n", err)
		os.Exit(1)
	}

	value := body["id"]
	str := value.(string)

	homepagevariables := pageVariables{
		OrderId: str,
	}

	c.HTML(200, "app.html", gin.H{
		"userid":      user.User_Id,
		"totalprice":  amountt,
		"total":       amountt,
		"orderid":     homepagevariables.OrderId,
		"email":       user.Email,
		"phonenumber": user.Phone,
	})
}

// >>>>>>>>>>>>>> Razorpay <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func RazorpaySuccess(c *gin.Context) {
	userid := c.Query("user_id")
	userID, _ := strconv.Atoi(userid)
	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	paymentAmount := c.Query("total")
	amount, _ := strconv.Atoi(paymentAmount)
	Rpay := models.RazorPay{
		UserID:          uint(userID),
		RazorPaymentId:  paymentid,
		Signature:       signature,
		RazorPayOrderID: orderid,
		AmountPaid:      float64(amount),
	}

	result := Init.DB.Create(&Rpay)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error,
		})

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
	roomidstr, err := initiializer.ReddisClient.Get(context.Background(), "roomid").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'roomid' from Redis client"})
		return
	}
	roomid, err := strconv.Atoi(roomidstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
	}
	var room models.Rooms
	if err := Init.DB.Where("id = ?", roomid).First(&room).Error; err != nil {
		c.JSON(400, gin.H{"error": "error fetching room"})
	}

	var owner models.Owner
	if err := Init.DB.Where("user_name = ?", room.OwnerUsername).First(&owner).Error; err != nil {
		c.JSON(400, gin.H{"error": "error fetchin owner"})
	}
	duration := toDate.Sub(fromDate)
	days := duration.Hours() / 24

	owneramount := (room.DiscountPrice * days) * 0.7

	var booking models.Booking
	booking.UserID = uint(userID)
	booking.HotelID = room.HotelsId
	booking.RoomID = room.ID
	booking.OwnerID = owner.ID
	booking.RoomNo = uint(room.RoomNo)
	booking.CheckInDate = fromDate
	booking.CheckOutDate = toDate
	booking.PaymentMethod = "RAZOR PAY"
	booking.AdminAmount = float64(amount)
	booking.PaymentAmount = float64(amount)
	booking.OwnerAmount = owneramount
	booking.TotalDays = uint(days)
	booking.RoomCategoryID = room.RoomCategoryId
	booking.BookedAt = time.Now()

	if err := Init.DB.Create(&booking).Error; err != nil {
		c.JSON(400, gin.H{"error": "booking error"})
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

	c.JSON(200, gin.H{"status": true})
}

func Success(c *gin.Context) {
	pid := c.Query("id")
	fmt.Println(pid)
	fmt.Printf("Fully success")

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
