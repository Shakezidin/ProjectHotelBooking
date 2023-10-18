package routes

import (
	"github.com/gin-gonic/gin"
	UserCtrl "github.com/shaikhzidhin/controllers/User"
	"github.com/shaikhzidhin/controllers/hotelowner"
	"github.com/shaikhzidhin/middleware"
)

// UserRoutes Set up the routes for the user section of the application.
func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	// User Authentication routes
	user.POST("/login", UserCtrl.Login)
	user.POST("/signup", UserCtrl.Signup)
	user.POST("/signup/verification", UserCtrl.SignupVerification)

	// Password Recovery routes
	user.POST("/password/forget", UserCtrl.ForgetPassword)
	user.POST("/password/forget/verifyotp", UserCtrl.VerifyOTP)
	user.POST("/password/set/new", UserCtrl.NewPassword)

	// User Home & Profile routes
	user.GET("/", UserCtrl.Home)
	user.GET("/profile", middleware.UserAuthMiddleware, UserCtrl.Profile)
	user.PATCH("/profile/edit", middleware.UserAuthMiddleware, UserCtrl.ProfileEdit)
	user.PUT("/profile/password/change", middleware.UserAuthMiddleware, UserCtrl.PasswordChange)
	user.GET("/booking/history", middleware.UserAuthMiddleware, UserCtrl.History)

	// Hotels routes
	user.GET("/home", UserCtrl.Home)
	user.GET("/home/banner", UserCtrl.BannerShowing)
	user.POST("/home/banner/hotel", hotelowner.ViewSpecificHotel)
	user.POST("/home/search", UserCtrl.Searching)
	user.POST("/home/search/hotel", UserCtrl.SearchHotelByName)

	// Room routes
	user.GET("/home/rooms", UserCtrl.RoomsView)
	user.GET("/home/rooms/room", UserCtrl.RoomDetails)
	user.POST("/home/rooms/filter", UserCtrl.RoomFilter)

	// Contact routes
	user.POST("/home/contact", middleware.UserAuthMiddleware, UserCtrl.SubmitContact)

	// Booking Management routes
	user.GET("/home/room/book", middleware.UserAuthMiddleware, UserCtrl.CalculateAmountForDays)
	user.GET("/coupons/view", middleware.UserAuthMiddleware, UserCtrl.ViewNonBlockedCoupons)
	user.GET("/coupon/apply", middleware.UserAuthMiddleware, UserCtrl.ApplyCoupon)
	user.GET("/wallet", middleware.UserAuthMiddleware, UserCtrl.ViewWallet)
	user.GET("/wallet/apply", middleware.UserAuthMiddleware, UserCtrl.ApplyWallet)
	user.GET("/payat/hotel", middleware.UserAuthMiddleware, UserCtrl.OfflinePayment)

	// Razorpay routes
	user.GET("/online/payment", UserCtrl.RazorpayPaymentGateway)
	user.GET("/payment/success", UserCtrl.RazorpaySuccess)
	user.GET("/success", UserCtrl.SuccessPage)
	user.GET("/cancel/booking", UserCtrl.CancelBooking)
}
