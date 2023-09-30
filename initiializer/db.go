package initiializer

import (
	"fmt"

	"github.com/shaikhzidhin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database_connection() {
	dsn := "host=localhost user=postgres password=Sinu1090. dbname=icrodebooking port=5432 sslmode=disable"

	// Assign the database connection to the package-level DB variable
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection to database failed")
	}
	DB.AutoMigrate(
		&models.Admin{},
		&models.Cancellation{},
		&models.RoomFecilities{},
		&models.RoomCategory{},
		&models.Review{},
		&models.Report{},
		&models.Rooms{},
		&models.HotelCategory{},
		&models.AvailableRoom{},
		&models.Owner{},
		&models.HotelAmenities{},
		&models.Hotels{},
		&models.User{},
		&models.Coupon{},
		&models.UsedCoupen{},
		&models.Booking{},
		&models.Banner{},
	)
}
