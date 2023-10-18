package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

//UserAuthMiddleware User verification
func UserAuthMiddleware(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	if header == "" {
		c.JSON(400, gin.H{"error": "token missing"})
		c.Abort()
		return
	}
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
