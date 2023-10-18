package initializer

import (
	"fmt"

	"github.com/shaikhzidhin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//DB Initialized with database
var DB *gorm.DB

// DatabaseConnection function to migrate models in database
func DatabaseConnection() {
	dsn := "host=localhost user=postgres password=Sinu1090. dbname=icrodebooking port=5432 sslmode=disable"

	// Assign the database connection to the package-level DB variable

	var err error
	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection to the database failed")
	}
	DB.AutoMigrate(
		&models.Contact{},
		&models.Admin{},
		&models.Cancellation{},
		&models.RoomFacilities{},
		&models.RoomCategory{},
		&models.Report{},
		&models.Rooms{},
		&models.HotelCategory{},
		&models.AvailableRoom{},
		&models.Owner{},
		&models.HotelAmenities{},
		&models.Hotels{},
		&models.User{},
		&models.Coupon{},
		&models.UsedCoupon{},
		&models.Booking{},
		&models.Banner{},
		&models.Wallet{},
		&models.Transaction{},
		&models.RazorPay{},
		&models.Revenue{},
	)
}
