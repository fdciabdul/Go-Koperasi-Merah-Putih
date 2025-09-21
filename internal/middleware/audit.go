package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"koperasi-merah-putih/internal/models/cassandra"
	cassandraRepo "koperasi-merah-putih/internal/repository/cassandra"
)

type AuditMiddleware struct {
	analyticsRepo *cassandraRepo.AnalyticsRepository
}

func NewAuditMiddleware(analyticsRepo *cassandraRepo.AnalyticsRepository) *AuditMiddleware {
	return &AuditMiddleware{analyticsRepo: analyticsRepo}
}

func (a *AuditMiddleware) AuditLogger() gin.HandlerFunc {
	return gin.LoggerWithWriter(gin.DefaultWriter)
}

func (a *AuditMiddleware) TransactionLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		userID := uint64(0)
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint64); ok {
				userID = id
			}
		}

		tenantID := uint64(1)
		if tid, exists := c.Get("tenant_id"); exists {
			if id, ok := tid.(uint64); ok {
				tenantID = id
			}
		}

		koperasiID := uint64(0)
		if kid, exists := c.Get("koperasi_id"); exists {
			if id, ok := kid.(uint64); ok {
				koperasiID = id
			}
		}

		now := time.Now()
		userActivity := &cassandra.UserActivityLog{
			ID:          gocql.TimeUUID(),
			TenantID:    tenantID,
			KoperasiID:  koperasiID,
			UserID:      userID,
			Activity:    method + " " + path,
			Module:      extractModule(path),
			Description: extractDescription(method, path, statusCode),
			IPAddress:   c.ClientIP(),
			UserAgent:   c.Request.UserAgent(),
			CreatedAt:   now,
			Year:        now.Year(),
			Month:       int(now.Month()),
			Day:         now.Day(),
			Hour:        now.Hour(),
		}

		go a.analyticsRepo.InsertUserActivityLog(userActivity)

		if isTransactionEndpoint(path) {
			transactionLog := &cassandra.TransactionLog{
				ID:              gocql.TimeUUID(),
				TenantID:        tenantID,
				KoperasiID:      koperasiID,
				TransactionType: extractTransactionType(path),
				UserID:          userID,
				Status:          getStatusFromCode(statusCode),
				Description:     string(requestBody),
				IPAddress:       c.ClientIP(),
				UserAgent:       c.Request.UserAgent(),
				CreatedAt:       now,
				Year:            now.Year(),
				Month:           int(now.Month()),
				Day:             now.Day(),
			}

			go a.analyticsRepo.InsertTransactionLog(transactionLog)
		}

		performanceMetrics := &cassandra.PerformanceMetrics{
			ID:           gocql.TimeUUID(),
			TenantID:     tenantID,
			KoperasiID:   koperasiID,
			MetricName:   "response_time",
			MetricValue:  float64(duration.Milliseconds()),
			MetricUnit:   "ms",
			EndpointPath: path,
			HttpMethod:   method,
			ResponseTime: float64(duration.Milliseconds()),
			StatusCode:   statusCode,
			RequestSize:  int64(len(requestBody)),
			ResponseSize: int64(c.Writer.Size()),
			CreatedAt:    now,
			Year:         now.Year(),
			Month:        int(now.Month()),
			Day:          now.Day(),
			Hour:         now.Hour(),
			Minute:       now.Minute(),
		}

		go a.analyticsRepo.InsertPerformanceMetrics(performanceMetrics)
	}
}

func (a *AuditMiddleware) ErrorLogger() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(error); ok {
			now := time.Now()
			errorLog := &cassandra.ErrorLog{
				ID:           gocql.TimeUUID(),
				TenantID:     1,
				KoperasiID:   0,
				UserID:       0,
				ErrorType:    "panic",
				ErrorMessage: err.Error(),
				Module:       extractModule(c.Request.URL.Path),
				Function:     c.Request.URL.Path,
				IPAddress:    c.ClientIP(),
				UserAgent:    c.Request.UserAgent(),
				CreatedAt:    now,
				Year:         now.Year(),
				Month:        int(now.Month()),
				Day:          now.Day(),
			}

			go a.analyticsRepo.InsertErrorLog(errorLog)
		}
		c.AbortWithStatus(500)
	})
}

func extractModule(path string) string {
	if len(path) < 8 {
		return "unknown"
	}

	switch {
	case contains(path, "/ppob"):
		return "ppob"
	case contains(path, "/users"):
		return "users"
	case contains(path, "/payments"):
		return "payments"
	case contains(path, "/koperasi"):
		return "koperasi"
	case contains(path, "/anggota"):
		return "anggota"
	case contains(path, "/simpan-pinjam"):
		return "simpan_pinjam"
	case contains(path, "/klinik"):
		return "klinik"
	case contains(path, "/financial"):
		return "financial"
	default:
		return "general"
	}
}

func extractDescription(method, path string, statusCode int) string {
	return method + " " + path + " - Status: " + string(rune(statusCode))
}

func extractTransactionType(path string) string {
	switch {
	case contains(path, "/ppob"):
		return "ppob"
	case contains(path, "/payments"):
		return "payment"
	case contains(path, "/simpan-pinjam"):
		return "simpan_pinjam"
	case contains(path, "/klinik"):
		return "klinik"
	default:
		return "other"
	}
}

func isTransactionEndpoint(path string) bool {
	transactionPaths := []string{
		"/ppob/transactions",
		"/payments",
		"/simpan-pinjam/transactions",
		"/klinik/kunjungan",
		"/users/register",
	}

	for _, tp := range transactionPaths {
		if contains(path, tp) {
			return true
		}
	}
	return false
}

func getStatusFromCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "success"
	case code >= 400 && code < 500:
		return "failed"
	case code >= 500:
		return "error"
	default:
		return "unknown"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr))
}