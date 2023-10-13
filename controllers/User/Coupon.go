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

func Book(c *gin.Context) {
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

	err = initiializer.ReddisClient.Set(context.Background(), "roomid", roomIDStr, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
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

	duration := toDate.Sub(fromDate)
	days := int(duration.Hours() / 24)

	var room models.Rooms
	if err := Init.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetching rooms"})
		return
	}

	totalPrice := days * int(room.DiscountPrice)

	GSTPercentage := 18.0 // Use a floating-point number for the percentage
	GSTAmount := totalPrice * int(GSTPercentage/100.0)
	payableamount := totalPrice + (totalPrice * 18 / 100)

	err = initiializer.ReddisClient.Set(context.Background(), "Amount", payableamount, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{
		"status":        "success",
		"roomPrice":     room.Price,
		"discount":      room.Discount,
		"TotalAmount":   totalPrice,
		"GSTAmount":     GSTAmount,
		"PayableAmount": payableamount,
	})
}

func ApplyCoupon(c *gin.Context) {
	code := c.Query("coupencode")
	if code == "" {
		c.JSON(400, gin.H{"error": "coupon query parameter is missing"})
		return
	}
	amountStr, err := initiializer.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'fromdate' from Redis client"})
		return
	}
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
	}
	var coupon models.Coupon
	var usedcoupen models.UsedCoupen

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	if err := Init.DB.Where("coupen_code = ?", code).First(&coupon).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching coupen"})
		return
	}

	currentTime := time.Now()
	if coupon.ExpiresAt.Before(currentTime) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon expired"})
		return
	}

	if amount < coupon.MinValue {
		c.JSON(400, gin.H{"Error": "Minimum amount required"})
		return
	}

	if amount > coupon.MaxValue {
		c.JSON(400, gin.H{"Error": "amount above max amount"})
		return
	}

	result := Init.DB.Where("coupon_id = ? AND username = ?", coupon.ID, username).First(&usedcoupen)
	if result.RowsAffected > 0 {
		c.JSON(400, gin.H{
			"error": "coupen already used ",
		})
		return
	}

	updatedTotal := amount - coupon.Discount
	err = initiializer.ReddisClient.Set(context.Background(), "Amount", updatedTotal, 1*time.Hour).Err()
	err = initiializer.ReddisClient.Set(context.Background(), "couponcode", code, 1*time.Hour).Err()
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

func Coupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := Init.DB.Where("is_block = ?", false).Find(&coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while fetching coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

func Wallet(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	var user models.User
	if err := Init.DB.Where("user_name = ?", username).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching user"})
		return
	}

	var wallet models.Wallet

	if err := Init.DB.Where("user_id = ?", user.User_Id).First(&wallet).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Error while fetching wallet"})
		return
	}

	err = initiializer.ReddisClient.Set(context.Background(), "WalletId", wallet.ID, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{"success": wallet})
}

func Applaywallet(c *gin.Context) {
	amountStr, err := initiializer.ReddisClient.Get(context.Background(), "Amount").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'amount' from Redis client"})
		return
	}
	amountt, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
	}
	walletIdstr, err := initiializer.ReddisClient.Get(context.Background(), "WalletId").Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting 'walletId' from Redis client"})
		return
	}
	walletId, err := strconv.Atoi(walletIdstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "string convertion"})
	}
	var wallet models.Wallet
	if err := Init.DB.Where("id = ?", uint(walletId)).First(&wallet).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching wallet"})
		return
	}

	amount := float64(amountt)
	if wallet.Balance <= 0 {
		c.JSON(400, gin.H{"error": "no balnce are there"})
		return
	} else {
		var balnc float64
		if amount > wallet.Balance {
			balnc = 0
		} else {
			balnc = amount - wallet.Balance
		}
		amount = amount - wallet.Balance
		wallet.Balance = balnc
		if err := Init.DB.Save(&wallet).Error; err != nil {
			c.JSON(400, gin.H{"error": "Error while updating wallet"})
			return
		}
	}

	var transaction models.Transaction
	if err := Init.DB.Where("u_ser_id =  ?", wallet.User_Id).First(&transaction).Error; err != nil {
		c.JSON(400, gin.H{"error": "transaction fetching error"})
		return
	}

	transaction.Date = time.Now()
	transaction.Details = "Booked room in"
	transaction.Amount = amount

	if err := Init.DB.Create(&transaction).Error; err != nil {
		c.JSON(400, gin.H{"error": "transaction adding error"})
		return
	}
	err = initiializer.ReddisClient.Set(context.Background(), "Amount", amount, 1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting in Redis client"})
		return
	}

	c.JSON(200, gin.H{
		"status":         "wallet applied",
		"wallet balance": wallet.Balance,
		"amount":         amount,
	})
}
