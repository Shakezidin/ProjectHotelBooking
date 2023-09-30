package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/routes"
)

var R = gin.Default()

func main() {
	Init.Database_connection()
	Init.Getenv()

	r := gin.Default()

	store := cookie.NewStore([]byte("iamsuperkey"))
	r.Use(sessions.Sessions("mysession", store))

	routes.OwnerRoutes(r)
	routes.UserRoutes(r)
	routes.AdminRoutes(r)

	r.Run("localhost:3000")
}
