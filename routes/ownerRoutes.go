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
		owner.GET("/home", middleware.AuthMiddleWare, HotelOwner.GetOwnerDashboard)
		owner.GET("/profile", middleware.AuthMiddleWare, HotelOwner.OwnerProfile)
		owner.PATCH("/profileedit", middleware.AuthMiddleWare, HotelOwner.ProfileEdit)
	}

	hotel := c.Group("/hotel")
	{
		hotel.Use(middleware.AuthMiddleWare)

		hotel.GET("/hotelfecilities", HotelOwner.ViewHotelFecilities)
		hotel.POST("/addhotel", HotelOwner.AddHotel)
		hotel.GET("/viewhotels", HotelOwner.ViewHotels)
		hotel.GET("/viewhotel", HotelOwner.ViewSpecificHotel)
		hotel.PATCH("/edit", HotelOwner.Hoteledit)
		hotel.DELETE("/delete", HotelOwner.DeleteHotel)
		hotel.GET("/hotelavailability", HotelOwner.HotelAvailability)
		hotel.GET("/viewbanners", HotelOwner.ViewBanners)
		hotel.GET("/canaddbanner", HotelOwner.CanAddBanner)
		hotel.POST("/addbanner", HotelOwner.AddBanner)
		hotel.GET("/editbanner", HotelOwner.EditBanner)
		hotel.POST("/updatebanner", HotelOwner.UpdateBanner)
		hotel.GET("/banneravailability", HotelOwner.AvailableBanner)
		hotel.DELETE("/deletebanner", HotelOwner.DeleteBanner)
	}
	room := c.Group("/room")
	{
		room.Use(middleware.AuthMiddleWare)

		room.GET("/fecilities", HotelOwner.ViewRoomfecilities)
		room.GET("/cancellation", HotelOwner.ViewCancellation)
		room.GET("/category", HotelOwner.ViewRoomCatagory)
		room.GET("/add", HotelOwner.AddRoom)
		room.POST("/adding", HotelOwner.AddingRoom)
		room.PATCH("/edit", HotelOwner.EditRoom)
		room.GET("/viewrooms", HotelOwner.ViewRooms)
		room.GET("/viewspecificroom", HotelOwner.ViewspecificRoom)
		room.DELETE("/delete", HotelOwner.DeleteRoom)
		room.GET("/availability", HotelOwner.RoomAvailability)
	}

}
