package models

import "gorm.io/gorm"

type Report struct {
	gorm.Model
	Report        string  `json:"report" gorm:"not null"`
	BookingId     uint    `gorm:"not null"`
	Booking       Booking `gorm:"ForeignKey:BookingId"`
	UserId        uint    `gorm:"not null"`
	User          User    `gorm:"ForeignKey:UserId"`
	AdminResponse string  `gorm:"defaul:no Responce"`
}
