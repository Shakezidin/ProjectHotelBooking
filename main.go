package main

import (
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/routes"
)

var R = gin.Default()

func main() {
	Init.Database_connection()
	Init.Getenv()

	r := gin.Default()

	routes.OwnerRoutes(r)
	routes.UserRoutes(r)
	routes.AdminRoutes(r)

	r.Run("localhost:3000")
}
