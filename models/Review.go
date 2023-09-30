package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Review  string `json:"review" gorm:"not null"`
	Rating  int    `json:"rating" gorm:"not null"`
	HotelId uint   `gorm:"not null"`
	UserId  uint   `gorm:"not null"`
}
