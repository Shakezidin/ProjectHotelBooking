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
	DB.AutoMigrate(&models.Cancellation{})
	DB.AutoMigrate(&models.RoomFecilities{})
	DB.AutoMigrate(&models.HotelCategory{})
	DB.AutoMigrate(&models.Owner{})
	DB.AutoMigrate(&models.HotelAmenities{})
	DB.AutoMigrate(&models.Hotel{})
	DB.AutoMigrate(&models.Room{})

}
