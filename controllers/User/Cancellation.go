package user

import (
	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
	"gorm.io/gorm"
)

func CancelBooking(c *gin.Context) {
	bookingID := c.Query("id")
	var booking models.Booking

	if err := Init.DB.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		c.JSON(400, gin.H{"error": "No bookings found"})
		return
	}

	if booking.Cancelled == "Cancelled" {
		c.JSON(400, gin.H{"Error": "This order has already been canceled"})
		return
	}

	tx := Init.DB.Begin() // Start a database transaction

	ownerID := booking.OwnerID
	var owner models.Owner
	if err := tx.Where("id = ?", ownerID).First(&owner).Error; err != nil {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"Error": "Error while fetching owner"})
		return
	}

	cancellationID := booking.CancellationID
	var cancellation models.Cancellation
	if err := tx.Where("id = ?", cancellationID).First(&cancellation).Error; err != nil {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"error": "Cancellation fetching error"})
		return
	}

	refundPercentage := cancellation.RefundAmountPercentage

	// Calculate the refund amount for the user's wallet based on the cancellation scheme
	refundAmount := (float64(refundPercentage) / 100) * booking.PaymentAmount

	// Calculate admin and owner revenue adjustments
	adminRevenueAdjustment := (refundAmount * 1 / 4) // Admin gets 1/4 of the refund
	ownerRevenueAdjustment := (refundAmount * 3 / 4) // Owner gets 3/4 of the refund

	// Update admin and owner revenues
	if err := tx.Model(&models.Revenue{}).Where("owner_id = ?", ownerID).Update("admin_revenue", gorm.Expr("admin_revenue - ?", adminRevenueAdjustment)).Error; err != nil {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"error": "Admin revenue updating error"})
		return
	}

	if err := tx.Model(&models.Owner{}).Where("id = ?", ownerID).Update("Revenue", gorm.Expr("Revenue - ?", ownerRevenueAdjustment)).Error; err != nil {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"error": "Owner revenue updating error"})
		return
	}

	// Update the booking's cancellation status and save it to the database
	if err := tx.Model(&models.Booking{}).Where("id = ?", bookingID).Update("Cancelled", "Cancelled").Error; err != nil {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"error": "Booking status updating error"})
		return
	}

	// Calculate the refund amount for the user's wallet
	userID := booking.UserID
	userWallet := GetUserWallet(tx, userID)

	if userWallet != nil {
		userWallet.Balance += refundAmount

		// Update the user's wallet balance in the database
		if err := tx.Save(userWallet).Error; err != nil {
			tx.Rollback() // Rollback the transaction
			c.JSON(400, gin.H{"error": "User wallet balance updating error"})
			return
		}
	} else {
		tx.Rollback() // Rollback the transaction
		c.JSON(400, gin.H{"error": "User wallet not found"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(400, gin.H{"error": "Transaction commit error"})
		return
	}

	// Respond with a success message
	c.JSON(200, gin.H{"message": "Booking canceled successfully"})
}

// GetUserWallet helps to get user wallet by user id
func GetUserWallet(tx *gorm.DB, userID uint) *models.Wallet {
	var wallet models.Wallet
	if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil
	}
	return &wallet
}
