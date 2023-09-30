package routes

import (
	"github.com/gin-gonic/gin"
	User "github.com/shaikhzidhin/controllers/User"
	"github.com/shaikhzidhin/middleware"
)

func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")
	{
		user.POST("/signup", User.UserSignup)
		user.POST("/signup-verification", User.SingupVerification)
		user.POST("/login", User.UserLogin)
		user.POST("/forgetpassword",User.ForgetPassword)
		user.POST("/verifyotp",User.VerifyOTP)
		user.POST("/setnewpassword",User.Newpassword)
	
		user.GET("/profile", middleware.UserAuthMiddleWare, User.Profile)
		user.PATCH("/editprofile", middleware.UserAuthMiddleWare, User.ProfileEdit)
		user.PUT("/changepassword", middleware.UserAuthMiddleWare, User.PasswordChange)
		
		user.GET("/homepage", User.UserHome)
		user.GET("/banner",User.BannerShowing)
		user.POST("/searchresult",User.Searching)
		user.GET("/seerooms",User.RoomsView)
		user.GET("/seeroom",User.RoomDetails)

		user.POST("/roomfilter",User.RoomFilter)

		user.POST("/searchhotel",User.SearchHotel)

		user.GET("/book",middleware.UserAuthMiddleWare,User.Book)
		user.GET("/viewcoupens",middleware.UserAuthMiddleWare,User.Coupons)
		user.GET("/applycoupen",middleware.UserAuthMiddleWare,User.ApplyCoupon)
	}
}
