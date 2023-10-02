package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

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

type JSONB []interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
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
	Fecility        JSONB         `gorm:"type:jsonb" json:"fecilities"`
	Revenue         float64       `json:"revenue" gorm:"default=0"`
	IsAvailable     bool          `json:"isAvailable" gorm:"default=false"`
	IsBlock         bool          `json:"isBlock"`
	AdminApproval   bool          `json:"adminApproval" gorm:"default=false"`
	HotelCategoryID uint          `json:"category_id" gorm:"not null"`
	HotelCategory   HotelCategory `gorm:"ForeignKey:HotelCategoryID"`
	OwnerUsername   string
}
