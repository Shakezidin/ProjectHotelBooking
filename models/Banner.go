package models

import (
	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null"`
	Subtitle  string `json:"subtitle" gorm:"not null"`
	ImageURL  string `json:"image_url" gorm:"not null"`
	LinkTo    string `json:"link_to" gorm:"not null"`
	Available bool   `json:"available" gorm:"default:true"`
	Active    bool   `json:"active" gorm:"default:false"`
	OwnerID   Owner  `gorm:"foreignkey:User_Id"`
	User_Id   uint   `json:"owner_id" `
	HotelID   Hotels  `gorm:"foreignkey:ID"`
	ID        uint   `json:"hotel_id" `
}
