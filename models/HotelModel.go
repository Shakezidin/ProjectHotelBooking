package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// HotelCategory model
type HotelCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
}

// HotelAmenities model
type HotelAmenities struct {
	FacilityID     uint   `json:"facility_id" gorm:"primaryKey;autoIncrement"`
	HotelAmenities string `json:"amenities" gorm:"not null"`
}

// JSONB type for handling JSON data in the database
type JSONB []interface{}

//Value used to retrive value
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

//Scan helps to scan values
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// Hotels model
type Hotels struct {
	gorm.Model
	Name            string  `json:"name" validate:"required"`
	Title           string  `json:"title" validate:"required"`
	Description     string  `json:"description" validate:"required"`
	StartingPrice   float64 `json:"starting_price" validate:"required"`
	City            string  `json:"city" validate:"required"`
	Pincode         string  `json:"pincode" validate:"required"`
	Address         string  `json:"address" validate:"required"`
	Images          string  `json:"images" validate:"required"`
	TypesOfRoom     int     `json:"types_of_room" validate:"required"`
	Facility        JSONB   `gorm:"type:jsonb" json:"facilities"`
	Revenue         float64 `json:"revenue" gorm:"default=0"`
	IsAvailable     bool    `json:"is_available" gorm:"default=false"`
	IsBlock         bool    `json:"is_block"`
	AdminApproval   bool    `json:"admin_approval" gorm:"default=false"`
	HotelCategoryID uint    `json:"category_id" gorm:"not null"`
	HotelCategory   HotelCategory `gorm:"ForeignKey:HotelCategoryID"`
	OwnerUsername   string
}
