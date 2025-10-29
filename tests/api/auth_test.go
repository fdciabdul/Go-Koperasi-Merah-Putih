package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"koperasi-merah-putih/tests/helpers"
)

type AuthTestSuite struct {
	helpers.BaseTestSuite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

// Test Login Endpoint
func (s *AuthTestSuite) TestLogin_Success() {
	// Prepare request body
	loginReq := map[string]interface{}{
		"email":    "admin@demo.local",
		"password": "admin123",
	}

	// Mock database expectations
	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs("admin@demo.local").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "email", "password_hash", "nama_lengkap", "role", "is_active",
		}).AddRow(1, 1, "admin@demo.local", "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", "Administrator", "super_admin", true))

	// Make request
	w := s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")

	// Assert response
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "token")
	s.Contains(w.Body.String(), "Login successful")
}

func (s *AuthTestSuite) TestLogin_InvalidCredentials() {
	loginReq := map[string]interface{}{
		"email":    "admin@demo.local",
		"password": "wrongpassword",
	}

	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs("admin@demo.local").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "email", "password_hash", "nama_lengkap", "role", "is_active",
		}).AddRow(1, 1, "admin@demo.local", "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", "Administrator", "super_admin", true))

	w := s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "invalid email or password")
}

func (s *AuthTestSuite) TestLogin_UserNotFound() {
	loginReq := map[string]interface{}{
		"email":    "notfound@demo.local",
		"password": "password123",
	}

	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs("notfound@demo.local").
		WillReturnError(gorm.ErrRecordNotFound)

	w := s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "invalid email or password")
}

func (s *AuthTestSuite) TestLogin_InactiveUser() {
	loginReq := map[string]interface{}{
		"email":    "inactive@demo.local",
		"password": "test123",
	}

	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs("inactive@demo.local").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "email", "password_hash", "nama_lengkap", "role", "is_active",
		}).AddRow(1, 1, "inactive@demo.local", "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", "Inactive User", "user", false))

	w := s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "account is not active")
}

func (s *AuthTestSuite) TestLogin_InvalidInput() {
	// Test missing email
	loginReq := map[string]interface{}{
		"password": "test123",
	}

	w := s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	s.Equal(http.StatusBadRequest, w.Code)

	// Test missing password
	loginReq = map[string]interface{}{
		"email": "test@demo.local",
	}

	w = s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	s.Equal(http.StatusBadRequest, w.Code)

	// Test invalid email format
	loginReq = map[string]interface{}{
		"email":    "invalid-email",
		"password": "test123",
	}

	w = s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	s.Equal(http.StatusBadRequest, w.Code)

	// Test short password
	loginReq = map[string]interface{}{
		"email":    "test@demo.local",
		"password": "123",
	}

	w = s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	s.Equal(http.StatusBadRequest, w.Code)
}

// Test Register Endpoint
func (s *AuthTestSuite) TestRegister_Success() {
	registerReq := map[string]interface{}{
		"koperasi_id":     1,
		"nik":             "1234567890123457",
		"nama_lengkap":    "New User",
		"email":           "newuser@demo.local",
		"no_telepon":      "081234567892",
		"alamat":          "Jl. Test No. 123",
		"kelurahan_id":    1,
		"password":        "password123",
		"jenis_keanggotaan": "anggota",
		"rencana_simpanan": 1000000,
	}

	// Mock database expectations
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"user_registrations\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/users/register", registerReq, "")

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Registration created successfully")
}

func (s *AuthTestSuite) TestRegister_DuplicateEmail() {
	registerReq := map[string]interface{}{
		"koperasi_id":     1,
		"nik":             "1234567890123458",
		"nama_lengkap":    "Duplicate User",
		"email":           "existing@demo.local",
		"no_telepon":      "081234567893",
		"alamat":          "Jl. Test No. 124",
		"kelurahan_id":    1,
		"password":        "password123",
		"jenis_keanggotaan": "anggota",
		"rencana_simpanan": 1000000,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"user_registrations\"").
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	s.MockDB.ExpectRollback()

	w := s.MakeRequest("POST", "/api/v1/users/register", registerReq, "")

	s.Equal(http.StatusInternalServerError, w.Code)
}

func (s *AuthTestSuite) TestRegister_InvalidInput() {
	// Test missing required fields
	registerReq := map[string]interface{}{
		"nama_lengkap": "Test User",
	}

	w := s.MakeRequest("POST", "/api/v1/users/register", registerReq, "")
	s.Equal(http.StatusBadRequest, w.Code)

	// Test invalid NIK
	registerReq = map[string]interface{}{
		"koperasi_id":     1,
		"nik":             "123", // Too short
		"nama_lengkap":    "Test User",
		"email":           "test@demo.local",
		"no_telepon":      "081234567892",
		"alamat":          "Jl. Test No. 123",
		"kelurahan_id":    1,
		"password":        "password123",
		"jenis_keanggotaan": "anggota",
		"rencana_simpanan": 1000000,
	}

	w = s.MakeRequest("POST", "/api/v1/users/register", registerReq, "")
	s.Equal(http.StatusBadRequest, w.Code)
}

