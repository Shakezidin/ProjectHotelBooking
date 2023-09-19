package routes

import (
	"github.com/gin-gonic/gin"
	Auth "github.com/shaikhzidhin/Auth"
	User "github.com/shaikhzidhin/controllers/User"
)

func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")
	{
		user.POST("/signup", User.UserSignup)
		user.POST("/signup-verification", User.SingupVerification)
		user.POST("/login", User.UserLogin)
	}

	profile := c.Group("/userprofile")
	{
		profile.GET("/profile", Auth.UserAuthMiddleWare, User.Profile)
		profile.PATCH("/editprofile",Auth.UserAuthMiddleWare, User.ProfileEdit)
		profile.POST("/changepassword",Auth.UserAuthMiddleWare,User.PasswordChange)
	}

	home:=c.Group("/userhome")
	{
		home.GET("/homepage",Auth.UserAuthMiddleWare,User.UserHome)
		home.GET("/searchresult",Auth.UserAuthMiddleWare,User.Searching)
	}
}
