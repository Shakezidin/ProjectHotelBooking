package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func Book(c *gin.Context) {
	roomId := c.GetUint("roomid")

	fromdatestr := c.GetString("fromdate")
	if fromdatestr == "" {
		c.JSON(400, gin.H{"error": "from date query parameter is missing"})
		return
	}
	fromDate, err := time.Parse("2006-01-02", fromdatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	todatestr := c.GetString("todate")
	if todatestr == "" {
		c.JSON(400, gin.H{"error": "To date query parameter is missing"})
		return
	}
	toDate, err := time.Parse("2006-01-02", todatestr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}

	duration := toDate.Sub(fromDate)
	days := duration.Hours() / 24

	var room models.Rooms
	if err := Init.DB.Where("Rooms_id = ?", roomId).First(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error while fetching rooms"})
		return
	}

	totalPrice := room.DiscountPrice * days

	c.JSON(200, gin.H{
		"status":        "success",
		"roomPrice":     room.Price,
		"discount":      room.Discount,
		"TotalAmount":   totalPrice,
		"GSTAmount":     totalPrice * (18 / 100),
		"PayableAmount": totalPrice + (totalPrice * 18 / 100),
	})
}

func Coupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := Init.DB.Where("isblock = ?", false).Find(&coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while fetching coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

func ApplyCoupon(c *gin.Context) {
	code := c.GetString("coupencode")
	TotalAmount := c.GetInt("Totalamount")
	var coupon models.Coupon
	var usedcoupen models.UsedCoupen

	header := c.Request.Header.Get("Authorization")
	username, err := Auth.Trim(header)
	if err != nil {
		c.JSON(404, gin.H{"error": "email didnt get"})
		return
	}

	if err := Init.DB.Where("coupencode = ?", code).First(&coupon).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error while fetching coupen"})
		return
	}

	currentTime := time.Now()
	if coupon.ExpiresAt.Before(currentTime) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon expired"})
		return
	}

	if TotalAmount < coupon.MinValue {
		c.JSON(400, gin.H{"Error": "Minimum amount required"})
		return
	}

	if TotalAmount > coupon.MaxValue {
		c.JSON(400, gin.H{"Error": "amount above max amount"})
		return
	}

	result := Init.DB.Where("couponId = ? AND username = ?", coupon.ID, username).First(&usedcoupen)
	if result.RowsAffected > 0 {
		c.JSON(400, gin.H{
			"error": "coupen already used ",
		})
		return
	}

	updatedTotal := TotalAmount - coupon.Discount

	c.JSON(http.StatusOK, gin.H{
		"status":         "coupon applied",
		"couponDiscount": coupon.Discount,
		"current total":  updatedTotal,
	})
}
