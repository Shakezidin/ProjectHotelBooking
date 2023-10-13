package models

import "gorm.io/gorm"

type Contact struct{
	gorm.Model
	Message string `json:"message"`
	User_Id uint `gorm:"not null"`
	User User `gorm:"ForeignKey:User_Id"`
}