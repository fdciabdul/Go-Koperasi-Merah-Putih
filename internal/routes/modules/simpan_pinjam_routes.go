package modules

import (
	"github.com/gin-gonic/gin"
	"go_koperasi/internal/handlers"
	"go_koperasi/internal/middleware"
)

type SimpanPinjamRoutes struct {
	simpanPinjamHandler *handlers.SimpanPinjamHandler
	rbacMiddleware      *middleware.RBACMiddleware
}

func NewSimpanPinjamRoutes(simpanPinjamHandler *handlers.SimpanPinjamHandler, rbacMiddleware *middleware.RBACMiddleware) *SimpanPinjamRoutes {
	return &SimpanPinjamRoutes{
		simpanPinjamHandler: simpanPinjamHandler,
		rbacMiddleware:      rbacMiddleware,
	}
}

func (r *SimpanPinjamRoutes) SetupRoutes(router *gin.RouterGroup) {
	simpanPinjam := router.Group("/simpan-pinjam")
	simpanPinjam.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		// Produk Simpan Pinjam
		simpanPinjam.POST("/produk", r.rbacMiddleware.AdminOnly(), r.simpanPinjamHandler.CreateProduk)
		simpanPinjam.GET("/:koperasi_id/produk", r.simpanPinjamHandler.GetProdukList)

		// Rekening Management
		simpanPinjam.POST("/rekening", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.CreateRekening)
		simpanPinjam.GET("/anggota/:anggota_id/rekening", r.simpanPinjamHandler.GetRekeningByAnggota)

		// Transaksi
		simpanPinjam.POST("/transaksi", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.CreateTransaksi)
		simpanPinjam.GET("/rekening/:rekening_id/transaksi", r.simpanPinjamHandler.GetTransaksiByRekening)

		// Reports & Statistics
		simpanPinjam.GET("/:koperasi_id/statistik", r.rbacMiddleware.AdminOnly(), r.simpanPinjamHandler.GetStatistik)
		simpanPinjam.GET("/pinjaman/jatuh-tempo", r.rbacMiddleware.FinancialAccess(), r.simpanPinjamHandler.GetPinjamanJatuhTempo)
	}
}