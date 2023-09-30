package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	User          User      `gorm:"ForeignKey:UserID"`
	HotelID       uint      `gorm:"not null"`
	Hotels        Hotels    `gorm:"ForeignKey:HotelID"`
	RoomID        uint      `gorm:"not null"`
	Rooms         Rooms     `gorm:"ForeignKey:RoomID"`
	OwnerID       uint      `gorm:"not null"`
	Owner         Owner     `gorm:"ForeignKey:OwnerID"`
	RoomNo        uint      `gorm:"not null"`
	CheckInDate   time.Time `gorm:"not null"`
	CheckOutDate  time.Time `gorm:"not null"`
	PaymentMethod string
	PaymentAmount float64
	TotalDays     uint
	AdminAmount   float64   `gorm:"not null"`
	OwnerAmount   float64   `gorm:"not null"`
	RoomCategory  string    `gorm:"not null"`
	Review        bool      `gorm:"default:false"`
	Report        bool      `gorm:"default:false"`
	Refund        bool      `gorm:"default:false"`
	Cancel        bool      `gorm:"default:false"`
	AdminResponse string    `gorm:"default:'No Action'"`
	BookedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// func (Booking) TableName() string {
// 	return "bookings" // Adjust the table name to match your database schema
// }
