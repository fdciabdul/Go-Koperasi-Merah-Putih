package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type KlinikRoutes struct {
	klinikHandler  *handlers.KlinikHandler
	rbacMiddleware *middleware.RBACMiddleware
}

func NewKlinikRoutes(klinikHandler *handlers.KlinikHandler, rbacMiddleware *middleware.RBACMiddleware) *KlinikRoutes {
	return &KlinikRoutes{
		klinikHandler:  klinikHandler,
		rbacMiddleware: rbacMiddleware,
	}
}

func (r *KlinikRoutes) SetupRoutes(router *gin.RouterGroup) {
	klinik := router.Group("/klinik")
	klinik.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess(), r.rbacMiddleware.KlinikAccess())
	{
		// Pasien Management
		klinik.POST("/pasien", r.klinikHandler.CreatePasien)
		klinik.GET("/:koperasi_id/pasien", r.klinikHandler.GetPasienList)
		klinik.GET("/pasien/:id", r.klinikHandler.GetPasien)
		klinik.GET("/:koperasi_id/pasien/search", r.klinikHandler.SearchPasien)

		// Tenaga Medis Management
		klinik.POST("/tenaga-medis", r.rbacMiddleware.AdminOnly(), r.klinikHandler.CreateTenagaMedis)
		klinik.GET("/:koperasi_id/tenaga-medis", r.klinikHandler.GetTenagaMedisList)

		// Kunjungan Management
		klinik.POST("/kunjungan", r.klinikHandler.CreateKunjungan)
		klinik.GET("/kunjungan/:id", r.klinikHandler.GetKunjungan)
		klinik.GET("/pasien/:id/kunjungan", r.klinikHandler.GetKunjunganByPasien)

		// Obat Management
		klinik.POST("/obat", r.rbacMiddleware.AdminOnly(), r.klinikHandler.CreateObat)
		klinik.GET("/:koperasi_id/obat", r.klinikHandler.GetObatList)
		klinik.GET("/:koperasi_id/obat/search", r.klinikHandler.SearchObat)
		klinik.GET("/:koperasi_id/obat/stok-rendah", r.klinikHandler.GetObatStokRendah)

		// Statistics
		klinik.GET("/:koperasi_id/statistik", r.rbacMiddleware.AdminOnly(), r.klinikHandler.GetStatistik)
	}
}