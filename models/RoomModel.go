package models

import (
	"time"

	"gorm.io/gorm"
)

type Cancellation struct {
	gorm.Model
	CancellationPolicy string `json:"cancellation_policy" gorm:"not null"`
}

type AvailableRoom struct {
	gorm.Model
	RoomID      uint        `json:"room_id"`
	RoomNo      uint        `json:"room_no"`
	CheckIn     []time.Time `gorm:"type:jsonb" json:"check_in"`
	Checkout    []time.Time `gorm:"type:jsonb" json:"checkout"`
	IsAvailable bool        `json:"is_available"`
}

type RoomFecilities struct {
	gorm.Model
	RoomAmanities string `json:"roomamanities" gorm:"not null" validate:"required"`
}

type RoomCategory struct {
	gorm.Model
	Name            string `json:"name" gorm:"not null"`
}

type Rooms struct {
	gorm.Model
	Description    string          `json:"description" gorm:"not null"`
	Price          float64         `json:"price" gorm:"not null"`
	Adults         uint            `json:"adults" gorm:"not null"`
	Children       uint            `json:"children" gorm:"not null"`
	Bed            string          `json:"bed" gorm:"not null"`
	Images         string          `json:"images" validate:"required"`
	CancellationId uint            `json:"cancellation_id" gorm:"not null"`
	Cancellation   Cancellation    `gorm:"ForeignKey:CancellationId"`
	NoOfRooms      uint            `json:"number_of_rooms" gorm:"default:1"`
	AvailableRooms []AvailableRoom `json:"available_rooms" gorm:"foreignKey:RoomID"`
	Fecilities     []string        `json:"facilities" gorm:"type:jsonb"`
	IsAvailable    bool            `json:"is_available"`
	IsBlocked      bool            `json:"is_blocked"`
	DiscountPrice  float64         `json:"discount_price"`
	Discount       float64         `json:"discount"`
	AdminApproval  bool            `json:"admin_approval" gorm:"default=false"`
	HotelID        uint            `json:"hotel_id"`
	OwnerUsername  string          `json:"owner_username"`
	RoomCategoryId uint            `json:"category_id" gorm:"not null"`
	RoomCategory   RoomCategory    `gorm:"ForeignKey:RoomCategoryId"`
}
