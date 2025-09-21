package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type KoperasiRoutes struct {
	koperasiHandler *handlers.KoperasiHandler
	rbacMiddleware  *middleware.RBACMiddleware
}

func NewKoperasiRoutes(koperasiHandler *handlers.KoperasiHandler, rbacMiddleware *middleware.RBACMiddleware) *KoperasiRoutes {
	return &KoperasiRoutes{
		koperasiHandler: koperasiHandler,
		rbacMiddleware:  rbacMiddleware,
	}
}

func (r *KoperasiRoutes) SetupRoutes(router *gin.RouterGroup) {
	koperasi := router.Group("/koperasi")
	koperasi.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		// Koperasi CRUD
		koperasi.POST("", r.rbacMiddleware.SuperAdminOnly(), r.koperasiHandler.CreateKoperasi)
		koperasi.GET("", r.koperasiHandler.GetKoperasiList)
		koperasi.GET("/:id", r.koperasiHandler.GetKoperasi)
		koperasi.PUT("/:id", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.UpdateKoperasi)
		koperasi.DELETE("/:id", r.rbacMiddleware.SuperAdminOnly(), r.koperasiHandler.DeleteKoperasi)

		// Anggota Management
		koperasi.POST("/anggota", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.CreateAnggota)
		koperasi.GET("/:id/anggota", r.koperasiHandler.GetAnggotaList)
		koperasi.GET("/anggota/:id", r.koperasiHandler.GetAnggota)
		koperasi.PUT("/anggota/:id/status", r.rbacMiddleware.AdminOnly(), r.koperasiHandler.UpdateAnggotaStatus)
	}
}