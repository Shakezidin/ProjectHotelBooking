package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewCoupons returns a list of all coupons.
func ViewCoupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := Init.DB.Find(&coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while fetching coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

// CreateCoupon creates a new coupon.
func CreateCoupon(c *gin.Context) {
	var req struct {
		CouponCode string `json:"couponCode" binding:"required"`
		Discount   int    `json:"discount" binding:"required"`
		MinVal     int    `json:"minVal" binding:"required"`
		MaxVal     int    `json:"maxVal" binding:"required"`
		ExpireAt   string `json:"expireAt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	// Convert the ExpireAt date to time.Time
	expiresAt, err := time.Parse(layout, req.ExpireAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date conversion error"})
		return
	}

	var existingCoupon models.Coupon
	if result := Init.DB.Where("coupon_code = ?", req.CouponCode).First(&existingCoupon); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The coupon already exists"})
		return
	}

	coupon := models.Coupon{
		CouponCode: req.CouponCode,
		Discount:   req.Discount,
		MinValue:   req.MinVal,
		MaxValue:   req.MaxVal,
		ExpiresAt:  expiresAt,
	}

	if result := Init.DB.Create(&coupon); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Coupon added"})
}

// BlockCoupon toggles the 'isBlock' field of a coupon.
func BlockCoupon(c *gin.Context) {
	couponID := c.Query("id")

	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	coupon.IsBlock = !coupon.IsBlock

	if err := Init.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Coupon status updated"})
}

// GetCoupon retrieves a coupon by ID.
func GetCoupon(c *gin.Context) {
	couponID := c.Query("id")

	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

// UpdateCoupon updates a coupon by ID.
func UpdateCoupon(c *gin.Context) {
	couponID := c.Query("id") // Use Param() to get the ID from the route

	var req struct {
		CouponCode string `json:"couponCode" binding:"required"`
		Discount   int    `json:"discount" binding:"required"`
		MinVal     int    `json:"minVal" binding:"required"`
		MaxVal     int    `json:"maxVal" binding:"required"`
		ExpireAt   string `json:"expireAt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"

	ExpireAt, err := time.Parse(layout, req.ExpireAt)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	// Update the coupon in the database
	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	// Update the coupon fields
	coupon.CouponCode = req.CouponCode
	coupon.Discount = req.Discount
	coupon.MinValue = req.MinVal
	coupon.MaxValue = req.MaxVal
	coupon.ExpiresAt = ExpireAt

	if err := Init.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Coupon updated"})
}

// DeleteCoupon deletes a coupon by ID.
func DeleteCoupon(c *gin.Context) {
	couponIDStr := c.Query("id")
	if couponIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon ID is missing"})
		return
	}

	couponID, err := strconv.Atoi(couponIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Coupon ID"})
		return
	}

	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	if err := Init.DB.Delete(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Coupon deleted"})
}
