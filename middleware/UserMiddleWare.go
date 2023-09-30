package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func UserAuthMiddleWare(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	rslt, err := auth.Trim(header)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "trim error"})
		c.Abort()
		return
	}
	var user models.User
	result := Init.DB.Where("user_name = ?", rslt).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		c.Abort()
		return
	}

	c.Next()
}
