package models

import (
	"time"
)

type User struct {
	User_Id      uint64    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	UserName     string    `json:"username" gorm:"not null;unique" validate:"required"`
	Name         string    `json:"name" gorm:"not null;unique" validate:"required"`
	Email        string    `json:"email" gorm:"not null;unique"`
	Phone        string    `json:"phone" gorm:"not null;unique"`
	Password     string    `json:"password" gorm:"not null"`
	Is_Block     bool      `json:"is_block" gorm:"default:false"`
	Validation   bool      `json:"validation" gorm:"default:false"`
	Wallet       int       `json:"wallet" gorm:"default=0"`
	ReferralCode string    `json:"referral_code" gorm:"unique"`
	JoinedAt     time.Time `json:"joined_at" gorm:"default:now()"`
}
