package models

import "gorm.io/gorm"

type Owner struct {
	gorm.Model
	User_Id  int    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	UserName string `json:"username"  gorm:"not null;unique" validate:"required,min=2,max=50"`
	Email    string `json:"email" gorm:"not null;unique" validate:"required,min=2,max=50"`
	Phone    string `json:"phone" gorm:"not null;unique" validate:"required"`
	Revenue  int    `json:"revenue" gorm:"default=0"`
	Password string `json:"password" gorm:"not noll" validate:"required"`
	Is_Block bool   `json:"is_block" gorm:"default=false"`
}

