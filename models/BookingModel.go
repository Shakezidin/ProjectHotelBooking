package models

import (
	"time"

	"gorm.io/gorm"
)

// Booking Model
type Booking struct {
	gorm.Model
	UserID         uint      `gorm:"not null"`
	User           User      `gorm:"ForeignKey:UserID"`
	HotelID        uint      `gorm:"not null"`
	Hotels         Hotels    `gorm:"ForeignKey:HotelID"`
	RoomID         uint      `gorm:"not null"`
	Rooms          Rooms     `gorm:"ForeignKey:RoomID"`
	OwnerID        uint      `gorm:"not null"`
	Owner          Owner     `gorm:"ForeignKey:OwnerID"`
	RoomNo         uint      `gorm:"not null"`
	CheckInDate    time.Time `gorm:"not null"`
	CheckOutDate   time.Time `gorm:"not null"`
	PaymentMethod  string
	PaymentAmount  float64
	TotalDays      uint
	AdminAmount    float64 `gorm:"not null"`
	OwnerAmount    float64 `gorm:"not null"`
	CancellationID uint    `gorm:"not null"`
	Cancellation   Cancellation
	RoomCategoryID uint         `gorm:"not null"`
	RoomCategory   RoomCategory `gorm:"ForeignKey:RoomCategoryID"`
	Review         bool         `gorm:"default:false"`
	Report         bool         `gorm:"default:false"`
	Refund         bool         `gorm:"default:false"`
	Cancel         bool         `gorm:"default:false"`
	AdminResponse  string       `gorm:"default:'No Action'"`
	BookedAt       time.Time    `gorm:"default:CURRENT_TIMESTAMP"`
	Cancelled      string
}

// RazorPay Model
type RazorPay struct {
	UserID          uint    `JSON:"userid"`
	RazorPaymentID  string  `JSON:"razorpaymentid" gorm:"primaryKey"`
	RazorPayOrderID string  `JSON:"razorpayorderid"`
	Signature       string  `JSON:"signature"`
	AmountPaid      float64 `JSON:"amountpaid"`
}
