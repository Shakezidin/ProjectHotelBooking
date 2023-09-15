package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	RoomID           uint          `json:"room_id" gorm:"primaryKey;autoIncrement"`
	Description      string        `json:"description" gorm:"not null"`
	Price            float64       `json:"price" gorm:"not null"`
	Adults           uint          `json:"adults" gorm:"not null"`
	Children         uint          `json:"children" gorm:"not null"`
	Bed              string        `json:"bed" gorm:"not null"`
	Images           string        `json:"images" validate:"required"`
	Cancellation     Cancellation  `json:"cancellation" gorm:"Cancellation_Id"`
	Cancellation_Id  uint          `json:"cancellation_id"`
	NoOfRooms        uint          `json:"number_of_rooms" gorm:"default:1"`
	AvailableRooms   AvailableRoom `json:"available_rooms" gorm:"foreignKey:AvailableRoom_Id"`
	AvailableRoom_Id uint          `json:"availableroom_id"`
	Fecilities       []string      `json:"facilities" gorm:"type:jsonb"`
	IsAvailable      bool          `json:"is_available"`
	IsBlocked        bool          `json:"is_blocked"`
	DiscountPrice    float64       `json:"discount_price"`
	Discount         float64       `json:"discount"`
	AdminApproval    bool          `json:"admin_approval" gorm:"default=false"`
	HotelID          Hotel         `json:"hotel_id" gorm:"foreignkey:ID"`
	ID               uint          `json:"id"`
	OwnerUsername    string        `json:"owner_username"`
	CategoryID       RoomCategory  `json:"category_id" gorm:"foreignKey:Category_Id"`
	Category_Id      uint          `json:"categoryid"`
}

type AvailableRoom struct {
	AvailableRoom_Id uint        `json:"availableroom_id" gorm:"primaryKey;autoIncrement"`
	Room_id          uint        `json:"room_id"`
	RoomNo           uint        `json:"room_no"`
	CheckIn          []time.Time `gorm:"type:array"`
	Checkout         []time.Time `gorm:"type:array"`
	IsAvailable      bool        `json:"is_available"`
}

type Cancellation struct {
	Cancellation_Id     uint   `json:"cancellation_id" gorm:"primaryKey;autoIncrement"`
	Cancellation_Policy string `json:"cancellation_Policy" gorm:"not null"`
}

type RoomFecilities struct {
	RoomAmanities_Id uint   `json:"roomamanities_id" gorm:"primaryKey;autoIncrement"`
	RoomAmanities    string `json:"roomamanities" gorm:"not null" validate:"required"`
}

type RoomCategory struct {
	Category_Id uint   `json:"category_id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"not null"`
}
