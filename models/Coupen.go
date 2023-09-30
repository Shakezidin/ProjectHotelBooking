package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	CoupenCode string    `gorm:"not null"`
	Discount   int       `gorm:"not null"`
	MinValue   int       `gorm:"not null"`
	MaxValue   int       `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	IsBlock    bool      `gomr:"default:false"`
}

type UsedCoupen struct {
	gorm.Model
	username string
	CouponId uint
}
