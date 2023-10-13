package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewCoupons(c *gin.Context) {
	var coupons []models.Coupon

	if err := Init.DB.Find(&coupons).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error while fetching coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

// CreateCoupon creates a new coupon
func CreateCoupon(c *gin.Context) {
	var req struct {
		CouponCode string `json:"couponCode"`
		Discount   int    `json:"discount"`
		MinVal     int    `json:"minVal"`
		MaxVal     int    `json:"maxVal"`
		ExpireAt   string `json:"expireAt"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	layout := "2006-01-02"
	// Convert the FromDate and ToDate to time.Time with default values
	expiresAt, err := time.Parse(layout, req.ExpireAt)
	if err != nil {
		c.JSON(400, gin.H{"Error": "date convertion error"})
		return
	}

	var existingCoupon models.Coupon
	if result := Init.DB.Where("coupen_code = ?", req.CouponCode).First(&existingCoupon); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The coupon already exists"})
		return
	}

	coupon := models.Coupon{
		CoupenCode: req.CouponCode,
		Discount:   req.Discount,
		MinValue:   req.MinVal,
		MaxValue:   req.MaxVal,
		ExpiresAt:  expiresAt,
	}

	if result := Init.DB.Create(&coupon); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "coupon added"})
}

// BlockCoupon toggles the 'isBlock' field of a coupon
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

	c.JSON(http.StatusOK, gin.H{"Status": "Coupon status updated"})
}

// GetCoupon retrieves a coupon by ID
func GetCoupon(c *gin.Context) {
	couponID := c.Query("id")

	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

// UpdateCoupon updates a coupon by ID
func UpdateCoupon(c *gin.Context) {
	couponID := c.Query("id")
	var req struct {
		CouponCode string `json:"couponCode"`
		Discount   int    `json:"discount"`
		MinVal     int    `json:"minVal"`
		MaxVal     int    `json:"maxVal"`
		ExpireAt   string `json:"expireAt"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := Init.DB.Model(&models.Coupon{}).Where("id = ?", couponID).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "coupon updated"})
}

// DeleteCoupon deletes a coupon by ID
func DeleteCoupon(c *gin.Context) {
	couponID := c.GetUint("id")

	var coupon models.Coupon
	if err := Init.DB.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	if err := Init.DB.Delete(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Status": "deleted"})
}
