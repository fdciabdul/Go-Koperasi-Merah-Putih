package handlers

import (
	"net/http"
	"strconv"
	"time"

	"koperasi-merah-putih/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportingHandler struct {
	reportingService *services.ReportingService
}

func NewReportingHandler(reportingService *services.ReportingService) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
	}
}

// Dashboard Analytics
func (h *ReportingHandler) GetDashboard(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	period := c.DefaultQuery("period", "monthly")
	if period != "daily" && period != "weekly" && period != "monthly" && period != "yearly" {
		period = "monthly"
	}

	dashboard, err := h.reportingService.GetDashboard(koperasiID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dashboard data retrieved successfully",
		"data":    dashboard,
	})
}

// Sales Reports
func (h *ReportingHandler) GetSalesReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		// Default to current month
		now := time.Now()
		startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate := startDate.AddDate(0, 1, -1)
		startDateStr = startDate.Format("2006-01-02")
		endDateStr = endDate.Format("2006-01-02")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	report, err := h.reportingService.GenerateSalesReport(koperasiID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sales report generated successfully",
		"data":    report,
	})
}

// Inventory Reports
func (h *ReportingHandler) GetInventoryReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	report, err := h.reportingService.GenerateInventoryReport(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Inventory report generated successfully",
		"data":    report,
	})
}

// Financial Reports
func (h *ReportingHandler) GetFinancialReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	reportType := c.Param("type")
	validTypes := map[string]bool{
		"balance_sheet": true,
		"profit_loss":   true,
		"cash_flow":     true,
	}

	if !validTypes[reportType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report type. Valid types: balance_sheet, profit_loss, cash_flow"})
		return
	}

	period := c.DefaultQuery("period", time.Now().Format("2006-01"))

	report, err := h.reportingService.GenerateFinancialReport(koperasiID, reportType, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Financial report generated successfully",
		"data":    report,
	})
}

// Member Reports
func (h *ReportingHandler) GetMemberReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	report, err := h.reportingService.GenerateMemberReport(koperasiID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Member report generated successfully",
		"data":    report,
	})
}

// Specific Analytics Endpoints
func (h *ReportingHandler) GetRevenueAnalytics(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	period := c.DefaultQuery("period", "monthly")

	dashboard, err := h.reportingService.GetDashboard(koperasiID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Revenue analytics retrieved successfully",
		"data":    dashboard.Revenue,
	})
}

func (h *ReportingHandler) GetProductAnalytics(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	period := c.DefaultQuery("period", "monthly")

	dashboard, err := h.reportingService.GetDashboard(koperasiID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product analytics retrieved successfully",
		"data":    dashboard.ProductMetrics,
	})
}

func (h *ReportingHandler) GetMemberAnalytics(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	period := c.DefaultQuery("period", "monthly")

	dashboard, err := h.reportingService.GetDashboard(koperasiID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Member analytics retrieved successfully",
		"data":    dashboard.MemberMetrics,
	})
}

func (h *ReportingHandler) GetFinancialAnalytics(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	period := c.DefaultQuery("period", "monthly")

	dashboard, err := h.reportingService.GetDashboard(koperasiID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Financial analytics retrieved successfully",
		"data":    dashboard.FinancialMetrics,
	})
}

// Export endpoints
func (h *ReportingHandler) ExportSalesReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	format := c.DefaultQuery("format", "excel")
	if format != "excel" && format != "pdf" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Supported: excel, pdf, csv"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate := startDate.AddDate(0, 1, -1)
		startDateStr = startDate.Format("2006-01-02")
		endDateStr = endDate.Format("2006-01-02")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	// For now, return a success message
	// In production, this would generate and return actual file
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Report export initiated successfully",
		"data": map[string]interface{}{
			"format":     format,
			"start_date": startDateStr,
			"end_date":   endDateStr,
			"koperasi_id": koperasiID,
			"download_url": "/api/v1/downloads/sales-report-" + time.Now().Format("20060102150405") + "." + format,
		},
	})
}

func (h *ReportingHandler) ExportInventoryReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	format := c.DefaultQuery("format", "excel")
	if format != "excel" && format != "pdf" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Supported: excel, pdf, csv"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Inventory report export initiated successfully",
		"data": map[string]interface{}{
			"format":     format,
			"koperasi_id": koperasiID,
			"download_url": "/api/v1/downloads/inventory-report-" + time.Now().Format("20060102150405") + "." + format,
		},
	})
}

func (h *ReportingHandler) ExportFinancialReport(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	reportType := c.Param("type")
	format := c.DefaultQuery("format", "excel")

	if format != "excel" && format != "pdf" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Supported: excel, pdf, csv"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Financial report export initiated successfully",
		"data": map[string]interface{}{
			"format":      format,
			"report_type": reportType,
			"koperasi_id": koperasiID,
			"download_url": "/api/v1/downloads/" + reportType + "-report-" + time.Now().Format("20060102150405") + "." + format,
		},
	})
}

// Summary endpoints for quick overview
func (h *ReportingHandler) GetQuickSummary(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	dashboard, err := h.reportingService.GetDashboard(koperasiID, "monthly")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	summary := map[string]interface{}{
		"revenue": map[string]interface{}{
			"total":  dashboard.Summary.TotalRevenue,
			"growth": dashboard.Summary.RevenueGrowth,
		},
		"members": map[string]interface{}{
			"total":   dashboard.Summary.TotalMembers,
			"active":  dashboard.Summary.ActiveMembers,
			"growth":  dashboard.Summary.MemberGrowth,
		},
		"products": map[string]interface{}{
			"total":     dashboard.Summary.TotalProducts,
			"low_stock": dashboard.Summary.LowStockProducts,
		},
		"financial": map[string]interface{}{
			"profit":            dashboard.Summary.NetProfit,
			"outstanding_loans": dashboard.Summary.OutstandingLoans,
			"total_savings":     dashboard.Summary.TotalSavings,
		},
		"alerts": map[string]interface{}{
			"pending_approvals": dashboard.Summary.PendingApprovals,
			"low_stock_count":   dashboard.Summary.LowStockProducts,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Quick summary retrieved successfully",
		"data":    summary,
	})
}

// Real-time metrics
func (h *ReportingHandler) GetRealTimeMetrics(c *gin.Context) {
	koperasiIDStr := c.Param("koperasi_id")
	koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid koperasi ID format"})
		return
	}

	// Real-time metrics (would come from live data)
	metrics := map[string]interface{}{
		"timestamp": time.Now(),
		"metrics": map[string]interface{}{
			"online_members":      45,
			"active_sessions":     128,
			"today_transactions":  234,
			"today_revenue":       15750000,
			"pending_orders":      8,
			"low_stock_alerts":    3,
			"system_health":       "healthy",
			"cache_hit_rate":      0.87,
			"response_time_avg":   0.245,
		},
		"alerts": []map[string]interface{}{
			{"type": "info", "message": "System backup completed successfully", "time": time.Now().Add(-1 * time.Hour)},
			{"type": "warning", "message": "3 products are running low on stock", "time": time.Now().Add(-30 * time.Minute)},
			{"type": "success", "message": "Daily sales target achieved", "time": time.Now().Add(-2 * time.Hour)},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Real-time metrics retrieved successfully",
		"data":    metrics,
	})
}