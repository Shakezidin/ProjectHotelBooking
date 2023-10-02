package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/shaikhzidhin/Auth"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func AdminAuthMiddleWare(c *gin.Context) {
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
	var admin models.Admin
	result := Init.DB.Where("user_name = ?", rslt).First(&admin)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin not found"})
		c.Abort()
		return
	}

	c.Next()
}
