package routes

import (
	"github.com/gin-gonic/gin"
	"go_koperasi/internal/handlers"
	"go_koperasi/internal/middleware"
	"go_koperasi/internal/routes/modules"
)

type Routes struct {
	// Route modules
	authRoutes       *modules.AuthRoutes
	koperasiRoutes   *modules.KoperasiRoutes
	wilayahRoutes    *modules.WilayahRoutes
	simpanPinjamRoutes *modules.SimpanPinjamRoutes
	ppobRoutes       *modules.PPOBRoutes
	klinikRoutes     *modules.KlinikRoutes
	produkRoutes     *modules.ProdukRoutes
	financialRoutes  *modules.FinancialRoutes
	masterDataRoutes *modules.MasterDataRoutes
	adminRoutes      *modules.AdminRoutes

	// Middleware
	auditMiddleware *middleware.AuditMiddleware
}

func NewRoutes(
	userHandler *handlers.UserHandler,
	paymentHandler *handlers.PaymentHandler,
	ppobHandler *handlers.PPOBHandler,
	koperasiHandler *handlers.KoperasiHandler,
	simpanPinjamHandler *handlers.SimpanPinjamHandler,
	klinikHandler *handlers.KlinikHandler,
	financialHandler *handlers.FinancialHandler,
	wilayahHandler *handlers.WilayahHandler,
	masterDataHandler *handlers.MasterDataHandler,
	sequenceHandler *handlers.SequenceHandler,
	produkHandler *handlers.ProdukHandler,
	rbacMiddleware *middleware.RBACMiddleware,
	auditMiddleware *middleware.AuditMiddleware,
) *Routes {
	return &Routes{
		authRoutes:       modules.NewAuthRoutes(userHandler, paymentHandler),
		koperasiRoutes:   modules.NewKoperasiRoutes(koperasiHandler, rbacMiddleware),
		wilayahRoutes:    modules.NewWilayahRoutes(wilayahHandler),
		simpanPinjamRoutes: modules.NewSimpanPinjamRoutes(simpanPinjamHandler, rbacMiddleware),
		ppobRoutes:       modules.NewPPOBRoutes(ppobHandler, rbacMiddleware),
		klinikRoutes:     modules.NewKlinikRoutes(klinikHandler, rbacMiddleware),
		produkRoutes:     modules.NewProdukRoutes(produkHandler, rbacMiddleware),
		financialRoutes:  modules.NewFinancialRoutes(financialHandler, rbacMiddleware),
		masterDataRoutes: modules.NewMasterDataRoutes(masterDataHandler, rbacMiddleware),
		adminRoutes:      modules.NewAdminRoutes(sequenceHandler, rbacMiddleware),
		auditMiddleware:  auditMiddleware,
	}
}

func (r *Routes) SetupRoutes(router *gin.Engine) {
	// Global middleware
	api := router.Group("/api/v1")
	api.Use(r.auditMiddleware.TransactionLogger())

	// Setup public routes (no authentication required)
	r.authRoutes.SetupPublicRoutes(api)

	// Setup protected routes (authentication required)
	r.authRoutes.SetupProtectedRoutes(api)

	// Setup domain-specific routes
	r.koperasiRoutes.SetupRoutes(api)
	r.wilayahRoutes.SetupRoutes(api)
	r.simpanPinjamRoutes.SetupRoutes(api)
	r.ppobRoutes.SetupRoutes(api)
	r.klinikRoutes.SetupRoutes(api)
	r.produkRoutes.SetupRoutes(api)
	r.financialRoutes.SetupRoutes(api)
	r.masterDataRoutes.SetupRoutes(api)
	r.adminRoutes.SetupRoutes(api)
}