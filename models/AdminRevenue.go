package models

import "gorm.io/gorm"

type Revenue struct {
	gorm.Model
	AdminRevenue int   `gorm:"not null"`
	OwnerId      uint  `gorm:"not null"`
	Owner        Owner `gorm:"ForeignKey:OwnerId"`
}
