package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

// ViewReports returns a list of all reports.
func ViewReports(c *gin.Context) {
	var reports []models.Report
	if err := Init.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reports"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

// DeleteReport deletes a report by ID.
func DeleteReport(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID query missing"})
		return
	}
	var report models.Report
	if err := Init.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error fetching report"})
		return
	}
	if err := Init.DB.Delete(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Report deleted successfully"})
}

// ReportDetails retrieves details of a single report.
func ReportDetails(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID query missing"})
		return
	}
	var report models.Report
	if err := Init.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error fetching report"})
		return
	}

	var booking models.Booking
	if err := Init.DB.First(&booking, report.BookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error fetching booking"})
		return
	}

	var hotel models.Hotels
	if err := Init.DB.First(&hotel, booking.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error fetching hotel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"report": report, "booking": booking, "hotel": hotel})
}

// ReportStatus updates the status of a report and associated booking.
func ReportStatus(c *gin.Context) {
	var input struct {
		BookingID uint   `json:"bookingId"`
		ReportID  uint   `json:"reportId"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var booking models.Booking
	if err := Init.DB.First(&booking, input.BookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	var report models.Report
	if err := Init.DB.First(&report, input.ReportID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	booking.AdminResponse = input.Status
	report.AdminResponse = input.Status
	if err := Init.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}
	if err := Init.DB.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update report"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}
