package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type ReportingRoutes struct {
	reportingHandler *handlers.ReportingHandler
	rbacMiddleware   *middleware.RBACMiddleware
}

func NewReportingRoutes(reportingHandler *handlers.ReportingHandler, rbacMiddleware *middleware.RBACMiddleware) *ReportingRoutes {
	return &ReportingRoutes{
		reportingHandler: reportingHandler,
		rbacMiddleware:   rbacMiddleware,
	}
}

func (r *ReportingRoutes) SetupRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	reports.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		// Dashboard & Analytics
		reports.GET("/:koperasi_id/dashboard", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetDashboard)
		reports.GET("/:koperasi_id/quick-summary", r.reportingHandler.GetQuickSummary)
		reports.GET("/:koperasi_id/real-time", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetRealTimeMetrics)

		// Analytics by Domain
		reports.GET("/:koperasi_id/analytics/revenue", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetRevenueAnalytics)
		reports.GET("/:koperasi_id/analytics/products", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetProductAnalytics)
		reports.GET("/:koperasi_id/analytics/members", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetMemberAnalytics)
		reports.GET("/:koperasi_id/analytics/financial", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetFinancialAnalytics)

		// Standard Reports
		reports.GET("/:koperasi_id/sales", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetSalesReport)
		reports.GET("/:koperasi_id/inventory", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetInventoryReport)
		reports.GET("/:koperasi_id/financial/:type", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetFinancialReport)
		reports.GET("/:koperasi_id/members", r.rbacMiddleware.AdminOnly(), r.reportingHandler.GetMemberReport)

		// Export Endpoints
		reports.GET("/:koperasi_id/export/sales", r.rbacMiddleware.AdminOnly(), r.reportingHandler.ExportSalesReport)
		reports.GET("/:koperasi_id/export/inventory", r.rbacMiddleware.AdminOnly(), r.reportingHandler.ExportInventoryReport)
		reports.GET("/:koperasi_id/export/financial/:type", r.rbacMiddleware.AdminOnly(), r.reportingHandler.ExportFinancialReport)
	}
}