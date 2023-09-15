package routes

import (
	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	"github.com/shaikhzidhin/controllers"
	"github.com/shaikhzidhin/controllers/HotelOwner"
)

func OwnerRoutes(c *gin.Engine) {

	owner := c.Group("/owner")
	{
		owner.POST("/singup", HotelOwner.OwnerSignUp)
		owner.POST("/login", HotelOwner.OwnerLogin)
		owner.POST("/login/otp", controllers.OtpLog)
		owner.POST("/login/otpvalidate", controllers.CheckOtp)
	}
	hotel := c.Group("/hotel")
	{
		hotel.GET("/hotelfecilities", Auth.AuthMiddleWare, HotelOwner.ViewHotelFecilities)
		hotel.POST("/addhotel", Auth.AuthMiddleWare, HotelOwner.AddHotel)
		hotel.GET("/viewhotels", Auth.AuthMiddleWare, HotelOwner.ViewHotels)
		hotel.GET("/viewhotel", Auth.AuthMiddleWare, HotelOwner.ViewSpecificHotel)
		hotel.PATCH("/edit", Auth.AuthMiddleWare, HotelOwner.Hoteledit)
		hotel.DELETE("/delete", Auth.AuthMiddleWare, HotelOwner.DeleteHotel)
		hotel.GET("/hotelavailability", Auth.AuthMiddleWare, HotelOwner.HotelAvailability)
	}
	room := c.Group("/room")
	{
		room.GET("/fecilities", Auth.AuthMiddleWare, HotelOwner.ViewRoomfecilities)
		room.GET("/cancellation", Auth.AuthMiddleWare, HotelOwner.ViewCancellation)
		room.GET("/category", Auth.AuthMiddleWare, HotelOwner.ViewRoomCatagory)
		room.GET("/add", Auth.AuthMiddleWare, HotelOwner.AddRoom)
		room.POST("/adding", Auth.AuthMiddleWare, HotelOwner.AddingRoom)
		room.PATCH("/edit", Auth.AuthMiddleWare, HotelOwner.EditRoom)
		room.GET("/viewrooms", Auth.AuthMiddleWare, HotelOwner.ViewRooms)
		room.GET("/viewspecificroom", Auth.AuthMiddleWare, HotelOwner.ViewspecificRoom)
		room.DELETE("/delete", Auth.AuthMiddleWare, HotelOwner.DeleteRoom)
		room.GET("/availability", Auth.AuthMiddleWare, HotelOwner.RoomAvailability)
	}

}
