package models

import "gorm.io/gorm"

//Owner Model
type Owner struct {
	gorm.Model
	Name     string `json:"name"`
	UserName string `json:"username"  gorm:"not null;unique" validate:"required,min=2,max=50"`
	Email    string `json:"email" gorm:"not null;unique" validate:"required,min=2,max=50"`
	Phone    string `json:"phone" gorm:"not null;unique" validate:"required"`
	Revenue  int    `json:"revenue" gorm:"default=0"`
	Password string `json:"password" gorm:"not noll" validate:"required"`
	IsBlocked bool   `json:"is_block" gorm:"default=false"`
}
