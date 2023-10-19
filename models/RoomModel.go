package models

import (
	"time"

	"gorm.io/gorm"
)

// Cancellation Model
type Cancellation struct {
	gorm.Model
	CancellationPolicy     string `json:"cancellation_policy" gorm:"not null"`
	RefundAmountPercentage int    `json:"refund_amount_percentage" gorm:"not null"`
}

// AvailableRoom Model
type AvailableRoom struct {
	gorm.Model
	RoomID    uint `json:"room_id"`
	BookingID uint
	CheckIn   time.Time `json:"check_in" time_format:"2006-01-02"`
	CheckOut  time.Time `json:"check_out" time_format:"2006-01-02"`
}

// RoomFacilities Model
type RoomFacilities struct {
	gorm.Model
	RoomAmenities string `json:"room_amenities" gorm:"not null" validate:"required"`
}

// RoomCategory Model
type RoomCategory struct {
	gorm.Model
	Name string `json:"name" gorm:"not null"`
}

// Rooms Model
type Rooms struct {
	gorm.Model
	Description    string       `json:"description" validate:"required" gorm:"not null"`
	Price          float64      `json:"price" gorm:"not null" validate:"required"`
	Adults         int          `json:"adults" gorm:"not null" validate:"required"`
	Children       int          `json:"children" gorm:"not null" validate:"required"`
	Bed            string       `json:"bed" gorm:"not null" validate:"required"`
	Images         string       `json:"images" validate:"required"`
	CancellationID uint         `json:"cancellation_id" gorm:"not null"`
	Cancellation   Cancellation `gorm:"ForeignKey:CancellationID"`
	Facility       JSONB        `gorm:"type:jsonb" json:"facilities"`
	RoomNo         int
	IsAvailable    bool         `json:"is_available" validate:"required"`
	IsBlocked      bool         `json:"is_blocked"`
	DiscountPrice  float64      `json:"discount_price"`
	Discount       float64      `json:"discount"`
	AdminApproval  bool         `json:"admin_approval" gorm:"default=false"`
	HotelsID       uint         `json:"hotel_id" gorm:"not null"`
	Hotels         Hotels       `gorm:"ForeignKey:HotelsID"`
	OwnerUsername  string       `json:"owner_username"`
	RoomCategoryID uint         `json:"category_id" gorm:"not null"`
    RoomCategory   RoomCategory `json:"category" gorm:"ForeignKey:RoomCategoryID"`
}
