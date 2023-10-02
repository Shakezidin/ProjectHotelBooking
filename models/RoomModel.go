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
	CheckIn     []time.Time `gorm:"type:timestamp[]" json:"check_in"`
	Checkout    []time.Time `gorm:"type:timestamp[]" json:"checkout"`
	IsAvailable bool        `json:"is_available"`
}

type RoomFecilities struct {
	gorm.Model
	RoomAmanities string `json:"roomamanities" gorm:"not null" validate:"required"`
}

type RoomCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

type Rooms struct {
	gorm.Model
	Description    string       `json:"description" validate:"required" gorm:"not null"`
	Price          float64      `json:"price" gorm:"not null" validate:"required"`
	Adults         int          `json:"adults" gorm:"not null" validate:"required"`
	Children       int          `json:"children" gorm:"not null" validate:"required"`
	Bed            string       `json:"bed" gorm:"not null" validate:"required"`
	Images         string       `json:"images" validate:"required"`
	CancellationId uint         `json:"cancellation_id" gorm:"not null"`
	Cancellation   Cancellation `gorm:"ForeignKey:CancellationId"`
	Fecility       JSONB        `gorm:"type:jsonb" json:"fecilities"`
	RoomNo         int
	IsAvailable    bool         `json:"is_available" validate:"required"`
	IsBlocked      bool         `json:"is_blocked"`
	DiscountPrice  float64      `json:"discount_price"`
	Discount       float64      `json:"discount"`
	AdminApproval  bool         `json:"admin_approval" gorm:"default=false"`
	HotelsId       uint         `json:"hotel_id" gorm:"not null"`
	Hotels         Hotels       `gorm:"ForeignKey:HotelsId" `
	OwnerUsername  string       `json:"owner_username"`
	RoomCategoryId uint         `json:"category_id" gorm:"not null"`
	RoomCategory   RoomCategory `gorm:"ForeignKey:RoomCategoryId"`
}
