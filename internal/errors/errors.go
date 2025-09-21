package errors

import (
	"fmt"
	"net/http"
)

// Custom error types
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Predefined error codes
const (
	// Authentication & Authorization
	ErrUnauthorized     = "UNAUTHORIZED"
	ErrForbidden        = "FORBIDDEN"
	ErrInvalidToken     = "INVALID_TOKEN"
	ErrTokenExpired     = "TOKEN_EXPIRED"
	ErrInvalidCredentials = "INVALID_CREDENTIALS"

	// Validation
	ErrValidation       = "VALIDATION_ERROR"
	ErrRequiredField    = "REQUIRED_FIELD"
	ErrInvalidFormat    = "INVALID_FORMAT"
	ErrInvalidValue     = "INVALID_VALUE"
	ErrDuplicateEntry   = "DUPLICATE_ENTRY"

	// Business Logic
	ErrInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrInsufficientStock = "INSUFFICIENT_STOCK"
	ErrTransactionFailed = "TRANSACTION_FAILED"
	ErrBusinessRule      = "BUSINESS_RULE_VIOLATION"
	ErrOperationFailed   = "OPERATION_FAILED"

	// Resource
	ErrNotFound         = "NOT_FOUND"
	ErrAlreadyExists    = "ALREADY_EXISTS"
	ErrConflict         = "CONFLICT"
	ErrResourceLocked   = "RESOURCE_LOCKED"

	// System
	ErrInternalServer   = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError    = "DATABASE_ERROR"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrTimeoutError     = "TIMEOUT_ERROR"
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"

	// Payment
	ErrPaymentFailed    = "PAYMENT_FAILED"
	ErrPaymentTimeout   = "PAYMENT_TIMEOUT"
	ErrPaymentCancelled = "PAYMENT_CANCELLED"
	ErrInvalidPayment   = "INVALID_PAYMENT"

	// File & Upload
	ErrFileTooBig       = "FILE_TOO_BIG"
	ErrInvalidFileType  = "INVALID_FILE_TYPE"
	ErrUploadFailed     = "UPLOAD_FAILED"
)

// Error constructors
func NewAppError(code, message string, status int, details ...string) *AppError {
	err := &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// Authentication & Authorization errors
func NewUnauthorizedError(message ...string) *AppError {
	msg := "Unauthorized access"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrUnauthorized, msg, http.StatusUnauthorized)
}

func NewForbiddenError(message ...string) *AppError {
	msg := "Access forbidden"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrForbidden, msg, http.StatusForbidden)
}

func NewInvalidTokenError(message ...string) *AppError {
	msg := "Invalid token"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrInvalidToken, msg, http.StatusUnauthorized)
}

func NewTokenExpiredError(message ...string) *AppError {
	msg := "Token expired"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrTokenExpired, msg, http.StatusUnauthorized)
}

func NewInvalidCredentialsError(message ...string) *AppError {
	msg := "Invalid credentials"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrInvalidCredentials, msg, http.StatusUnauthorized)
}

// Validation errors
func NewValidationError(message string, details ...string) *AppError {
	return NewAppError(ErrValidation, message, http.StatusBadRequest, details...)
}

func NewRequiredFieldError(field string) *AppError {
	return NewAppError(ErrRequiredField, fmt.Sprintf("Field '%s' is required", field), http.StatusBadRequest)
}

func NewInvalidFormatError(field string) *AppError {
	return NewAppError(ErrInvalidFormat, fmt.Sprintf("Invalid format for field '%s'", field), http.StatusBadRequest)
}

func NewInvalidValueError(field, value string) *AppError {
	return NewAppError(ErrInvalidValue, fmt.Sprintf("Invalid value '%s' for field '%s'", value, field), http.StatusBadRequest)
}

func NewDuplicateEntryError(field string) *AppError {
	return NewAppError(ErrDuplicateEntry, fmt.Sprintf("Duplicate entry for field '%s'", field), http.StatusConflict)
}

// Business Logic errors
func NewInsufficientFundsError(available, required float64) *AppError {
	return NewAppError(ErrInsufficientFunds,
		fmt.Sprintf("Insufficient funds. Available: %.2f, Required: %.2f", available, required),
		http.StatusBadRequest)
}

func NewInsufficientStockError(product string, available, required int) *AppError {
	return NewAppError(ErrInsufficientStock,
		fmt.Sprintf("Insufficient stock for %s. Available: %d, Required: %d", product, available, required),
		http.StatusBadRequest)
}

