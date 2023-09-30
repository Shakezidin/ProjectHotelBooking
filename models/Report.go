package models

import "gorm.io/gorm"

type Report struct {
	gorm.Model
	Report        string `json:"report" gorm:"not null"`
	BookingId     uint   `gorm:"not null"`
	UserId        uint   `gorm:"not null"`
	AdminResponse string `gorm:"defaul:no Responce"`
}
