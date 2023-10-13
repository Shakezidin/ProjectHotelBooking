package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	User "github.com/shaikhzidhin/controllers/User"
	"github.com/shaikhzidhin/middleware"
)

func UserRoutes(c *gin.Engine) {

	user := c.Group("/user")
	{
		store := cookie.NewStore([]byte("iamsuperkey"))
		user.Use(sessions.Sessions("mysession", store))

		user.POST("/signup", User.UserSignup)
		user.POST("/signup-verification", User.SingupVerification)
		user.POST("/login", User.UserLogin)
		user.POST("/forgetpassword", User.ForgetPassword)
		user.POST("/verifyotp", User.VerifyOTP)
		user.POST("/setnewpassword", User.Newpassword)

		user.GET("/profile", middleware.UserAuthMiddleWare, User.Profile)
		user.PATCH("/editprofile", middleware.UserAuthMiddleWare, User.ProfileEdit)
		user.PUT("/changepassword", middleware.UserAuthMiddleWare, User.PasswordChange)
		user.GET("/history", middleware.UserAuthMiddleWare, User.History)

		user.GET("/homepage", User.UserHome)
		user.GET("/banner", User.BannerShowing)
		user.POST("/searchresult", User.Searching)
		user.GET("/seerooms", User.RoomsView)
		user.GET("/seeroom", User.RoomDetails)

		user.POST("/roomfilter", User.RoomFilter)

		user.POST("/searchhotel", User.SearchHotel)

		user.GET("/book", middleware.UserAuthMiddleWare, User.Book)
		user.GET("/viewcoupens", middleware.UserAuthMiddleWare, User.Coupons)
		user.GET("/applycoupen", middleware.UserAuthMiddleWare, User.ApplyCoupon)
		user.GET("/wallet",middleware.UserAuthMiddleWare,User.Wallet)
		user.GET("/applaywallet",middleware.UserAuthMiddleWare,User.Applaywallet)
		user.GET("/payathotel",middleware.UserAuthMiddleWare,User.OfflinePayment)
		user.GET("/onlinepayment",User.Razorpay)
		user.GET("/payment/success",User.RazorpaySuccess)
		user.GET("/success",User.Success)
	}
}
