package models

import (
	"gorm.io/gorm"
)

type HotelCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}

type HotelAmenities struct {
	FecilityId     uint   `json:"fecilityid" gorm:"primaryKey;autoIncrement"`
	HotelAmenities string `json:"amenities" gorm:"not null"`
}

type Hotels struct {
	gorm.Model
	Name            string        `json:"name" validate:"required"`
	Title           string        `json:"title" validate:"required"`
	Description     string        `json:"description" validate:"required"`
	StartingPrice   float64       `json:"startingprice" validate:"required"`
	City            string        `json:"city" validate:"required"`
	Pincode         string        `json:"pincode" validate:"required"`
	Address         string        `json:"address" validate:"required"`
	Images          string        `json:"images" validate:"required"`
	TypesOfRoom     int           `json:"typesofroom" validate:"required"`
	Fecilities      []string      `json:"facilities" gorm:"type:jsonb"`
	Revenue         float64       `json:"revenue" gorm:"default=0"`
	IsAvailable     bool          `json:"isAvailable" gorm:"default=false"`
	IsBlock         bool          `json:"isBlock"`
	AdminApproval   bool          `json:"adminApproval" gorm:"default=false"`
	HotelCategoryId uint          `json:"category_id" gorm:"not null"`
	HotelCategory   HotelCategory `gorm:"ForeignKey:HotelCategoryId"`
	OwnerUsername   string
}