func NewTransactionFailedError(message string) *AppError {
	return NewAppError(ErrTransactionFailed, message, http.StatusBadRequest)
}

func NewBusinessRuleError(message string) *AppError {
	return NewAppError(ErrBusinessRule, message, http.StatusBadRequest)
}

func NewOperationFailedError(operation, reason string) *AppError {
	return NewAppError(ErrOperationFailed,
		fmt.Sprintf("Operation '%s' failed: %s", operation, reason),
		http.StatusBadRequest)
}

// Resource errors
func NewNotFoundError(resource string, id interface{}) *AppError {
	return NewAppError(ErrNotFound,
		fmt.Sprintf("%s with ID %v not found", resource, id),
		http.StatusNotFound)
}

func NewAlreadyExistsError(resource string, identifier string) *AppError {
	return NewAppError(ErrAlreadyExists,
		fmt.Sprintf("%s with %s already exists", resource, identifier),
		http.StatusConflict)
}

func NewConflictError(message string) *AppError {
	return NewAppError(ErrConflict, message, http.StatusConflict)
}

func NewResourceLockedError(resource string) *AppError {
	return NewAppError(ErrResourceLocked,
		fmt.Sprintf("%s is currently locked by another operation", resource),
		http.StatusLocked)
}

// System errors
func NewInternalServerError(message ...string) *AppError {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}
	return NewAppError(ErrInternalServer, msg, http.StatusInternalServerError)
}

func NewDatabaseError(operation string, err error) *AppError {
	message := fmt.Sprintf("Database operation '%s' failed", operation)
	details := ""
	if err != nil {
		details = err.Error()
	}
	return NewAppError(ErrDatabaseError, message, http.StatusInternalServerError, details)
}

func NewServiceUnavailableError(service string) *AppError {
	return NewAppError(ErrServiceUnavailable,
		fmt.Sprintf("Service '%s' is temporarily unavailable", service),
		http.StatusServiceUnavailable)
}

func NewTimeoutError(operation string) *AppError {
	return NewAppError(ErrTimeoutError,
		fmt.Sprintf("Operation '%s' timed out", operation),
		http.StatusRequestTimeout)
}

func NewRateLimitExceededError(limit int, window string) *AppError {
	return NewAppError(ErrRateLimitExceeded,
		fmt.Sprintf("Rate limit exceeded. Limit: %d requests per %s", limit, window),
		http.StatusTooManyRequests)
}

// Payment errors
func NewPaymentFailedError(reason string) *AppError {
	return NewAppError(ErrPaymentFailed,
		fmt.Sprintf("Payment failed: %s", reason),
		http.StatusBadRequest)
}

func NewPaymentTimeoutError() *AppError {
	return NewAppError(ErrPaymentTimeout, "Payment timeout", http.StatusRequestTimeout)
}

func NewPaymentCancelledError() *AppError {
	return NewAppError(ErrPaymentCancelled, "Payment cancelled", http.StatusBadRequest)
}

func NewInvalidPaymentError(reason string) *AppError {
	return NewAppError(ErrInvalidPayment,
		fmt.Sprintf("Invalid payment: %s", reason),
		http.StatusBadRequest)
}

// File & Upload errors
func NewFileTooBigError(maxSize string) *AppError {
	return NewAppError(ErrFileTooBig,
		fmt.Sprintf("File size exceeds maximum limit of %s", maxSize),
		http.StatusBadRequest)
}

func NewInvalidFileTypeError(allowedTypes []string) *AppError {
	return NewAppError(ErrInvalidFileType,
		fmt.Sprintf("Invalid file type. Allowed types: %v", allowedTypes),
		http.StatusBadRequest)
}

func NewUploadFailedError(reason string) *AppError {
	return NewAppError(ErrUploadFailed,
		fmt.Sprintf("File upload failed: %s", reason),
		http.StatusInternalServerError)
}

// Error response format
type ErrorResponse struct {
	Success   bool              `json:"success"`
	Error     ErrorDetails      `json:"error"`
	RequestID string            `json:"request_id,omitempty"`
	Timestamp string            `json:"timestamp"`
}

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

type ValidationErrorResponse struct {
	Success    bool                `json:"success"`
	Error      ErrorDetails        `json:"error"`
	Validation []ValidationError   `json:"validation,omitempty"`
	RequestID  string              `json:"request_id,omitempty"`
	Timestamp  string              `json:"timestamp"`
}

