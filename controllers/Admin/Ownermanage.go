package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewOwners(c *gin.Context) {
	var owners []models.Owner

	if err := Init.DB.Find(&owners).Error; err != nil {
		c.JSON(400, gin.H{"error": "fetching owners error"})
		return
	}

	c.JSON(200, gin.H{"owners": owners})
}

func BlockOwner(c *gin.Context) {
	ownerIDStr := c.DefaultQuery("id", "")
	if ownerIDStr == "" {
		c.JSON(400, gin.H{"error": "hotelid query parameter is missing"})
		return
	}
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "convert error"})
		return
	}
	var owner models.Owner

	if err := Init.DB.First(&owner, uint(ownerID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "owner not found",
		})
		return
	}

	owner.Is_Block = !owner.Is_Block

	if err := Init.DB.Save(&owner).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save owner availability",
		})
		return
	}
	c.Status(200)
}
