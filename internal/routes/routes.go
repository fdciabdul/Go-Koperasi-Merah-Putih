package routes

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type Routes struct {
	userHandler    *handlers.UserHandler
	paymentHandler *handlers.PaymentHandler
	ppobHandler    *handlers.PPOBHandler
}

func NewRoutes(
	userHandler *handlers.UserHandler,
	paymentHandler *handlers.PaymentHandler,
	ppobHandler *handlers.PPOBHandler,
) *Routes {
	return &Routes{
		userHandler:    userHandler,
		paymentHandler: paymentHandler,
		ppobHandler:    ppobHandler,
	}
}

func (r *Routes) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	public := api.Group("")
	{
		public.POST("/users/register", r.userHandler.RegisterUser)
		public.POST("/payments/midtrans/callback", r.paymentHandler.HandleMidtransCallback)
		public.POST("/payments/xendit/callback", r.paymentHandler.HandleXenditCallback)
		public.PUT("/users/verify-payment/:payment_id", r.userHandler.VerifyPayment)
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/users/registrations/:id/approve", r.userHandler.ApproveRegistration)
		protected.PUT("/users/registrations/:id/reject", r.userHandler.RejectRegistration)
		protected.POST("/payments", r.paymentHandler.CreatePayment)
	}

	ppob := api.Group("/ppob")
	{
		ppob.GET("/kategoris", r.ppobHandler.GetKategoriList)
		ppob.GET("/kategoris/:kategori_id/produks", r.ppobHandler.GetProdukByKategori)
		ppob.POST("/transactions", r.ppobHandler.CreateTransaction)
	}

	ppobProtected := ppob.Group("")
	ppobProtected.Use(middleware.AuthMiddleware())
	{
		ppobProtected.POST("/settlements", r.ppobHandler.CreateSettlement)
	}

	koperasi := api.Group("/koperasi")
	koperasi.Use(middleware.AuthMiddleware())
	{
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
	}
}