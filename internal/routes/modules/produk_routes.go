package modules

import (
	"github.com/gin-gonic/gin"
	"go_koperasi/internal/handlers"
	"go_koperasi/internal/middleware"
)

type ProdukRoutes struct {
	produkHandler  *handlers.ProdukHandler
	rbacMiddleware *middleware.RBACMiddleware
}

func NewProdukRoutes(produkHandler *handlers.ProdukHandler, rbacMiddleware *middleware.RBACMiddleware) *ProdukRoutes {
	return &ProdukRoutes{
		produkHandler:  produkHandler,
		rbacMiddleware: rbacMiddleware,
	}
}

func (r *ProdukRoutes) SetupRoutes(router *gin.RouterGroup) {
	produk := router.Group("/produk")
	produk.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		// Master Data
		produk.POST("/kategori", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreateKategoriProduk)
		produk.GET("/kategori", r.produkHandler.GetAllKategoriProduk)
		produk.GET("/kategori/:id", r.produkHandler.GetKategoriProdukByID)

		produk.POST("/satuan", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreateSatuanProduk)
		produk.GET("/satuan", r.produkHandler.GetAllSatuanProduk)

		// Supplier Management
		produk.POST("/supplier", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreateSupplier)
		produk.GET("/:koperasi_id/supplier", r.produkHandler.GetSuppliersByKoperasi)

		// Product Management
		produk.POST("", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreateProduk)
		produk.GET("/:koperasi_id", r.produkHandler.GetProduksByKoperasi)
		produk.GET("/detail/:id", r.produkHandler.GetProdukByID)
		produk.GET("/barcode/:barcode", r.produkHandler.GetProdukByBarcode)
		produk.POST("/:id/generate-barcode", r.rbacMiddleware.AdminOnly(), r.produkHandler.GenerateBarcode)

		// Purchase Management
		produk.POST("/purchase-order", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreatePurchaseOrder)
		produk.GET("/:koperasi_id/purchase-order", r.produkHandler.GetPurchaseOrdersByKoperasi)

		// Transaction Management
		produk.POST("/pembelian", r.rbacMiddleware.AdminOnly(), r.produkHandler.CreatePembelian)
		produk.POST("/penjualan", r.produkHandler.CreatePenjualan)

		// Reports
		produk.GET("/:koperasi_id/stok-report", r.rbacMiddleware.AdminOnly(), r.produkHandler.GetStokReport)
		produk.GET("/:koperasi_id/stok-rendah", r.rbacMiddleware.AdminOnly(), r.produkHandler.GetProdukStokRendah)
		produk.GET("/:koperasi_id/expiring-soon", r.rbacMiddleware.AdminOnly(), r.produkHandler.GetProdukExpiringSoon)
	}
}