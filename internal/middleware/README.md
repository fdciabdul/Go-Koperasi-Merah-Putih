# Middleware Layer - Cross-cutting Concerns

## Purpose
Middleware layer menangani cross-cutting concerns seperti authentication, authorization, logging, validation, dan security. Middleware berjalan sebelum request mencapai handler dan dapat memodifikasi request/response.

## Structure Pattern
```
middleware/
├── auth.go          # Authentication middleware
├── rbac.go          # Role-based access control
├── audit.go         # Audit logging middleware
├── cors.go          # CORS handling
├── rate_limit.go    # Rate limiting
├── validation.go    # Request validation
└── recovery.go      # Panic recovery
```

## Middleware Architecture

### Basic Middleware Pattern
```go
func MiddlewareName() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Pre-processing
        // Validate, authenticate, log, etc.

        // Check conditions
        if !isValid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort() // Stop request processing
            return
        }

        // Set context values
        c.Set("key", value)

        // Continue to next middleware/handler
        c.Next()

        // Post-processing (optional)
        // Cleanup, additional logging, etc.
    }
}
```

### Middleware with Configuration
```go
type MiddlewareConfig struct {
    Timeout time.Duration
    MaxRetries int
    SkipPaths []string
}

func ConfigurableMiddleware(config MiddlewareConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip middleware untuk certain paths
        for _, path := range config.SkipPaths {
            if c.Request.URL.Path == path {
                c.Next()
                return
            }
        }

        // Apply middleware logic with config
        // ...

        c.Next()
    }
}
```

## Authentication Middleware

### JWT Authentication
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token dari header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort()
            return
        }

        // Parse Bearer token
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid authorization format",
            })
            c.Abort()
            return
        }

        // Validate JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            c.Abort()
            return
        }

        // Extract claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", uint64(claims["user_id"].(float64)))
            c.Set("tenant_id", uint64(claims["tenant_id"].(float64)))
            c.Set("role", claims["role"].(string))
            c.Set("email", claims["email"].(string))
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token claims",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Optional Authentication
```go
func OptionalAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            // No auth header, continue without setting user context
            c.Next()
            return
        }

        // Try to authenticate, but don't fail if invalid
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err == nil && token.Valid {
            if claims, ok := token.Claims.(jwt.MapClaims); ok {
                c.Set("user_id", uint64(claims["user_id"].(float64)))
                c.Set("tenant_id", uint64(claims["tenant_id"].(float64)))
            }
        }

        c.Next()
    }
}
```

## RBAC Middleware

### Role-Based Access Control
```go
type RBACMiddleware struct {
    userService *services.UserService
}

func NewRBACMiddleware(userService *services.UserService) *RBACMiddleware {
    return &RBACMiddleware{userService: userService}
}

func (m *RBACMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "User role not found in context",
            })
            c.Abort()
            return
        }

        role := userRole.(string)
        for _, allowedRole := range roles {
            if role == allowedRole {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "error": "Insufficient permissions",
        })
        c.Abort()
    }
}

func (m *RBACMiddleware) AdminOnly() gin.HandlerFunc {
    return m.RequireRole("admin", "super_admin")
}

func (m *RBACMiddleware) SuperAdminOnly() gin.HandlerFunc {
    return m.RequireRole("super_admin")
}
```

### Permission-Based Access Control
```go
func (m *RBACMiddleware) RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "User not authenticated",
            })
            c.Abort()
            return
        }

        // Check if user has permission
        hasPermission, err := m.userService.HasPermission(userID.(uint64), permission)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to check permissions",
            })
            c.Abort()
            return
        }

        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Permission denied",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Koperasi Access Control
```go
func (m *RBACMiddleware) RequireKoperasiAccess() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")
        tenantID, _ := c.Get("tenant_id")

        // Extract koperasi_id dari URL parameter
        koperasiIDStr := c.Param("koperasi_id")
        if koperasiIDStr != "" {
            koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)
            if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "Invalid koperasi ID",
                })
                c.Abort()
                return
            }

            // Check if user has access to this koperasi
            hasAccess, err := m.userService.HasKoperasiAccess(
                userID.(uint64),
                tenantID.(uint64),
                koperasiID,
            )
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Failed to check koperasi access",
                })
                c.Abort()
                return
            }

            if !hasAccess {
                c.JSON(http.StatusForbidden, gin.H{
                    "error": "No access to this koperasi",
                })
                c.Abort()
                return
            }

            c.Set("koperasi_id", koperasiID)
        }

        c.Next()
    }
}
```

## Audit Middleware

### Transaction Logging
```go
type AuditMiddleware struct {
    auditService *services.AuditService
}

func NewAuditMiddleware(auditService *services.AuditService) *AuditMiddleware {
    return &AuditMiddleware{auditService: auditService}
}

func (m *AuditMiddleware) TransactionLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Create request ID
        requestID := generateRequestID()
        c.Set("request_id", requestID)

        // Log request
        m.logRequest(c, requestID)

        // Process request
        c.Next()

        // Log response
        duration := time.Since(start)
        m.logResponse(c, requestID, duration)
    }
}

func (m *AuditMiddleware) logRequest(c *gin.Context, requestID string) {
    var body []byte
    if c.Request.Body != nil {
        body, _ = io.ReadAll(c.Request.Body)
        c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
    }

    auditLog := &models.AuditLog{
        RequestID:    requestID,
        UserID:       getUserID(c),
        TenantID:     getTenantID(c),
        Method:       c.Request.Method,
        Path:         c.Request.URL.Path,
        QueryParams:  c.Request.URL.RawQuery,
        RequestBody:  string(body),
        UserAgent:    c.Request.UserAgent(),
        IP:           c.ClientIP(),
        Timestamp:    time.Now(),
    }

    // Async logging to avoid blocking request
    go m.auditService.LogRequest(auditLog)
}

