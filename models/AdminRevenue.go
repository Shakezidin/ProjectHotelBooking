package models

import "gorm.io/gorm"

//Revenue model of admin
type Revenue struct {
	gorm.Model
	AdminRevenue int   `gorm:"not null"`
	OwnerID      uint  `gorm:"not null"`
	Owner        Owner `gorm:"ForeignKey:OwnerID"`
}
