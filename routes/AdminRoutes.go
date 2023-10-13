package routes

import (
	"github.com/gin-gonic/gin"
	Admin "github.com/shaikhzidhin/controllers/Admin"
	"github.com/shaikhzidhin/controllers/HotelOwner"
	"github.com/shaikhzidhin/middleware"
)

func AdminRoutes(c *gin.Engine) {

	//<<<<<<<<<<< Admin Login 7 dashbord >>>>>>>>>>>>>>>>
	c.POST("/admin/login", Admin.AdminLogin)
	c.GET("/admin/dashbord", Admin.Dashboard)

	admin := c.Group("/admin")
	{
		admin.Use(middleware.AdminAuthMiddleWare)

		//<<<<<<<<<<<<<<< Owner Management >>>>>>>>>>>>>>>>>>>
		admin.GET("/owners", Admin.ViewOwners)
		admin.GET("/owner/details", HotelOwner.OwnerProfile)
		admin.GET("/owner/block", Admin.BlockOwner)

		//<<<<<<<<<<<<<<<hotel Management>>>>>>>>>>>>>>>>>>>>
		admin.GET("/owner/hotels", Admin.OwnerHotels)
		admin.GET("/hotel/block", Admin.BlockandUnblockhotel)
		admin.GET("/hotels/blocked", Admin.BlockedHotels)
		admin.GET("/hotel/pending", Admin.HotelforApproval)
		admin.GET("/hotel/approving", Admin.HotelsApproval)
		admin.GET("/hotel/viewdetails", HotelOwner.ViewSpecificHotel)

		//<<<<<<<<<hotel fecilities controlls>>>>>>>>>>>>>>>>>>
		admin.GET("/hotel/fecilities/view", Admin.ViewHotelFecilities)
		admin.POST("/hotel/fecility/add", Admin.AddHotelfecilility)
		admin.DELETE("/hotel/fecility/delete", Admin.DeleteHotelfecility)

		//<<<<<<<<<hotel catagories controlls>>>>>>>>>>>>>>>>>>
		admin.GET("/hotel/catagories/view", Admin.ViewHotelCatagories)
		admin.POST("/hotel/catagories/add", Admin.AddHotlecatagory)
		admin.DELETE("/hotel/catagory/delete", Admin.DeleteHotelcatagory)

		//<<<<<<<<<<<<<<<Room Management>>>>>>>>>>>>>>>>>>>>
		admin.GET("/owner/hotel/rooms", HotelOwner.ViewRooms)
		admin.GET("/hotel/blocked/rooms", Admin.BlockedRooms)
		admin.GET("/hotel/room/block", Admin.BlockandUnblockRooms)
		admin.GET("/hotel/room/pending", Admin.RoomsforApproval)
		admin.GET("/hotel/room/approving", Admin.RoomsApproval)
		admin.GET("/hotel/room/details", HotelOwner.ViewspecificRoom)

		//<<<<<<<<<room fecilities controlls>>>>>>>>>>>>>>>>>>
		admin.GET("/room/fecilities/view", Admin.ViewRoomFecilities)
		admin.POST("/room/fecility/add", Admin.AddRoomfecilility)
		admin.DELETE("/room/fecility/delete", Admin.DeleteRoomFecility)

		//<<<<<<<<< Room cancellation controlls>>>>>>>>>>>>>>>>>>
		admin.GET("/room/cancellation/view", Admin.ViewRoomCancellation)
		admin.POST("/room/cancellation/add", Admin.Addcancellation)
		admin.DELETE("/room/cancellation/delete", Admin.Deletecancellation)

		//<<<<<<<<<<<room catagory controlls>>>>>>>>>>>>>>>>>>
		admin.GET("/room/catagories/view", Admin.ViewRoomCatagory)
		admin.POST("/room/catagory/add", Admin.AddRoomCatagory)
		admin.DELETE("/room/catagory/delete", Admin.DeleteRoomCatagories)

		//<<<<<<<<<<<<<<<Message Management>>>>>>>>>>>>>>>>>>>>>
		admin.GET("/message", Admin.GetMessages)
		admin.GET("/message/delete", Admin.DeleteMessage)

		//<<<<<<<<<<<<<<<<Report Management >>>>>>>>>>>>>>>>>
		admin.GET("/reports", Admin.ViewReports)
		admin.GET("/report/details", Admin.ReportDetails)
		admin.POST("/report/status", Admin.ReportStatus)
		admin.DELETE("/report/delete", Admin.DeleteReport)

		//<<<<<<<<<<<<< Banner Management >>>>>>>>>>>>>>>>>>>>
		admin.GET("/banners", Admin.BannerView)
		admin.GET("/banner/details", Admin.BannerDetails)
		admin.GET("/banner/activate", Admin.BannerSetActive)
		admin.GET("/banner/delete", Admin.DeleteBanner)

		//<<<<<<<<<<<<<<<<<<User Management>>>>>>>>>>>>>>>>>>>>
		admin.GET("/users", Admin.ViewUser)
		admin.GET("/users/blocked", Admin.ViewBlockedUser)
		admin.GET("/users/unblocked", Admin.ViewUnblockedUsers)
		admin.GET("/user/block", Admin.BlockandUnblockUser)

		//<<<<<<<<<<<<<<<< Coupon Management >>>>>>>>>>>>>>>>>>>
		admin.GET("/coupons/view", Admin.ViewCoupons)
		admin.POST("/coupon/add", Admin.CreateCoupon)
		admin.GET("/coupon/block", Admin.BlockCoupon)
		admin.GET("/coupon/getcoupon", Admin.GetCoupon)
		admin.PATCH("/coupon/update", Admin.UpdateCoupon)
		admin.DELETE("/coupon/delete", Admin.DeleteCoupon)
	}
}