func (m *AuditMiddleware) logResponse(c *gin.Context, requestID string, duration time.Duration) {
    responseLog := &models.ResponseLog{
        RequestID:    requestID,
        StatusCode:   c.Writer.Status(),
        ResponseSize: c.Writer.Size(),
        Duration:     duration,
        Timestamp:    time.Now(),
    }

    // Async logging
    go m.auditService.LogResponse(responseLog)
}
```

### Sensitive Data Masking
```go
func (m *AuditMiddleware) maskSensitiveData(body string) string {
    // Mask password fields
    re := regexp.MustCompile(`"password":\s*"[^"]*"`)
    body = re.ReplaceAllString(body, `"password":"***"`)

    // Mask credit card numbers
    re = regexp.MustCompile(`"card_number":\s*"[^"]*"`)
    body = re.ReplaceAllString(body, `"card_number":"***"`)

    // Mask NIK
    re = regexp.MustCompile(`"nik":\s*"(\d{4})\d{8}(\d{4})"`)
    body = re.ReplaceAllString(body, `"nik":"$1********$2"`)

    return body
}
```

## Validation Middleware

### Request Validation
```go
func ValidateJSON() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "POST" || c.Request.Method == "PUT" {
            contentType := c.GetHeader("Content-Type")
            if !strings.Contains(contentType, "application/json") {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "Content-Type must be application/json",
                })
                c.Abort()
                return
            }

            // Validate JSON syntax
            var json map[string]interface{}
            if err := c.ShouldBindJSON(&json); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "Invalid JSON format",
                })
                c.Abort()
                return
            }

            // Reset request body untuk handler
            c.Request.Body = io.NopCloser(strings.NewReader(
                fmt.Sprintf("%v", json),
            ))
        }

        c.Next()
    }
}
```

### Tenant Validation
```go
func ValidateTenant() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID, exists := c.Get("tenant_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Tenant ID not found",
            })
            c.Abort()
            return
        }

        // Validate tenant is active
        // This could be cached untuk performance
        isActive := validateTenantActive(tenantID.(uint64))
        if !isActive {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Tenant is inactive",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## Security Middleware

### CORS Middleware
```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

### Rate Limiting
```go
func RateLimitMiddleware(requests int, window time.Duration) gin.HandlerFunc {
    var mu sync.Mutex
    clients := make(map[string]*rate.Limiter)

    return func(c *gin.Context) {
        ip := c.ClientIP()

        mu.Lock()
        limiter, exists := clients[ip]
        if !exists {
            limiter = rate.NewLimiter(rate.Every(window/time.Duration(requests)), requests)
            clients[ip] = limiter
        }
        mu.Unlock()

        if !limiter.Allow() {
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### Security Headers
```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

        c.Next()
    }
}
```

## Recovery Middleware

### Panic Recovery
```go
func RecoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        if err, ok := recovered.(string); ok {
            log.Printf("Panic recovered: %s\n%s", err, debug.Stack())
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
                "request_id": c.GetString("request_id"),
            })
        }
        c.Abort()
    })
}
```

## Middleware Composition

### Middleware Chain
```go
func SetupMiddlewareChain(router *gin.Engine, deps *Dependencies) {
    // Global middleware (order matters!)
    router.Use(
        RecoveryMiddleware(),
        CORSMiddleware(),
        SecurityHeadersMiddleware(),
        RateLimitMiddleware(100, time.Minute),
        deps.AuditMiddleware.TransactionLogger(),
    )

    // API routes with specific middleware
    api := router.Group("/api/v1")
    {
        // Public routes
        public := api.Group("")
        public.Use(OptionalAuthMiddleware())

        // Protected routes
        protected := api.Group("")
        protected.Use(
            AuthMiddleware(),
            ValidateTenant(),
            deps.RBACMiddleware.RequireTenantAccess(),
        )

        // Admin routes
        admin := protected.Group("/admin")
        admin.Use(deps.RBACMiddleware.AdminOnly())
    }
}
```

## Testing Middleware

### Middleware Testing Pattern
```go
func TestAuthMiddleware(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(AuthMiddleware())
    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "success"})
    })

    tests := []struct {
        name           string
        authHeader     string
        expectedStatus int
    }{
        {
            name:           "Valid token",
            authHeader:     "Bearer " + generateValidToken(),
            expectedStatus: 200,
        },
        {
            name:           "No token",
            authHeader:     "",
            expectedStatus: 401,
        },
        {
            name:           "Invalid token",
            authHeader:     "Bearer invalid",
            expectedStatus: 401,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/test", nil)
            if tt.authHeader != "" {
                req.Header.Set("Authorization", tt.authHeader)
            }

            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
```

## Performance Considerations

1. **Use async logging untuk audit middleware**
2. **Cache permission checks untuk RBAC**
3. **Implement circuit breakers untuk external services**
4. **Use connection pooling untuk database checks**
5. **Minimize middleware overhead**
6. **Profile middleware performance**

## Security Best Practices

1. **Always validate and sanitize inputs**
2. **Use HTTPS only in production**
3. **Implement proper CORS policies**
4. **Add security headers**
5. **Rate limit all endpoints**
6. **Log security events**
7. **Use secure JWT practices**
8. **Implement proper session management**