package models

import (
	"gorm.io/gorm"
)

// Banner Model
type Banner struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null"`
	Subtitle  string `json:"subtitle" gorm:"not null"`
	ImageURL  string `json:"image_url" gorm:"not null"`
	LinkTo    string `json:"link_to" gorm:"not null"`
	Available bool   `json:"available" gorm:"default:true"`
	Active    bool   `json:"active" gorm:"default:false"`
	OwnerID   uint   `json:"owner_id" `
	Owner     Owner  `gorm:"ForeignKey:OwnerID"`
	HotelsID  uint   `json:"hotel_id"`
	Hotels    Hotels `gorm:"ForeignKey:HotelsID"`
}
