package main

import (
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/routes"
)

func main() {

	Init.DatabaseConnection()
	Init.LoadEnvironmentVariables()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	routes.OwnerRoutes(r)
	routes.UserRoutes(r)
	routes.AdminRoutes(r)

	r.Run("localhost:3000")
}
