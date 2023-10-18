package models

import "gorm.io/gorm"

// Contact Model
type Contact struct {
	gorm.Model
	Message string `json:"message"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"ForeignKey:UserID"`
}
