package initializer

import (
	"fmt"
	"os"

	"github.com/shaikhzidhin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the GORM database instance.
var DB *gorm.DB

// DatabaseConnection initializes and returns a GORM database connection.
func DatabaseConnection() *gorm.DB {
	host := os.Getenv("HOST")
	user := os.Getenv("PSQLUSER")
	password := os.Getenv("PSQLPASSWORD")
	dbname := os.Getenv("DATABASENAME")
	port := os.Getenv("PORT")
	sslmode := os.Getenv("SSLMODE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Connection to the database failed:", err)
	}

	// AutoMigrate all models
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

	return DB
}
