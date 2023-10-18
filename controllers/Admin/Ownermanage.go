package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewOwners returns a list of all owners.
func ViewOwners(c *gin.Context) {
	var owners []models.Owner

	if err := Init.DB.Find(&owners).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error fetching owners"})
		return
	}

	c.JSON(200, gin.H{"owners": owners})
}

//ViewOwner returns a single owner
func ViewOwner(c *gin.Context) {
	ownerIDStr := c.DefaultQuery("id", "")
	if ownerIDStr == "" {
		c.JSON(400, gin.H{"error": "Owner ID query parameter is missing"})
		return
	}

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Conversion error"})
		return
	}

	var owner models.Owner
	if err := Init.DB.First(&owner, uint(ownerID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Owner not found",
		})
		return
	}
	c.JSON(200, gin.H{"status": owner})
}

// BlockOwner toggles the 'IsBlocked' field of an owner.
func BlockOwner(c *gin.Context) {
	ownerIDStr := c.DefaultQuery("id", "")
	if ownerIDStr == "" {
		c.JSON(400, gin.H{"error": "Owner ID query parameter is missing"})
		return
	}

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Conversion error"})
		return
	}

	var owner models.Owner
	if err := Init.DB.First(&owner, uint(ownerID)).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Owner not found",
		})
		return
	}

	owner.IsBlocked = !owner.IsBlocked

	if err := Init.DB.Save(&owner).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to save owner availability",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "block updated"})
}
