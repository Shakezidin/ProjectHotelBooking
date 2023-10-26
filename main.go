package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/routes"
)

func main() {

	Init.LoadEnvironmentVariables()
	Init.DatabaseConnection()

	r := gin.Default()
	r.Use(cors.Default())
	r.LoadHTMLGlob("templates/*")
	routes.OwnerRoutes(r)
	routes.UserRoutes(r)
	routes.AdminRoutes(r)

	r.Run("localhost:3000")
}
