package modules

import (
	"github.com/gin-gonic/gin"
	"go_koperasi/internal/handlers"
	"go_koperasi/internal/middleware"
)

type MasterDataRoutes struct {
	masterDataHandler *handlers.MasterDataHandler
	rbacMiddleware    *middleware.RBACMiddleware
}

func NewMasterDataRoutes(masterDataHandler *handlers.MasterDataHandler, rbacMiddleware *middleware.RBACMiddleware) *MasterDataRoutes {
	return &MasterDataRoutes{
		masterDataHandler: masterDataHandler,
		rbacMiddleware:    rbacMiddleware,
	}
}

func (r *MasterDataRoutes) SetupRoutes(router *gin.RouterGroup) {
	masterData := router.Group("/master-data")
	masterData.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		// KBLI Management
		masterData.POST("/kbli", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateKBLI)
		masterData.GET("/kbli", r.masterDataHandler.GetKBLIList)
		masterData.GET("/kbli/:id", r.masterDataHandler.GetKBLI)
		masterData.PUT("/kbli/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateKBLI)

		// Jenis Koperasi Management
		masterData.POST("/jenis-koperasi", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateJenisKoperasi)
		masterData.GET("/jenis-koperasi", r.masterDataHandler.GetJenisKoperasiList)
		masterData.GET("/jenis-koperasi/:id", r.masterDataHandler.GetJenisKoperasi)
		masterData.PUT("/jenis-koperasi/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateJenisKoperasi)

		// Bentuk Koperasi Management
		masterData.POST("/bentuk-koperasi", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.CreateBentukKoperasi)
		masterData.GET("/bentuk-koperasi", r.masterDataHandler.GetBentukKoperasiList)
		masterData.GET("/bentuk-koperasi/:id", r.masterDataHandler.GetBentukKoperasi)
		masterData.PUT("/bentuk-koperasi/:id", r.rbacMiddleware.SuperAdminOnly(), r.masterDataHandler.UpdateBentukKoperasi)
	}
}