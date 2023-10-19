package user

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// CalculateAmountForDays calculates the amount for booking based on selected dates and room.
func CalculateAmountForDays(c *gin.Context) {
	roomIDStr := c.Query("id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomid query parameter is missing"})
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "convert error"})
		return
	}

	err = Init.ReddisClient.Set(context.Background(), "roomid", roomIDStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	fromdateStr, err := Init.ReddisClient.Get(context.Background(), "fromdate").Result()
	fmt.Println(fromdateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}

	todateStr, err := Init.ReddisClient.Get(context.Background(), "todate").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'todate' from Redis client"})
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromdateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fromdate format"})
		return
	}

	toDate, err := time.Parse("2006-01-02", todateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid toDate format"})
		return
	}

	duration := toDate.Sub(fromDate)
	days := int(duration.Hours() / 24)

	var room models.Rooms
	if err := Init.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetching rooms"})
		return
	}

	totalPrice := days * int(room.DiscountPrice)

	GSTPercentage := 18.0 // Use a floating-point number for the percentage
	GSTAmount := (GSTPercentage / 100.0) * float64(totalPrice)
	payableAmount := totalPrice + int(GSTAmount)

	err = Init.ReddisClient.Set(context.Background(), "Amount", payableAmount, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"roomPrice":     room.Price,
		"discount":      room.Discount,
		"TotalAmount":   totalPrice,
		"GSTAmount":     int(GSTAmount),
		"PayableAmount": payableAmount,
	})
}

// ViewNonBlockedCoupons retrieves non-blocked coupons.
func ViewNonBlockedCoupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := Init.DB.Where("is_block = ?", false).Find(&coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while fetching coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

// ApplyCoupon applies a coupon to the booking.
func ApplyCoupon(c *gin.Context) {
	couponID, _ := strconv.Atoi(c.Query("id"))
	amountStr, err := Init.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'Amount' from Redis client: " + err.Error()})
		return
	}

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}

	CouponIDstr, err := Init.ReddisClient.Get(context.Background(), "couponID").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "No coupon used"})
	}
	if CouponIDstr != "" {
		oldcouponID, err := strconv.Atoi(CouponIDstr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
			return
		}

		if oldcouponID == couponID {
			c.JSON(400, gin.H{"error": "alredy used coupon"})
			return
		}
	}
	var coupon models.Coupon
	var usedCoupon models.UsedCoupon

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "username not found"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "user fetching error"})
		return
	}

	if err := Init.DB.Where("id = ?", couponID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching coupon"})
		return
	}

	currentTime := time.Now()
	if coupon.ExpiresAt.Before(currentTime) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon expired"})
		return
	}

	if amount < coupon.MinValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum amount required"})
		return
	}

	if amount > coupon.MaxValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount above max amount"})
		return
	}

	result := Init.DB.Where("coupon_id = ? AND user_id = ?", coupon.ID, user.ID).First(&usedCoupon)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon already used"})
		return
	}

	updatedTotal := amount - coupon.Discount
	err = Init.ReddisClient.Set(context.Background(), "Amount", updatedTotal, 1*time.Hour).Err()
	err = Init.ReddisClient.Set(context.Background(), "couponID", couponID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "coupon applied",
		"couponDiscount": coupon.Discount,
		"current total":  updatedTotal,
	})
}

// ViewWallet retrieves user's wallet information.
func ViewWallet(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching user"})
		return
	}

	var wallet models.Wallet

	if err := Init.DB.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching wallet"})
		return
	}

	err = Init.ReddisClient.Set(context.Background(), "WalletId", wallet.ID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": wallet})
}

// ApplyWallet applies user's wallet balance.
func ApplyWallet(c *gin.Context) {
	amountStr, err := Init.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amountt, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}
	walletIDStr, err := Init.ReddisClient.Get(context.Background(), "WalletId").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'walletId' from Redis client"})
		return
	}
	walletID, err := strconv.Atoi(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "string conversion"})
		return
	}

	var wallet models.Wallet
	if err := Init.DB.Where("id = ?", uint(walletID)).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching wallet"})
		return
	}

	amount := float64(amountt)
	if wallet.Balance <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no balance available"})
		return
	}

	var balance float64
	if amount > wallet.Balance {
		balance = 0
	} else {
		balance = amount - wallet.Balance
	}
	amount = amount - wallet.Balance
	wallet.Balance = balance
	if err := Init.DB.Save(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while updating wallet"})
		return
	}

	var transaction models.Transaction

	transaction.Date = time.Now()
	transaction.Details = "Booked room in"
	transaction.Amount = amount
	transaction.UserID = wallet.UserID

	if err := Init.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transaction adding error"})
		return
	}

	err = Init.ReddisClient.Set(context.Background(), "Amount", amount, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "wallet applied",
		"wallet balance": wallet.Balance,
		"amount":         amount,
	})
}
