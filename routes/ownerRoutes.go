package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/controllers/HotelOwner"
	"github.com/shaikhzidhin/middleware"
)

func OwnerRoutes(c *gin.Engine) {

	owner := c.Group("/owner")
	{
		owner.POST("/singup", HotelOwner.OwnerSignUp)
		owner.POST("/signup-verification", HotelOwner.OwnerSingupVerification)
		owner.POST("/login", HotelOwner.OwnerLogin)

	}
	hotel := c.Group("/hotel")
	{
		hotel.GET("/hotelfecilities", middleware.AuthMiddleWare, HotelOwner.ViewHotelFecilities)
		hotel.POST("/addhotel", middleware.AuthMiddleWare, HotelOwner.AddHotel)
		hotel.GET("/viewhotels", middleware.AuthMiddleWare, HotelOwner.ViewHotels)
		hotel.GET("/viewhotel", middleware.AuthMiddleWare, HotelOwner.ViewSpecificHotel)
		hotel.PATCH("/edit", middleware.AuthMiddleWare, HotelOwner.Hoteledit)
		hotel.DELETE("/delete", middleware.AuthMiddleWare, HotelOwner.DeleteHotel)
		hotel.GET("/hotelavailability", middleware.AuthMiddleWare, HotelOwner.HotelAvailability)
	}
	room := c.Group("/room")
	{
		room.GET("/fecilities", middleware.AuthMiddleWare, HotelOwner.ViewRoomfecilities)
		room.GET("/cancellation", middleware.AuthMiddleWare, HotelOwner.ViewCancellation)
		room.GET("/category", middleware.AuthMiddleWare, HotelOwner.ViewRoomCatagory)
		room.GET("/add", middleware.AuthMiddleWare, HotelOwner.AddRoom)
		room.POST("/adding", middleware.AuthMiddleWare, HotelOwner.AddingRoom)
		room.PATCH("/edit", middleware.AuthMiddleWare, HotelOwner.EditRoom)
		room.GET("/viewrooms", middleware.AuthMiddleWare, HotelOwner.ViewRooms)
		room.GET("/viewspecificroom", middleware.AuthMiddleWare, HotelOwner.ViewspecificRoom)
		room.DELETE("/delete", middleware.AuthMiddleWare, HotelOwner.DeleteRoom)
		room.GET("/availability", middleware.AuthMiddleWare, HotelOwner.RoomAvailability)
	}

}
