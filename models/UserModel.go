package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	User_Id      uint      `json:"user_id" gorm:"primaryKey;autoIncrement"`
	UserName     string    `json:"username" gorm:"not null;unique" validate:"required"`
	Name         string    `json:"name" gorm:"not null;unique" validate:"required"`
	Email        string    `json:"email" gorm:"not null;unique" validate:"required"`
	Phone        string    `json:"phone" gorm:"not null;unique" validate:"required"`
	Password     string    `json:"password" gorm:"not null" validate:"required"`
	Is_Block     bool      `json:"is_block" gorm:"default:false"`
	Validation   bool      `json:"validation" gorm:"default:false"`
	Wallet       Wallet    `json:"wallet"`
	ReferralCode string    `json:"referral_code"`
	JoinedAt     time.Time `json:"joined_at" gorm:"default:now()"`
}

type Wallet struct {
	gorm.Model
	Balance float64
	User_Id uint `gorm:"unique"`
}

type Transaction struct {
	gorm.Model
	Date    time.Time
	Details string
	Amount  float64
	USer_Id uint
}