// Error handler utilities
func HandleError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Convert common errors
	errMsg := err.Error()
	switch {
	case contains(errMsg, "duplicate key"):
		return NewDuplicateEntryError("unique constraint")
	case contains(errMsg, "foreign key"):
		return NewValidationError("Invalid reference to related data")
	case contains(errMsg, "not found"):
		return NewNotFoundError("Resource", "specified")
	case contains(errMsg, "timeout"):
		return NewTimeoutError("database operation")
	case contains(errMsg, "connection"):
		return NewServiceUnavailableError("database")
	default:
		return NewInternalServerError("An unexpected error occurred")
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		   (str == substr ||
		    str[:len(substr)] == substr ||
		    str[len(str)-len(substr):] == substr ||
		    containsSubstring(str, substr))
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// HTTP status mapping
func GetHTTPStatus(errCode string) int {
	statusMap := map[string]int{
		ErrUnauthorized:       http.StatusUnauthorized,
		ErrForbidden:          http.StatusForbidden,
		ErrInvalidToken:       http.StatusUnauthorized,
		ErrTokenExpired:       http.StatusUnauthorized,
		ErrInvalidCredentials: http.StatusUnauthorized,
		ErrValidation:         http.StatusBadRequest,
		ErrRequiredField:      http.StatusBadRequest,
		ErrInvalidFormat:      http.StatusBadRequest,
		ErrInvalidValue:       http.StatusBadRequest,
		ErrDuplicateEntry:     http.StatusConflict,
		ErrInsufficientFunds:  http.StatusBadRequest,
		ErrInsufficientStock:  http.StatusBadRequest,
		ErrTransactionFailed:  http.StatusBadRequest,
		ErrBusinessRule:       http.StatusBadRequest,
		ErrOperationFailed:    http.StatusBadRequest,
		ErrNotFound:           http.StatusNotFound,
		ErrAlreadyExists:      http.StatusConflict,
		ErrConflict:           http.StatusConflict,
		ErrResourceLocked:     http.StatusLocked,
		ErrInternalServer:     http.StatusInternalServerError,
		ErrDatabaseError:      http.StatusInternalServerError,
		ErrServiceUnavailable: http.StatusServiceUnavailable,
		ErrTimeoutError:       http.StatusRequestTimeout,
		ErrRateLimitExceeded:  http.StatusTooManyRequests,
		ErrPaymentFailed:      http.StatusBadRequest,
		ErrPaymentTimeout:     http.StatusRequestTimeout,
		ErrPaymentCancelled:   http.StatusBadRequest,
		ErrInvalidPayment:     http.StatusBadRequest,
		ErrFileTooBig:         http.StatusBadRequest,
		ErrInvalidFileType:    http.StatusBadRequest,
		ErrUploadFailed:       http.StatusInternalServerError,
	}

	if status, exists := statusMap[errCode]; exists {
		return status
	}
	return http.StatusInternalServerError
}

// Error context for logging
type ErrorContext struct {
	RequestID   string                 `json:"request_id"`
	UserID      uint64                 `json:"user_id,omitempty"`
	KoperasiID  uint64                 `json:"koperasi_id,omitempty"`
	Operation   string                 `json:"operation"`
	Resource    string                 `json:"resource,omitempty"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Timestamp   string                 `json:"timestamp"`
	Error       *AppError              `json:"error"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
}

func NewErrorContext(requestID, operation string) *ErrorContext {
	return &ErrorContext{
		RequestID: requestID,
		Operation: operation,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
}

func (ec *ErrorContext) WithUser(userID uint64) *ErrorContext {
	ec.UserID = userID
	return ec
}

func (ec *ErrorContext) WithKoperasi(koperasiID uint64) *ErrorContext {
	ec.KoperasiID = koperasiID
	return ec
}

func (ec *ErrorContext) WithResource(resource string) *ErrorContext {
	ec.Resource = resource
	return ec
}

func (ec *ErrorContext) WithRequest(method, path string) *ErrorContext {
	ec.Method = method
	ec.Path = path
	return ec
}

func (ec *ErrorContext) WithParameters(params map[string]interface{}) *ErrorContext {
	ec.Parameters = params
	return ec
}

func (ec *ErrorContext) WithError(err *AppError) *ErrorContext {
	ec.Error = err
	return ec
}

func (ec *ErrorContext) WithStackTrace(trace string) *ErrorContext {
	ec.StackTrace = trace
	return ec
}