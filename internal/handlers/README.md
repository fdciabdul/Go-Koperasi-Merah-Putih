# Handlers Layer - HTTP Request Processing

## Purpose
Layer ini bertanggung jawab untuk menangani HTTP requests, validation input, dan formatting responses. Handlers adalah entry point dari API yang menerima request dari client.

## Structure Pattern
```
handlers/
├── user_handler.go          # User management endpoints
├── koperasi_handler.go       # Koperasi CRUD operations
├── financial_handler.go      # Financial management endpoints
├── klinik_handler.go         # Healthcare services endpoints
└── ...                      # Other domain handlers
```

## Code Standards

### Handler Structure
```go
type EntityHandler struct {
    entityService *services.EntityService
    // other service dependencies
}

func NewEntityHandler(entityService *services.EntityService) *EntityHandler {
    return &EntityHandler{entityService: entityService}
}
```

### Method Pattern
```go
func (h *EntityHandler) CreateEntity(c *gin.Context) {
    // 1. Input validation
    var req services.CreateEntityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Extract context (user_id, tenant_id, etc.)
    userID, _ := c.Get("user_id")
    req.CreatedBy = userID.(uint64)

    // 3. Call service
    result, err := h.entityService.CreateEntity(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 4. Return success response
    c.JSON(http.StatusCreated, gin.H{
        "message": "Entity created successfully",
        "data":    result,
    })
}
```

## Response Format Standards

### Success Response
```go
// Single resource
c.JSON(http.StatusOK, gin.H{
    "message": "Success message",
    "data":    result,
})

// List response dengan pagination
c.JSON(http.StatusOK, gin.H{
    "data": results,
    "page": page,
    "limit": limit,
    "total": total,
})
```

### Error Response
```go
// Bad request
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Invalid input: " + err.Error(),
})

// Not found
c.JSON(http.StatusNotFound, gin.H{
    "error": "Resource not found",
})

// Server error
c.JSON(http.StatusInternalServerError, gin.H{
    "error": err.Error(),
})
```

## HTTP Status Codes

| Operation | Success Code | Description |
|-----------|-------------|-------------|
| GET       | 200         | Resource found |
| POST      | 201         | Resource created |
| PUT       | 200         | Resource updated |
| DELETE    | 200         | Resource deleted |
| GET (list)| 200         | List retrieved |

| Error Type | Code | When to Use |
|-----------|------|-------------|
| Bad Request | 400 | Invalid input, validation errors |
| Unauthorized | 401 | Missing/invalid auth token |
| Forbidden | 403 | No permission for action |
| Not Found | 404 | Resource doesn't exist |
| Conflict | 409 | Duplicate resource |
| Server Error | 500 | Internal application errors |

## Parameter Handling

### Path Parameters
```go
func (h *Handler) GetEntity(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }
    // Continue processing...
}
```

### Query Parameters
```go
func (h *Handler) GetEntityList(c *gin.Context) {
    // Pagination
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    // Filters
    search := c.Query("search")
    status := c.Query("status")

    // Convert and validate
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }
    // Continue processing...
}
```

## Context Values

### Standard Context Keys
- `user_id` - ID of authenticated user
- `tenant_id` - Multi-tenant isolation
- `koperasi_id` - Current koperasi context (when applicable)
- `role` - User role for authorization

### Extracting Context
```go
func (h *Handler) SomeMethod(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }

    tenantID, _ := c.Get("tenant_id")
    // Use in service calls...
}
```

## Validation Rules

### Request Validation
- Always use `c.ShouldBindJSON()` untuk JSON payloads
- Validate required fields dengan struct tags
- Handle binding errors gracefully
- Return clear error messages

### Business Validation
- Delegate business validation ke service layer
- Handle service errors dan convert ke appropriate HTTP responses
- Don't implement business rules di handlers

## Indonesian Specific Patterns

### Koperasi Context
```go
// For koperasi-specific endpoints
func (h *Handler) GetKoperasiData(c *gin.Context) {
    koperasiIDStr := c.Param("koperasi_id")
    koperasiID, err := strconv.ParseUint(koperasiIDStr, 10, 64)

    // Validate user has access to this koperasi
    // This should be handled by RBAC middleware
}
```

### Indonesian ID Validation
```go
// NIK validation example
func validateNIK(nik string) bool {
    return len(nik) == 16 && isNumeric(nik)
}
```

## Common Patterns

### File Upload
```go
func (h *Handler) UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Validate file type, size, etc.
    // Process upload via service
}
```

### Bulk Operations
```go
func (h *Handler) BulkCreate(c *gin.Context) {
    var req []services.CreateEntityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    results, err := h.entityService.BulkCreate(req)
    // Handle response...
}
```

## Error Handling Patterns

### Service Error Mapping
```go
func (h *Handler) handleServiceError(c *gin.Context, err error) {
    switch {
    case strings.Contains(err.Error(), "not found"):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case strings.Contains(err.Error(), "duplicate"):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    case strings.Contains(err.Error(), "validation"):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }
}
```

## Testing Guidelines

### Handler Testing Pattern
```go
func TestCreateEntity(t *testing.T) {
    // Setup
    mockService := &mocks.EntityService{}
    handler := NewEntityHandler(mockService)

    // Test data
    reqBody := `{"name": "test"}`

    // Execute
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/entities", strings.NewReader(reqBody))

    handler.CreateEntity(c)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Security Considerations

- Never expose internal errors ke client
- Validate all inputs
- Use middleware untuk authentication/authorization
- Sanitize output data
- Log security-relevant events

## Performance Tips

- Use pagination untuk list endpoints
- Implement caching where appropriate
- Avoid N+1 queries (handle di repository layer)
- Use streaming untuk large responses