package routes

import (
	"github.com/gin-gonic/gin"
	AdminCtrl "github.com/shaikhzidhin/controllers/Admin"
	HotelOwnerCtrl "github.com/shaikhzidhin/controllers/hotelowner"
	"github.com/shaikhzidhin/middleware"
)

// AdminRoutes Set up the routes for the admin section of the application.
func AdminRoutes(c *gin.Engine) {

	// Admin Login and Dashboard routes
	c.POST("/admin/login", AdminCtrl.Login)
	c.GET("/admin/dashboard", AdminCtrl.Dashboard)

	admin := c.Group("/admin")
	admin.Use(middleware.AdminAuthMiddleWare)

	// Owner Management routes
	admin.GET("/owners", AdminCtrl.ViewOwners)
	admin.GET("/owner/details", AdminCtrl.ViewOwner)
	admin.GET("/owner/block", AdminCtrl.BlockOwner)

	// Hotel Management routes
	admin.GET("/owner/hotels", AdminCtrl.OwnerHotels)
	admin.GET("/hotel/block", AdminCtrl.BlockAndUnblockHotel)
	admin.GET("/hotels/blocked", AdminCtrl.BlockedHotels)
	admin.GET("/hotel/pending", AdminCtrl.HotelsForApproval)
	admin.GET("/hotel/approving", AdminCtrl.HotelsApproval)
	admin.GET("/hotel/viewdetails", HotelOwnerCtrl.ViewSpecificHotel)

	// Hotel Facilities controls
	admin.GET("/hotel/facilities/view", AdminCtrl.ViewHotelFacilities)
	admin.POST("/hotel/facility/add", AdminCtrl.AddHotelFacility)
	admin.DELETE("/hotel/facility/delete", AdminCtrl.DeleteHotelFacility)

	// Hotel Categories controls
	admin.GET("/hotel/categories/view", AdminCtrl.ViewHotelCategories)
	admin.POST("/hotel/categories/add", AdminCtrl.AddHotelCategory)
	admin.DELETE("/hotel/category/delete", AdminCtrl.DeleteHotelCategory)

	// Room Management routes
	admin.GET("/owner/hotel/rooms", AdminCtrl.ViewRooms)
	admin.GET("/hotel/blocked/rooms", AdminCtrl.BlockedRooms)
	admin.GET("/hotel/room/block", AdminCtrl.BlockAndUnblockRooms)
	admin.GET("/hotel/room/pending", AdminCtrl.RoomsForApproval)
	admin.GET("/hotel/room/approving", AdminCtrl.RoomsApproval)
	admin.GET("/hotel/room/details", HotelOwnerCtrl.ViewSpecificRoom)

	// Room Facilities controls
	admin.GET("/room/facilities/view", AdminCtrl.ViewRoomFacilities)
	admin.POST("/room/facility/add", AdminCtrl.AddRoomFacility)
	admin.DELETE("/room/facility/delete", AdminCtrl.DeleteRoomFacility)

	// Room Cancellation controls
	admin.GET("/room/cancellation/view", AdminCtrl.ViewRoomCancellation)
	admin.POST("/room/cancellation/add", AdminCtrl.AddCancellation)
	admin.DELETE("/room/cancellation/delete", AdminCtrl.DeleteCancellation)

	// Room Category controls
	admin.GET("/room/categories/view", AdminCtrl.ViewRoomCategories)
	admin.POST("/room/category/add", AdminCtrl.AddRoomCategory)
	admin.DELETE("/room/category/delete", AdminCtrl.DeleteRoomCategory)

	// Message Management routes
	admin.GET("/messages", AdminCtrl.GetMessages)
	admin.DELETE("/message/delete", AdminCtrl.DeleteMessage)

	// Report Management routes
	admin.GET("/reports", AdminCtrl.ViewReports)
	admin.GET("/report/details", AdminCtrl.ReportDetails)
	admin.POST("/report/status", AdminCtrl.ReportStatus)
	admin.DELETE("/report/delete", AdminCtrl.DeleteReport)

	// Banner Management routes
	admin.GET("/banners", AdminCtrl.BannerView)
	admin.GET("/banner/details", AdminCtrl.BannerDetails)
	admin.GET("/banner/activate", AdminCtrl.BannerSetActive)
	admin.DELETE("/banner/delete", AdminCtrl.DeleteBanner)

	// User Management routes
	admin.GET("/users", AdminCtrl.ViewUsers)
	admin.GET("/users/blocked", AdminCtrl.ViewBlockedUsers)
	admin.GET("/users/unblocked", AdminCtrl.ViewUnblockedUsers)
	admin.GET("/user/block", AdminCtrl.BlockAndUnblockUser)

	// Coupon Management routes
	admin.GET("/coupons/view", AdminCtrl.ViewCoupons)
	admin.POST("/coupon/add", AdminCtrl.CreateCoupon)
	admin.GET("/coupon/block", AdminCtrl.BlockCoupon)
	admin.GET("/coupon/getcoupon", AdminCtrl.GetCoupon)
	admin.PATCH("/coupon/update", AdminCtrl.UpdateCoupon)
	admin.DELETE("/coupon/delete", AdminCtrl.DeleteCoupon)

}
