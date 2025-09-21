package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type FinancialRoutes struct {
	financialHandler *handlers.FinancialHandler
	rbacMiddleware   *middleware.RBACMiddleware
}

func NewFinancialRoutes(financialHandler *handlers.FinancialHandler, rbacMiddleware *middleware.RBACMiddleware) *FinancialRoutes {
	return &FinancialRoutes{
		financialHandler: financialHandler,
		rbacMiddleware:   rbacMiddleware,
	}
}

func (r *FinancialRoutes) SetupRoutes(router *gin.RouterGroup) {
	financial := router.Group("/financial")
	financial.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess(), r.rbacMiddleware.FinancialAccess())
	{
		// Chart of Accounts
		financial.POST("/coa/akun", r.rbacMiddleware.AdminOnly(), r.financialHandler.CreateCOAAkun)
		financial.GET("/:koperasi_id/coa/akun", r.financialHandler.GetCOAAkunList)
		financial.GET("/coa/kategori", r.financialHandler.GetCOAKategoriList)

		// Journal Management
		financial.POST("/jurnal", r.financialHandler.CreateJurnal)
		financial.GET("/:koperasi_id/jurnal", r.financialHandler.GetJurnalList)
		financial.GET("/jurnal/:id", r.financialHandler.GetJurnal)
		financial.PUT("/jurnal/:id/post", r.financialHandler.PostJurnal)
		financial.PUT("/jurnal/:id/cancel", r.financialHandler.CancelJurnal)

		// Financial Reports
		financial.GET("/:koperasi_id/neraca-saldo", r.financialHandler.GetNeracaSaldo)
		financial.GET("/:koperasi_id/laba-rugi", r.financialHandler.GetLabaRugi)
		financial.GET("/:koperasi_id/neraca", r.financialHandler.GetNeraca)
		financial.GET("/akun/:akun_id/saldo", r.financialHandler.GetSaldoAkun)
	}
}