package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"koperasi-merah-putih/internal/handlers"
	"koperasi-merah-putih/internal/middleware"
	"koperasi-merah-putih/internal/models/cassandra"
	postgresModel "koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
	"koperasi-merah-putih/internal/services"
)

// BaseTestSuite provides common test setup
type BaseTestSuite struct {
	suite.Suite
	Router         *gin.Engine
	DB             *gorm.DB
	MockDB         sqlmock.Sqlmock
	TestServer     *httptest.Server
	Token          string
	AdminToken     string
	SuperAdminToken string
}

// SetupSuite runs before all tests
func (s *BaseTestSuite) SetupSuite() {
	// Set test environment
	os.Setenv("APP_ENV", "test")
	os.Setenv("JWT_SECRET", "test-secret-key")

	// Setup mock database
	s.setupMockDB()

	// Setup router
	s.setupRouter()

	// Create test server
	s.TestServer = httptest.NewServer(s.Router)

	// Generate test tokens
	s.generateTestTokens()
}

// setupMockDB creates a mock database connection
func (s *BaseTestSuite) setupMockDB() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)

	s.MockDB = mock

	// Create GORM DB with mock
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	s.Require().NoError(err)

	s.DB = gormDB
}

// setupRouter creates a test router with all routes configured
func (s *BaseTestSuite) setupRouter() {
	gin.SetMode(gin.TestMode)
	s.Router = gin.New()

	// Initialize repositories
	userRepo := postgresRepo.NewUserRepository(s.DB)
	registrationRepo := postgresRepo.NewUserRegistrationRepository(s.DB)
	anggotaRepo := postgresRepo.NewAnggotaKoperasiRepository(s.DB)
	koperasiRepo := postgresRepo.NewKoperasiRepository(s.DB)
	produkRepo := postgresRepo.NewProdukRepository(s.DB)
	financialRepo := postgresRepo.NewFinancialRepository(s.DB)
	simpanPinjamRepo := postgresRepo.NewSimpanPinjamRepository(s.DB)
	ppobRepo := postgresRepo.NewPPOBRepository(s.DB)
	klinikRepo := postgresRepo.NewKlinikRepository(s.DB)
	sequenceRepo := postgresRepo.NewSequenceRepository(s.DB)
	wilayahRepo := postgresRepo.NewWilayahRepository(s.DB)
	masterDataRepo := postgresRepo.NewMasterDataRepository(s.DB)
	paymentRepo := postgresRepo.NewPaymentRepository(s.DB)
	paymentProviderRepo := postgresRepo.NewPaymentProviderRepository(s.DB)

	// Initialize services
	sequenceService := services.NewSequenceService(sequenceRepo)
	paymentService := services.NewPaymentService(paymentRepo, paymentProviderRepo, sequenceService)
	userService := services.NewUserService(userRepo, registrationRepo, anggotaRepo, paymentService, sequenceService)
	koperasiService := services.NewKoperasiService(koperasiRepo, anggotaRepo, wilayahRepo, sequenceService)
	produkService := services.NewProdukService(produkRepo, sequenceRepo)
	financialService := services.NewFinancialService(financialRepo, sequenceService)
	simpanPinjamService := services.NewSimpanPinjamService(simpanPinjamRepo, sequenceService)
	ppobService := services.NewPPOBService(ppobRepo, paymentService, sequenceService)
	klinikService := services.NewKlinikService(klinikRepo, sequenceService)
	wilayahService := services.NewWilayahService(wilayahRepo)
	masterDataService := services.NewMasterDataService(masterDataRepo)
	reportingService := services.NewReportingService(koperasiRepo, anggotaRepo, produkRepo, simpanPinjamRepo, financialRepo, klinikRepo, nil)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	paymentHandler := handlers.NewPaymentHandler(paymentService, userService, ppobService)
	koperasiHandler := handlers.NewKoperasiHandler(koperasiService)
	produkHandler := handlers.NewProdukHandler(produkService)
	financialHandler := handlers.NewFinancialHandler(financialService)
	simpanPinjamHandler := handlers.NewSimpanPinjamHandler(simpanPinjamService)
	ppobHandler := handlers.NewPPOBHandler(ppobService)
	klinikHandler := handlers.NewKlinikHandler(klinikService)
	wilayahHandler := handlers.NewWilayahHandler(wilayahService)
	masterDataHandler := handlers.NewMasterDataHandler(masterDataService)
	sequenceHandler := handlers.NewSequenceHandler(sequenceService)
	reportingHandler := handlers.NewReportingHandler(reportingService)

	// Initialize middleware
	rbacMiddleware := middleware.NewRBACMiddleware(s.DB)

	// For tests, we'll setup routes manually without audit middleware
	s.setupTestRoutes(
		userHandler,
		paymentHandler,
		ppobHandler,
		koperasiHandler,
		simpanPinjamHandler,
		klinikHandler,
		financialHandler,
		wilayahHandler,
		masterDataHandler,
		sequenceHandler,
		produkHandler,
		reportingHandler,
		rbacMiddleware,
	)
}

