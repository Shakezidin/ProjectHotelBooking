package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.
type User struct {
	gorm.Model
	UserName     string    `json:"username" gorm:"not null;unique" validate:"required"`
	Name         string    `json:"name" gorm:"not null" validate:"required"`
	Email        string    `json:"email" gorm:"not null;unique" validate:"required"`
	Phone        string    `json:"phone" gorm:"not null;unique" validate:"required"`
	Password     string    `json:"password" gorm:"not null" validate:"required"`
	IsBlocked    bool      `json:"is_blocked" gorm:"default:false"`
	Wallet       Wallet    `json:"wallet"`
	ReferralCode string    `json:"referral_code"`
	JoinedAt     time.Time `json:"joined_at" gorm:"default:now()"`
}

// Wallet represents a user's wallet balance.
type Wallet struct {
	gorm.Model
	Balance float64
	UserID  uint `gorm:"unique"`
}

// Transaction represents a financial transaction.
type Transaction struct {
	gorm.Model
	Date    time.Time
	Details string
	Amount  float64
	UserID  uint
}
