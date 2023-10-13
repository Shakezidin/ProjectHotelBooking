package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/controllers/HotelOwner"
	"github.com/shaikhzidhin/middleware"
)

func OwnerRoutes(c *gin.Engine) {

	owner := c.Group("/owner")
	{
		//<<<<<<<<<<<<<<<<<< Owner Signup & login >>>>>>>>>>>>>
		owner.POST("/singup", HotelOwner.OwnerSignUp)
		owner.POST("/signup/verification", HotelOwner.OwnerSingupVerification)
		owner.POST("/login", HotelOwner.OwnerLogin)

		//<<<<<<<<<<<<<<<<<< Owner DashBord >>>>>>>>>>>>>>>>>>
		owner.GET("/home", middleware.AuthMiddleWare, HotelOwner.GetOwnerDashboard)
		owner.GET("/profile", middleware.AuthMiddleWare, HotelOwner.OwnerProfile)
		owner.PATCH("/profile/edit", middleware.AuthMiddleWare, HotelOwner.ProfileEdit)

		hotel := owner.Group("/hotel")
		{
			//<<<<<<<<< MiddleWare >>>>>>>>>>>>
			hotel.Use(middleware.AuthMiddleWare)

			//<<<<<<<<<<<< Hotel Management >>>>>>>>>>>>>>>>
			hotel.POST("/add", HotelOwner.AddHotel)
			hotel.GET("/view/hotel", HotelOwner.ViewSpecificHotel)
			hotel.GET("/view/hotels", HotelOwner.ViewHotels)
			hotel.PATCH("/edit", HotelOwner.Hoteledit)
			hotel.DELETE("/delete", HotelOwner.DeleteHotel)
			hotel.GET("/availability", HotelOwner.HotelAvailability)
			hotel.GET("/fecilities", HotelOwner.ViewHotelFecilities)
		}
		room := owner.Group("/room")
		{
			//<<<<<<<<< MiddleWare >>>>>>>>>>>>
			room.Use(middleware.AuthMiddleWare)

			//<<<<<<<<<<<< Room Management >>>>>>>>>>>>>>>>>
			room.GET("/fecilities", HotelOwner.ViewRoomfecilities)
			room.GET("/cancellation", HotelOwner.ViewCancellation)
			room.GET("/category", HotelOwner.ViewRoomCatagory)

			room.POST("/add", HotelOwner.AddingRoom)
			room.PATCH("/edit", HotelOwner.EditRoom)
			room.GET("/view/rooms", HotelOwner.ViewRooms)
			room.GET("/view/room", HotelOwner.ViewspecificRoom)
			room.GET("/availability", HotelOwner.RoomAvailability)
			room.DELETE("/delete", HotelOwner.DeleteRoom)
		}
		banner := owner.Group("/banner")
		{
			//<<<<<<<<< MiddleWare >>>>>>>>>>>>
			banner.Use(middleware.AuthMiddleWare)

			//<<<<<<<<<Banner Management >>>>>>>>>>>>>>>
			banner.GET("/view", HotelOwner.ViewBanners)
			banner.POST("/add", HotelOwner.AddBanner)
			banner.PATCH("/edit", HotelOwner.UpdateBanner)
			banner.GET("/availability", HotelOwner.AvailableBanner)
			banner.DELETE("/delete", HotelOwner.DeleteBanner)
		}
	}

}