// setupTestRoutes sets up test routes without audit middleware
func (s *BaseTestSuite) setupTestRoutes(
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
	reportingHandler *handlers.ReportingHandler,
	rbacMiddleware *middleware.RBACMiddleware,
) {
	// Setup API routes
	api := s.Router.Group("/api/v1")

	// Auth routes (public)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/users/register", userHandler.RegisterUser)
	api.PUT("/users/verify-payment/:payment_id", userHandler.VerifyPayment)
	api.PUT("/users/registrations/:id/approve", userHandler.ApproveRegistration)
	api.PUT("/users/registrations/:id/reject", userHandler.RejectRegistration)
	api.POST("/payments", paymentHandler.CreatePayment)
	api.POST("/payments/midtrans/callback", paymentHandler.HandleMidtransCallback)
	api.POST("/payments/xendit/callback", paymentHandler.HandleXenditCallback)

	// Protected routes (basic auth)
	api.GET("/koperasi", koperasiHandler.GetKoperasiList)
	api.GET("/koperasi/:id", koperasiHandler.GetKoperasi)

	// Admin protected routes
	adminRoutes := api.Group("/")
	adminRoutes.POST("/koperasi", koperasiHandler.CreateKoperasi)
	adminRoutes.PUT("/koperasi/:id", koperasiHandler.UpdateKoperasi)
	adminRoutes.DELETE("/koperasi/:id", koperasiHandler.DeleteKoperasi)

	// Add more routes as needed for testing
	// This is a simplified setup for tests
}

// generateTestTokens creates tokens for different user roles
func (s *BaseTestSuite) generateTestTokens() {
	// Regular user token
	token, err := services.GenerateJWT(1, 1, "user", time.Now().Add(24*time.Hour))
	s.Require().NoError(err)
	s.Token = token

	// Admin token
	adminToken, err := services.GenerateJWT(2, 1, "admin", time.Now().Add(24*time.Hour))
	s.Require().NoError(err)
	s.AdminToken = adminToken

	// Super admin token
	superAdminToken, err := services.GenerateJWT(3, 1, "super_admin", time.Now().Add(24*time.Hour))
	s.Require().NoError(err)
	s.SuperAdminToken = superAdminToken
}

// TearDownSuite runs after all tests
func (s *BaseTestSuite) TearDownSuite() {
	if s.TestServer != nil {
		s.TestServer.Close()
	}
}

// MakeRequest helper function to make HTTP requests
func (s *BaseTestSuite) MakeRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var req *http.Request

	if body != nil {
		jsonBody, err := json.Marshal(body)
		s.Require().NoError(err)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()
	s.Router.ServeHTTP(w, req)
	return w
}

// AssertJSONResponse helper to validate JSON responses
func (s *BaseTestSuite) AssertJSONResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedBody interface{}) {
	s.Equal(expectedStatus, w.Code)

	if expectedBody != nil {
		var response interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		s.Require().NoError(err)

		expectedJSON, err := json.Marshal(expectedBody)
		s.Require().NoError(err)

		var expected interface{}
		err = json.Unmarshal(expectedJSON, &expected)
		s.Require().NoError(err)

		s.Equal(expected, response)
	}
}

// MockAnalyticsRepository for testing
type MockAnalyticsRepository struct{}

func (m *MockAnalyticsRepository) InsertTransactionLog(log *cassandra.TransactionLog) error {
	return nil
}

func (m *MockAnalyticsRepository) GetTransactionLogs(koperasiID uint64, year, month int, limit int) ([]cassandra.TransactionLog, error) {
	return nil, nil
}

func (m *MockAnalyticsRepository) InsertPaymentAnalytics(analytics *cassandra.PaymentAnalytics) error {
	return nil
}

func (m *MockAnalyticsRepository) InsertPPOBAnalytics(analytics *cassandra.PPOBAnalytics) error {
	return nil
}

func (m *MockAnalyticsRepository) GetPPOBAnalytics(koperasiID uint64, year, month int) ([]cassandra.PPOBAnalytics, error) {
	return nil, nil
}

func (m *MockAnalyticsRepository) InsertUserActivityLog(log *cassandra.UserActivityLog) error {
	return nil
}

func (m *MockAnalyticsRepository) InsertErrorLog(log *cassandra.ErrorLog) error {
	return nil
}

func (m *MockAnalyticsRepository) InsertPerformanceMetrics(metrics *cassandra.PerformanceMetrics) error {
	return nil
}

func (m *MockAnalyticsRepository) InsertFactKeuanganBulanan(fact *cassandra.FactKeuanganBulanan) error {
	return nil
}

func (m *MockAnalyticsRepository) GetFactKeuanganBulanan(koperasiID uint64, tahun int) ([]cassandra.FactKeuanganBulanan, error) {
	return nil, nil
}

// TestAuditMiddleware for testing
type TestAuditMiddleware struct{}

func (m *TestAuditMiddleware) AuditLogger() gin.HandlerFunc {
	return gin.Logger()
}

func (m *TestAuditMiddleware) TransactionLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func (m *TestAuditMiddleware) ErrorLogger() gin.HandlerFunc {
	return gin.Recovery()
}

// Test data fixtures
func GetTestUser() *postgresModel.User {
	return &postgresModel.User{
		ID:           1,
		TenantID:     1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: test123
		NamaLengkap:  "Test User",
		Role:         "user",
		IsActive:     true,
	}
}

func GetTestKoperasi() *postgresModel.Koperasi {
	return &postgresModel.Koperasi{
		ID:           1,
		TenantID:     1,
		NomorSK:      "TEST-001",
		NIK:          1234567890123456,
		NamaKoperasi: "Test Koperasi",
		NamaSK:       "Test Koperasi SK",
		Email:        "koperasi@test.com",
		Telepon:      "081234567890",
	}
}

func GetTestAnggota() *postgresModel.AnggotaKoperasi {
	return &postgresModel.AnggotaKoperasi{
		ID:            1,
		KoperasiID:    1,
		NIAK:          "TEST001",
		NIK:           "1234567890123456",
		Nama:          "Test Anggota",
		JenisKelamin:  "L",
		Email:         "anggota@test.com",
		Telepon:       "081234567891",
		StatusAnggota: "aktif",
	}
}