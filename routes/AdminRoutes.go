package routes

import (
	"github.com/gin-gonic/gin"
	Admin "github.com/shaikhzidhin/controllers/Admin"
	"github.com/shaikhzidhin/middleware"
)

func AdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")
	{
		admin.POST("/login", Admin.AdminLogin)
		admin.GET("/dashbord", Admin.Dashboard)

		//<<<<<<<<<other admin controlls>>>>>>>>>>>>>>>>>>

		admin.GET("/viewcancellation", middleware.AdminAuthMiddleWare, Admin.ViewRoomCancellation)
		admin.POST("/addcancellation", middleware.AdminAuthMiddleWare, Admin.Addcancellation)
		admin.DELETE("/deletecancellation", middleware.AdminAuthMiddleWare, Admin.Deletecancellation)

		admin.GET("/hotelcatagories", middleware.AdminAuthMiddleWare, Admin.ViewHotelCatagories)
		admin.POST("/addhotelcatagories", middleware.AdminAuthMiddleWare, Admin.AddHotlecatagory)
		admin.DELETE("/deletehotelcatagory", middleware.AdminAuthMiddleWare, Admin.DeleteHotelcatagory)

		admin.GET("/viewhotelfecilities", middleware.AdminAuthMiddleWare, Admin.ViewHotelFecilities)
		admin.POST("/addhotelfecility", middleware.AdminAuthMiddleWare, Admin.AddHotelfecilility)
		admin.DELETE("/deletehotelfecility", middleware.AdminAuthMiddleWare, Admin.DeleteHotelfecility)

		admin.GET("/viewroomfecilities", middleware.AdminAuthMiddleWare, Admin.ViewRoomFecilities)
		admin.POST("/addroomfecility", middleware.AdminAuthMiddleWare, Admin.AddRoomfecilility)
		admin.DELETE("/deleteroomfecility", middleware.AdminAuthMiddleWare, Admin.DeleteRoomFecility)

		admin.GET("/viewroomcatagories", middleware.AdminAuthMiddleWare, Admin.ViewRoomCatagory)
		admin.POST("/addroomcatagory", middleware.AdminAuthMiddleWare, Admin.AddRoomCatagory)
		admin.DELETE("/deleteroomcatagory", middleware.AdminAuthMiddleWare, Admin.DeleteRoomCatagories)

		admin.GET("/viewowners", middleware.AdminAuthMiddleWare, Admin.ViewOwner)
		admin.GET("/blockowner", middleware.AdminAuthMiddleWare, Admin.BlockOwner)

		admin.GET("/viewblockedhotels", middleware.AdminAuthMiddleWare, Admin.BlockedHotels)
		admin.GET("/ownershotels", middleware.AdminAuthMiddleWare, Admin.OwnerHotels)
		admin.GET("/hotelblockandunblock", middleware.AdminAuthMiddleWare, Admin.BlockandUnblockhotel)
		admin.GET("/approvalpending", middleware.AdminAuthMiddleWare, Admin.HotelforApproval)
		admin.GET("/hotelapproving", middleware.AdminAuthMiddleWare, Admin.HotelsApproval)

		admin.GET("/viewusers", middleware.AdminAuthMiddleWare, Admin.ViewUser)
		admin.GET("/viewblockedusers", middleware.AdminAuthMiddleWare, Admin.ViewBlockedUser)
		admin.GET("/viewunblockedusers", middleware.AdminAuthMiddleWare, Admin.ViewUnblockedUsers)
		admin.GET("/blockandunblockuser", middleware.AdminAuthMiddleWare, Admin.BlockandUnblockUser)
	}
}
