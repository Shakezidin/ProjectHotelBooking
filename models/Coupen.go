package models

import (
	"time"

	"gorm.io/gorm"
)

//Coupon Model
type Coupon struct {
	gorm.Model
	CouponCode string    `gorm:"not null"`
	Discount   int       `gorm:"not null"`
	MinValue   int       `gorm:"not null"`
	MaxValue   int       `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	IsBlock    bool      `gomr:"default:false"`
}

//UsedCoupon Model
type UsedCoupon struct {
	gorm.Model
	UserID uint
	CouponID uint
}
