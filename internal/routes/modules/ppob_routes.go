package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type PPOBRoutes struct {
	ppobHandler    *handlers.PPOBHandler
	rbacMiddleware *middleware.RBACMiddleware
}

func NewPPOBRoutes(ppobHandler *handlers.PPOBHandler, rbacMiddleware *middleware.RBACMiddleware) *PPOBRoutes {
	return &PPOBRoutes{
		ppobHandler:    ppobHandler,
		rbacMiddleware: rbacMiddleware,
	}
}

func (r *PPOBRoutes) SetupRoutes(router *gin.RouterGroup) {
	ppob := router.Group("/ppob")
	{
		// Public PPOB endpoints
		ppob.GET("/kategoris", r.ppobHandler.GetKategoriList)
		ppob.GET("/kategoris/:kategori_id/produks", r.ppobHandler.GetProdukByKategori)
		ppob.POST("/transactions", r.rbacMiddleware.PPOBAccess(), r.ppobHandler.CreateTransaction)
	}

	// Protected PPOB endpoints
	ppobProtected := ppob.Group("")
	ppobProtected.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireKoperasiAccess())
	{
		ppobProtected.POST("/settlements", r.rbacMiddleware.AdminOnly(), r.ppobHandler.CreateSettlement)
	}
}