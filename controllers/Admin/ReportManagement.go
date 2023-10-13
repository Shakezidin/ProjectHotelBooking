package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	Init "github.com/shaikhzidhin/initiializer"
	"github.com/shaikhzidhin/models"
)

func ViewReports(c *gin.Context) {
	var reports []models.Report
	if err := Init.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching reports"})
		return
	}
	c.HTML(http.StatusOK, "report", gin.H{"report": reports})
}

func DeleteReport(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "report id query missing"})
		return
	}
	var report models.Report
	if err := Init.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error fetching report"})
		return
	}
	if err := Init.DB.Delete(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "report deleted success"})
}

func ReportDetails(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "report id query missing"})
		return
	}
	var report models.Report
	if err := Init.DB.First(&report, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error fetching report"})
		return
	}

	var booking models.Booking
	if err := Init.DB.First(&booking, report.BookingId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error fetching booking"})
		return
	}

	var hotel models.Hotels
	if err := Init.DB.First(&hotel, booking.HotelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error fetching hotel"})
		return
	}

	c.HTML(http.StatusOK, "reportDetails", gin.H{"report": report, "hotel": hotel})
}

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

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
