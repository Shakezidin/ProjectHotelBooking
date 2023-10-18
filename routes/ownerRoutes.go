package routes

import (
	"github.com/gin-gonic/gin"
	Admin "github.com/shaikhzidhin/controllers/Admin"
	"github.com/shaikhzidhin/controllers/hotelowner"
	"github.com/shaikhzidhin/middleware"
)

//OwnerRoutes Set up routes for the owner section of the application.
func OwnerRoutes(c *gin.Engine) {
	owner := c.Group("/owner")

	// Owner Signup & login routes
	owner.POST("/signup", hotelowner.OwnerSignUp)
	owner.POST("/signup/verification", hotelowner.OwnerSignUpVerification)
	owner.POST("/login", hotelowner.OwnerLogin)

	// Owner Dashboard routes
	owner.GET("/home", middleware.AuthMiddleWare, hotelowner.GetOwnerDashboard)
	owner.GET("/profile", middleware.AuthMiddleWare, hotelowner.OwnerProfile)
	owner.PATCH("/profile/edit", middleware.AuthMiddleWare, hotelowner.ProfileEdit)

	hotel := owner.Group("/hotel")
	// hotel.Use(middleware.AuthMiddleWare)

	// Hotel Management routes
	hotel.POST("/add", hotelowner.AddHotel)
	hotel.GET("/view/hotel", hotelowner.ViewSpecificHotel)
	hotel.GET("/view/hotels", hotelowner.ViewHotels)
	hotel.PATCH("/edit", hotelowner.Hoteledit)
	hotel.DELETE("/delete", hotelowner.DeleteHotel)
	hotel.GET("/availability", hotelowner.HotelAvailability)
	hotel.GET("/fecilities", Admin.ViewHotelFacilities)

	room := owner.Group("/room")
	// room.Use(middleware.AuthMiddleWare)

	// Room Management routes
	room.GET("/fecilities", Admin.ViewRoomFacilities)
	room.GET("/cancellation", Admin.ViewRoomCancellation)
	room.GET("/category", Admin.ViewRoomCategories)

	room.POST("/add", hotelowner.AddingRoom)
	room.PATCH("/edit", hotelowner.EditRoom)
	room.GET("/view/rooms", hotelowner.ViewRooms)
	room.GET("/view/room", hotelowner.ViewSpecificRoom)
	room.GET("/availability", hotelowner.RoomAvailability)
	room.DELETE("/delete", hotelowner.DeleteRoom)

	banner := owner.Group("/banner")
	// banner.Use(middleware.AuthMiddleWare)

	// Banner Management routes
	banner.GET("/view", hotelowner.ViewBanners)
	banner.POST("/add", hotelowner.AddBanner)
	banner.PATCH("/edit", hotelowner.UpdateBanner)
	banner.GET("/availability", hotelowner.AvailableBanner)
	banner.DELETE("/delete", hotelowner.DeleteBanner)
}
