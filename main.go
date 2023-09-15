package main

import (
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/routes"
)

var R = gin.Default()

func main() {
	Init.Database_connection()

	r := gin.Default()

	routes.OwnerRoutes(r)

	r.Run("localhost:3000")
}