// Test Verify Payment Endpoint
func (s *AuthTestSuite) TestVerifyPayment_Success() {
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"payment_transactions\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "registration_id", "status",
		}).AddRow(1, 1, "pending"))

	s.MockDB.ExpectExec("UPDATE \"payment_transactions\"").
		WithArgs("paid", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectQuery("SELECT \\* FROM \"user_registrations\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status",
		}).AddRow(1, "pending_payment"))

	s.MockDB.ExpectExec("UPDATE \"user_registrations\"").
		WithArgs("pending_approval", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("PUT", "/api/v1/users/verify-payment/1", nil, "")

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Payment verified successfully")
}

func (s *AuthTestSuite) TestVerifyPayment_NotFound() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"payment_transactions\" WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	w := s.MakeRequest("PUT", "/api/v1/users/verify-payment/999", nil, "")

	s.Equal(http.StatusNotFound, w.Code)
}

// Test Approve Registration Endpoint
func (s *AuthTestSuite) TestApproveRegistration_Success() {
	s.MockDB.ExpectBegin()

	// Mock auth check
	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "role",
		}).AddRow(2, "admin"))

	// Mock get registration
	s.MockDB.ExpectQuery("SELECT \\* FROM \"user_registrations\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password_hash", "nama_lengkap", "koperasi_id", "status",
		}).AddRow(1, "newuser@demo.local", "hashedpass", "New User", 1, "pending_approval"))

	// Mock create user
	s.MockDB.ExpectQuery("INSERT INTO \"users\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	// Mock create anggota
	s.MockDB.ExpectQuery("INSERT INTO \"anggota_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Mock update registration
	s.MockDB.ExpectExec("UPDATE \"user_registrations\"").
		WithArgs("approved", 2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("PUT", "/api/v1/users/registrations/1/approve", nil, s.AdminToken)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Registration approved successfully")
}

func (s *AuthTestSuite) TestApproveRegistration_Unauthorized() {
	// Test without token
	w := s.MakeRequest("PUT", "/api/v1/users/registrations/1/approve", nil, "")
	s.Equal(http.StatusUnauthorized, w.Code)

	// Test with user token (non-admin)
	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "role",
		}).AddRow(1, "user"))

	w = s.MakeRequest("PUT", "/api/v1/users/registrations/1/approve", nil, s.Token)
	s.Equal(http.StatusForbidden, w.Code)
}

// Test Reject Registration Endpoint
func (s *AuthTestSuite) TestRejectRegistration_Success() {
	rejectReq := map[string]interface{}{
		"reason": "Invalid documents provided",
	}

	s.MockDB.ExpectBegin()

	// Mock auth check
	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "role",
		}).AddRow(2, "admin"))

	// Mock get registration
	s.MockDB.ExpectQuery("SELECT \\* FROM \"user_registrations\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status",
		}).AddRow(1, "pending_approval"))

	// Mock update registration
	s.MockDB.ExpectExec("UPDATE \"user_registrations\"").
		WithArgs("rejected", "Invalid documents provided", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("PUT", "/api/v1/users/registrations/1/reject", rejectReq, s.AdminToken)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Registration rejected")
}

func (s *AuthTestSuite) TestRejectRegistration_MissingReason() {
	w := s.MakeRequest("PUT", "/api/v1/users/registrations/1/reject", nil, s.AdminToken)
	s.Equal(http.StatusBadRequest, w.Code)
}

// Test Create Payment Endpoint
func (s *AuthTestSuite) TestCreatePayment_Success() {
	paymentReq := map[string]interface{}{
		"registration_id": 1,
		"amount":         1000000,
		"payment_method": "bank_transfer",
	}

	s.MockDB.ExpectBegin()

	// Mock auth check
	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "role",
		}).AddRow(1, "user"))

	// Mock create payment
	s.MockDB.ExpectQuery("INSERT INTO \"payment_transactions\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/payments", paymentReq, s.Token)

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Payment created successfully")
}

func (s *AuthTestSuite) TestCreatePayment_Unauthorized() {
	paymentReq := map[string]interface{}{
		"registration_id": 1,
		"amount":         1000000,
		"payment_method": "bank_transfer",
	}

	w := s.MakeRequest("POST", "/api/v1/payments", paymentReq, "")
	s.Equal(http.StatusUnauthorized, w.Code)
}

// Test Payment Callbacks
func (s *AuthTestSuite) TestMidtransCallback_Success() {
	callbackData := map[string]interface{}{
		"order_id":           "REG-001",
		"transaction_status": "settlement",
		"fraud_status":       "accept",
		"payment_type":       "bank_transfer",
		"gross_amount":       "1000000",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"payment_transactions\" WHERE order_id = \\$1").
		WithArgs("REG-001").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status",
		}).AddRow(1, "pending"))

	s.MockDB.ExpectExec("UPDATE \"payment_transactions\"").
		WithArgs("paid", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/payments/midtrans/callback", callbackData, "")

	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthTestSuite) TestXenditCallback_Success() {
	callbackData := map[string]interface{}{
		"external_id": "REG-002",
		"status":      "PAID",
		"amount":      1000000,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"payment_transactions\" WHERE external_id = \\$1").
		WithArgs("REG-002").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status",
		}).AddRow(2, "pending"))

	s.MockDB.ExpectExec("UPDATE \"payment_transactions\"").
		WithArgs("paid", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/payments/xendit/callback", callbackData, "")

	s.Equal(http.StatusOK, w.Code)
}