package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type RBACMiddleware struct {
	db *gorm.DB
}

func NewRBACMiddleware(db *gorm.DB) *RBACMiddleware {
	return &RBACMiddleware{db: db}
}

func (r *RBACMiddleware) RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}

func (r *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format"})
			c.Abort()
			return
		}

		hasPermission, err := r.checkPermission(userRole, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking permissions"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (r *RBACMiddleware) RequireKoperasiAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		uid, ok := userID.(uint64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		var user postgres.User
		err := r.db.First(&user, uid).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.KoperasiID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not assigned to any koperasi"})
			c.Abort()
			return
		}

		c.Set("koperasi_id", user.KoperasiID)
		c.Next()
	}
}

func (r *RBACMiddleware) RequireAnggotaAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		uid, ok := userID.(uint64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		var user postgres.User
		err := r.db.Preload("Anggota").First(&user, uid).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.AnggotaID == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member"})
			c.Abort()
			return
		}

		c.Set("anggota_id", user.AnggotaID)
		c.Next()
	}
}

func (r *RBACMiddleware) RequireTenantAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCode := c.GetHeader("X-Tenant-Code")
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant code required"})
			c.Abort()
			return
		}

		var tenant postgres.Tenant
		err := r.db.Where("tenant_code = ? AND is_active = ?", tenantCode, true).First(&tenant).Error
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or inactive tenant"})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenant.ID)
		c.Set("tenant", tenant)
		c.Next()
	}
}

func (r *RBACMiddleware) SuperAdminOnly() gin.HandlerFunc {
	return r.RequireRole("super_admin")
}

func (r *RBACMiddleware) AdminOnly() gin.HandlerFunc {
	return r.RequireRole("super_admin", "admin_koperasi")
}

func (r *RBACMiddleware) FinancialAccess() gin.HandlerFunc {
	return r.RequireRole("super_admin", "admin_koperasi", "bendahara")
}

func (r *RBACMiddleware) KlinikAccess() gin.HandlerFunc {
	return r.RequireRole("super_admin", "admin_koperasi", "operator")
}

func (r *RBACMiddleware) PPOBAccess() gin.HandlerFunc {
	return r.RequireRole("super_admin", "admin_koperasi", "operator", "anggota")
}

func (r *RBACMiddleware) checkPermission(role, permission string) (bool, error) {
	var count int64
	err := r.db.Table("role_permissions rp").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("rp.role = ? AND p.name = ?", role, permission).
		Count(&count).Error

	return count > 0, err
}