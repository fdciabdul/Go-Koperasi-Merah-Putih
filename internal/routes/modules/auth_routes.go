package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type AuthRoutes struct {
	userHandler    *handlers.UserHandler
	paymentHandler *handlers.PaymentHandler
}

func NewAuthRoutes(userHandler *handlers.UserHandler, paymentHandler *handlers.PaymentHandler) *AuthRoutes {
	return &AuthRoutes{
		userHandler:    userHandler,
		paymentHandler: paymentHandler,
	}
}

func (r *AuthRoutes) SetupPublicRoutes(router *gin.RouterGroup) {
	router.POST("/users/register", r.userHandler.RegisterUser)
	router.POST("/payments/midtrans/callback", r.paymentHandler.HandleMidtransCallback)
	router.POST("/payments/xendit/callback", r.paymentHandler.HandleXenditCallback)
	router.PUT("/users/verify-payment/:payment_id", r.userHandler.VerifyPayment)
}

func (r *AuthRoutes) SetupProtectedRoutes(router *gin.RouterGroup) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/users/registrations/:id/approve", r.userHandler.ApproveRegistration)
		protected.PUT("/users/registrations/:id/reject", r.userHandler.RejectRegistration)
		protected.POST("/payments", r.paymentHandler.CreatePayment)
	}
}