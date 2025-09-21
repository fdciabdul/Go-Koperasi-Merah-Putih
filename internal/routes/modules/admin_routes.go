package modules

import (
	"github.com/gin-gonic/gin"
	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
)

type AdminRoutes struct {
	sequenceHandler *handlers.SequenceHandler
	rbacMiddleware  *middleware.RBACMiddleware
}

func NewAdminRoutes(sequenceHandler *handlers.SequenceHandler, rbacMiddleware *middleware.RBACMiddleware) *AdminRoutes {
	return &AdminRoutes{
		sequenceHandler: sequenceHandler,
		rbacMiddleware:  rbacMiddleware,
	}
}

func (r *AdminRoutes) SetupRoutes(router *gin.RouterGroup) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), r.rbacMiddleware.RequireTenantAccess())
	{
		// Sequence Management
		admin.GET("/sequences", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.GetSequenceList)
		admin.PUT("/sequences/update-value", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.UpdateSequenceValue)
		admin.PUT("/sequences/reset", r.rbacMiddleware.AdminOnly(), r.sequenceHandler.ResetSequence)
	}
}