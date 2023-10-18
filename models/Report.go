package models

import "gorm.io/gorm"

//Report Model
type Report struct {
	gorm.Model
	Report        string  `json:"report" gorm:"not null"`
	BookingID     uint    `gorm:"not null"`
	Booking       Booking `gorm:"ForeignKey:BookingID"`
	UserID        uint    `gorm:"not null"`
	User          User    `gorm:"ForeignKey:UserID"`
	AdminResponse string  `gorm:"defaul:no Responce"`
}
