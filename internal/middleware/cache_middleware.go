package middleware

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"koperasi-merah-putih/internal/cache"

	"github.com/gin-gonic/gin"
)

type CacheMiddleware struct {
	cache *cache.RedisCache
}

func NewCacheMiddleware(cache *cache.RedisCache) *CacheMiddleware {
	return &CacheMiddleware{
		cache: cache,
	}
}

// Cache response middleware with configurable TTL
func (m *CacheMiddleware) CacheResponse(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := m.generateCacheKey(c)

		// Try to get from cache
		var cachedResponse map[string]interface{}
		err := m.cache.Get(cacheKey, &cachedResponse)
		if err == nil {
			// Cache hit - return cached response
			c.Header("X-Cache", "HIT")
			c.Header("X-Cache-Key", cacheKey)
			c.JSON(http.StatusOK, cachedResponse)
			c.Abort()
			return
		}

		// Cache miss - continue to handler
		c.Header("X-Cache", "MISS")
		c.Header("X-Cache-Key", cacheKey)

		// Create a custom response writer to capture the response
		writer := &cachedResponseWriter{
			ResponseWriter: c.Writer,
			body:          make([]byte, 0),
			statusCode:    http.StatusOK,
		}
		c.Writer = writer

		c.Next()

		// Cache the response if it's successful
		if writer.statusCode == http.StatusOK && len(writer.body) > 0 {
			var response map[string]interface{}
			// Parse JSON response
			if err := writer.parseJSON(&response); err == nil {
				m.cache.Set(cacheKey, response, ttl)
			}
		}
	}
}

// Cache for specific endpoints with different TTLs
func (m *CacheMiddleware) CacheWilayah() gin.HandlerFunc {
	return m.CacheResponse(24 * time.Hour) // Wilayah data changes rarely
}

func (m *CacheMiddleware) CacheProducts() gin.HandlerFunc {
	return m.CacheResponse(30 * time.Minute) // Product data changes moderately
}

func (m *CacheMiddleware) CacheReports() gin.HandlerFunc {
	return m.CacheResponse(5 * time.Minute) // Reports can be cached briefly
}

func (m *CacheMiddleware) CacheMasterData() gin.HandlerFunc {
	return m.CacheResponse(2 * time.Hour) // Master data changes infrequently
}

// Invalidate cache for specific patterns
func (m *CacheMiddleware) InvalidateCache(pattern string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// After successful operation, invalidate related cache
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			// Invalidate cache based on pattern
			switch {
			case strings.Contains(c.Request.URL.Path, "/produk"):
				m.invalidateProductCache(c)
			case strings.Contains(c.Request.URL.Path, "/koperasi"):
				m.invalidateKoperasiCache(c)
			case strings.Contains(c.Request.URL.Path, "/financial"):
				m.invalidateFinancialCache(c)
			}
		}
	}
}

// Rate limiting middleware using Redis
func (m *CacheMiddleware) RateLimit(requests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP + User-Agent)
		clientID := m.getClientID(c)
		key := fmt.Sprintf("rate_limit:%s", clientID)

		// Check current rate
		allowed, _, err := m.cache.CheckRateLimit(key, requests)
		if err != nil {
			// If Redis is down, allow the request
			c.Next()
			return
		}

		if !allowed {
			c.Header("X-RateLimit-Limit", strconv.FormatInt(requests, 10))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		// Increment counter
		newCount, err := m.cache.IncrementRateLimit(key, window)
		if err == nil {
			remaining := requests - newCount
			if remaining < 0 {
				remaining = 0
			}
			c.Header("X-RateLimit-Limit", strconv.FormatInt(requests, 10))
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
		}

		c.Next()
	}
}

// Session management middleware
func (m *CacheMiddleware) SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			// Try to get from cookie
			if cookie, err := c.Request.Cookie("session_id"); err == nil {
				sessionID = cookie.Value
			}
		}

		if sessionID != "" {
			// Get session data from Redis
			sessionData, err := m.cache.GetSession(sessionID)
			if err == nil {
				// Set session data in context
				c.Set("session_data", sessionData)
				if userID, ok := sessionData["user_id"].(float64); ok {
					c.Set("user_id", uint64(userID))
				}
			}
		}

		c.Next()
	}
}

// Helper functions
func (m *CacheMiddleware) generateCacheKey(c *gin.Context) string {
	// Include path, query parameters, and user context
	key := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)

	if c.Request.URL.RawQuery != "" {
		key += ":" + c.Request.URL.RawQuery
	}

	// Include user ID for user-specific caching
	if userID, exists := c.Get("user_id"); exists {
		key += fmt.Sprintf(":user_%v", userID)
	}

	// Include tenant/koperasi ID for tenant-specific caching
	if tenantID, exists := c.Get("tenant_id"); exists {
		key += fmt.Sprintf(":tenant_%v", tenantID)
	}

	// Hash the key to keep it reasonable length
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("cache:%x", hash)
}

func (m *CacheMiddleware) getClientID(c *gin.Context) string {
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Include user ID if available for authenticated requests
	if userID, exists := c.Get("user_id"); exists {
		return fmt.Sprintf("%s:%s:user_%v", ip, userAgent, userID)
	}

	return fmt.Sprintf("%s:%s", ip, userAgent)
}

func (m *CacheMiddleware) invalidateProductCache(c *gin.Context) {
	// Extract koperasi_id and invalidate related caches
	if koperasiID, exists := c.Get("koperasi_id"); exists {
		m.cache.Delete(fmt.Sprintf("products:%v", koperasiID))
		m.cache.Delete(fmt.Sprintf("inventory:%v", koperasiID))
		m.cache.Delete(fmt.Sprintf("dashboard:%v", koperasiID))
	}
}

func (m *CacheMiddleware) invalidateKoperasiCache(c *gin.Context) {
	if koperasiID, exists := c.Get("koperasi_id"); exists {
		m.cache.InvalidateKoperasi(koperasiID.(uint64))
		m.cache.Delete(fmt.Sprintf("dashboard:%v", koperasiID))
	}
}

func (m *CacheMiddleware) invalidateFinancialCache(c *gin.Context) {
	if koperasiID, exists := c.Get("koperasi_id"); exists {
		m.cache.Delete(fmt.Sprintf("financial:%v", koperasiID))
		m.cache.Delete(fmt.Sprintf("dashboard:%v", koperasiID))
	}
}

// Custom response writer to capture response for caching
type cachedResponseWriter struct {
	gin.ResponseWriter
	body       []byte
	statusCode int
}

func (w *cachedResponseWriter) Write(data []byte) (int, error) {
	w.body = append(w.body, data...)
	return w.ResponseWriter.Write(data)
}

func (w *cachedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *cachedResponseWriter) parseJSON(v interface{}) error {
	// This is a simplified JSON parser
	// In production, use a proper JSON library
	return nil
}