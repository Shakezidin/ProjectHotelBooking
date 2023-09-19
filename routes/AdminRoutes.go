package routes

import (
	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	Admin "github.com/shaikhzidhin/controllers/Admin"
)

func AdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")
	{
		admin.POST("/login", Admin.AdminLogin)
		admin.GET("/dashbord", Admin.Dashboard)

		//<<<<<<<<<other admin controlls>>>>>>>>>>>>>>>>>>

		admin.GET("/viewcancellation", Auth.AdminAuthMiddleWare, Admin.ViewRoomCancellation)
		admin.POST("/addcancellation", Auth.AdminAuthMiddleWare, Admin.Addcancellation)
		admin.DELETE("/deletecancellation", Auth.AdminAuthMiddleWare, Admin.Deletecancellation)

		admin.GET("/hotelcatagories", Auth.AdminAuthMiddleWare, Admin.ViewHotelCatagories)
		admin.POST("/addhotelcatagories", Auth.AdminAuthMiddleWare, Admin.AddHotlecatagory)
		admin.DELETE("/deletehotelcatagory", Auth.AdminAuthMiddleWare, Admin.DeleteHotelcatagory)

		admin.GET("/viewhotelfecilities", Auth.AdminAuthMiddleWare, Admin.ViewHotelFecilities)
		admin.POST("/addhotelfecility", Auth.AdminAuthMiddleWare, Admin.AddHotelfecilility)
		admin.DELETE("/deletehotelfecility", Auth.AdminAuthMiddleWare, Admin.DeleteHotelfecility)

		admin.GET("/viewowners", Auth.AdminAuthMiddleWare, Admin.ViewOwner)
		admin.GET("/blockowner", Auth.AdminAuthMiddleWare, Admin.BlockOwner)

		admin.GET("/viewblockedhotels", Auth.AdminAuthMiddleWare, Admin.BlockedHotels)
		admin.GET("/ownershotels", Auth.AdminAuthMiddleWare, Admin.OwnerHotels)
		admin.GET("/hotelblockandunblock", Auth.AdminAuthMiddleWare, Admin.BlockandUnblockhotel)
		admin.GET("/approvalpending", Auth.AdminAuthMiddleWare, Admin.HotelforApproval)
		admin.GET("/hotelapproving", Auth.AdminAuthMiddleWare, Admin.HotelsApproval)

		admin.GET("/viewusers", Auth.AdminAuthMiddleWare, Admin.ViewUser)
		admin.GET("/viewblockedusers", Auth.AdminAuthMiddleWare, Admin.ViewBlockedUser)
		admin.GET("/viewunblockedusers", Auth.AdminAuthMiddleWare, Admin.ViewUnblockedUsers)
		admin.GET("/blockandunblockuser", Auth.AdminAuthMiddleWare, Admin.BlockandUnblockUser)
	}
}
