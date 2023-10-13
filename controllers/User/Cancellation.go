package user

// import (
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/shaikhzidhin/models"
// )

// func Cancellation(c *gin.Context) {
// 	// Parse request body
// 	var requestBody struct {
// 		BookingID        uint   `json:"bookingId"`
// 		RoomCancellation string `json:"roomCancellation"`
// 	}
// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the booking
// 	var booking models.Booking
// 	if err := db.First(&booking, requestBody.BookingID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
// 		return
// 	}

// 	// Parse booking dates
// 	checkin := booking.CheckInDate
// 	checkout := booking.CheckOutDate
// 	currentDate := time.Now()
// 	userID := booking.User

// 	// Calculate amount to refund
// 	amountRefund := booking.PaymentAmount

// 	switch requestBody.RoomCancellation {
// 	case "Free cancellation upto 24hrs before checkin date":
// 		// Check the time difference
// 		timeDifference := checkin.Sub(currentDate)
// 		hoursDifference := int(timeDifference.Hours())

// 		if hoursDifference >= 24 {
// 			// Perform refund and return success
// 			err := userWallet(userID, amountRefund, booking.ID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.Status(http.StatusOK)
// 			return
// 		}

// 	case "Non-Refundable":
// 		// Perform cancellation (no refund) and return unauthorized status
// 		err := cancelBooking(booking.ID, booking.Room)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		c.Status(http.StatusUnauthorized)
// 		return

// 	case "Canceling within 7 days before checkin":
// 		// Check the time difference in days
// 		daysDifference := int(checkin.Sub(currentDate).Hours() / 24)

// 		if daysDifference >= 7 {
// 			// Perform refund and return success
// 			err := userWallet(userID, amountRefund, booking.ID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.Status(http.StatusOK)
// 			return
// 		}

// 	case "Free Cancellation before checkin date":
// 		if currentDate.Before(checkin) {
// 			// Perform refund and return success
// 			err := userWallet(userID, amountRefund, booking.ID)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 				return
// 			}
// 			c.Status(http.StatusOK)
// 			return
// 		}
// 	}

// 	// If the cancellation policy doesn't meet any conditions, cancel the booking (no refund)
// 	err := cancelBooking(booking.ID, booking.Room)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.Status(http.StatusUnauthorized)
// }

// func userWallet(userID, amount, bookingId) {
// 	date := time.Now()
// 	formattedDate := date.Format("2006-01-02")
// 	adminAmount := (15 / 100) * request.Amount
// 	ownerAmount := (85 / 100) * request.Amount

// 	// Update the user's wallet balance and transaction
// 	var user User
// 	if err := db.First(&user, request.UserID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	db.Model(&user).Update("Wallet.Balance", user.Wallet.Balance+request.Amount)
// 	db.Model(&user).Association("Wallet.Transactions").Append(&Transaction{
// 		Date:    formattedDate,
// 		Details: "Cancelled the booking",
// 		Amount:  request.Amount,
// 	})

// 	// Update the booking
// 	var booking Booking
// 	if err := db.First(&booking, request.BookingID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
// 		return
// 	}
// 	booking.Cancel = true
// 	booking.Refund = true
// 	db.Save(&booking)

// 	// Update the room
// 	var room Room
// 	if err := db.First(&room, booking.Room).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
// 		return
// 	}

// 	roomIndex := -1
// 	for i, roomInfo := range room.AvailableRooms {
// 		if roomInfo.RoomNo == booking.RoomNo {
// 			roomIndex = i
// 			break
// 		}
// 	}

// 	if roomIndex != -1 {
// 		room.AvailableRooms[roomIndex].RemoveCheckinCheckout(booking.CheckInDate, booking.CheckOutDate)
// 		db.Save(&room)
// 	}

// 	// Update the hotel
// 	var hotel Hotel
// 	if err := db.First(&hotel, booking.Hotel).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
// 		return
// 	}
// 	db.Model(&hotel).Update("Revenue", hotel.Revenue-ownerAmount)

// 	// Update the owner
// 	var owner Owner
// 	if err := db.First(&owner, hotel.Owner).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Owner not found"})
// 		return
// 	}
// 	db.Model(&owner).Update("Revenue", owner.Revenue-ownerAmount)

// 	// Update admin revenue (assuming you have an AdminRevenue model)
// 	var adminRevenue AdminRevenue
// 	if db.Where("Owner = ?", owner.ID).First(&adminRevenue).RecordNotFound() {
// 		// AdminRevenue record not found, create a new one
// 		adminRevenue.Owner = owner.ID
// 		adminRevenue.Revenue = -adminAmount
// 		db.Create(&adminRevenue)
// 	} else {
// 		// AdminRevenue record found, update the revenue
// 		db.Model(&adminRevenue).Update("Revenue", adminRevenue.Revenue-adminAmount)
// 	}

// 	c.Status(http.StatusOK)
// }

// func CancelBooking(c *gin.Context) {
// 	var request struct {
// 		BookingID uint `json:"bookingId"`
// 		RoomID    uint `json:"roomId"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Update the booking
// 	var booking Booking
// 	if err := db.First(&booking, request.BookingID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
// 		return
// 	}

// 	booking.Cancel = true
// 	db.Save(&booking)

// 	// Update the room
// 	var room Room
// 	if err := db.First(&room, request.RoomID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
// 		return
// 	}

// 	roomIndex := -1
// 	for i, roomInfo := range room.AvailableRooms {
// 		if roomInfo.RoomNo == booking.RoomNo {
// 			roomIndex = i
// 			break
// 		}
// 	}

// 	if roomIndex != -1 {
// 		room.AvailableRooms[roomIndex].RemoveCheckinCheckout(booking.CheckInDate, booking.CheckOutDate)
// 		db.Save(&room)
// 	}

// 	c.Status(http.StatusOK)
// }
