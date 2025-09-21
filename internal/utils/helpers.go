package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// String utilities
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func GenerateUUID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])
}

func TruncateString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[:maxLength-3] + "..."
}

func CleanString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func ToTitleCase(str string) string {
	return strings.Title(strings.ToLower(str))
}

// Validation utilities
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidPhoneNumber(phone string) bool {
	// Indonesian phone number validation
	phoneRegex := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)
	return phoneRegex.MatchString(phone)
}

func IsValidNIK(nik string) bool {
	// Indonesian NIK validation (16 digits)
	if len(nik) != 16 {
		return false
	}

	for _, char := range nik {
		if !unicode.IsDigit(char) {
			return false
		}
	}

	// Additional NIK validation logic
	return validateNIKLogic(nik)
}

func validateNIKLogic(nik string) bool {
	// Basic NIK validation
	// First 6 digits: area code
	// Next 6 digits: birth date (DDMMYY)
	// Last 4 digits: sequence number

	if len(nik) != 16 {
		return false
	}

	// Extract birth date part (positions 6-11)
	birthDateStr := nik[6:12]
	day, _ := strconv.Atoi(birthDateStr[0:2])
	month, _ := strconv.Atoi(birthDateStr[2:4])
	year, _ := strconv.Atoi(birthDateStr[4:6])

	// Validate day (1-31 for male, 41-71 for female)
	if (day < 1 || day > 31) && (day < 41 || day > 71) {
		return false
	}

	// Validate month (1-12)
	if month < 1 || month > 12 {
		return false
	}

	// Validate year (reasonable range)
	currentYear := time.Now().Year() % 100
	if year > currentYear+10 {
		return false
	}

	return true
}

// Password utilities
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateStrongPassword(length int) string {
	if length < 8 {
		length = 8
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomInt(len(charset))]
	}
	return string(b)
}

func randomInt(max int) int {
	bytes := make([]byte, 1)
	rand.Read(bytes)
	return int(bytes[0]) % max
}

// Number utilities
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %s", FormatNumber(amount))
}

func FormatNumber(num float64) string {
	str := fmt.Sprintf("%.0f", num)
	return addThousandSeparator(str)
}

func addThousandSeparator(str string) string {
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []string
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}

	return strings.Join(result, ".")
}

func ParseCurrency(currency string) (float64, error) {
	// Remove currency symbol and separators
	cleaned := strings.ReplaceAll(currency, "Rp", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", ".")
	cleaned = strings.TrimSpace(cleaned)

	return strconv.ParseFloat(cleaned, 64)
}

func CalculatePercentage(part, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (part / total) * 100
}

func CalculateGrowthRate(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}

// Date utilities
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func FormatDateTime(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func FormatDateIndonesia(date time.Time) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return fmt.Sprintf("%d %s %d", date.Day(), months[date.Month()-1], date.Year())
}

func ParseDateString(dateStr string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05",
	}

	for _, layout := range layouts {
		if date, err := time.Parse(layout, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func GetStartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func GetEndOfMonth(date time.Time) time.Time {
	return GetStartOfMonth(date).AddDate(0, 1, -1)
}

func GetStartOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
}

func GetEndOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), 12, 31, 23, 59, 59, 0, date.Location())
}

func DaysBetween(start, end time.Time) int {
	return int(math.Abs(end.Sub(start).Hours() / 24))
}

// Business logic utilities
func GenerateNIAK(koperasiID uint64, sequence uint64) string {
	// NIAK: Nomor Induk Anggota Koperasi
	// Format: KOOP + KoperasiID(4 digits) + Sequence(6 digits)
	return fmt.Sprintf("KOOP%04d%06d", koperasiID, sequence)
}

func GenerateNomorAnggota(koperasiID uint64, sequence uint64) string {
	// Format: A + KoperasiID(3 digits) + Sequence(4 digits)
	return fmt.Sprintf("A%03d%04d", koperasiID, sequence)
}

func GenerateNomorRekening(produkID uint64, sequence uint64) string {
	// Format: ProdukID(3 digits) + Sequence(6 digits)
	return fmt.Sprintf("%03d%06d", produkID, sequence)
}

func GenerateNomorTransaksi(koperasiID uint64, sequence uint64) string {
	// Format: TRX + Date(YYYYMMDD) + KoperasiID(3 digits) + Sequence(4 digits)
	today := time.Now().Format("20060102")
	return fmt.Sprintf("TRX%s%03d%04d", today, koperasiID, sequence)
}

func CalculateInterest(principal float64, rate float64, months int) float64 {
	// Simple interest calculation
	return principal * (rate / 100) * (float64(months) / 12)
}

func CalculateInstallment(principal float64, rate float64, months int) float64 {
	// Monthly installment calculation
	monthlyRate := rate / 100 / 12
	if monthlyRate == 0 {
		return principal / float64(months)
	}

	return principal * (monthlyRate * math.Pow(1+monthlyRate, float64(months))) /
		(math.Pow(1+monthlyRate, float64(months)) - 1)
}

// Array utilities
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

func ChunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// File utilities
func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return strings.ToLower(parts[len(parts)-1])
	}
	return ""
}

func IsImageFile(filename string) bool {
	imageExts := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp"}
	ext := GetFileExtension(filename)
	return Contains(imageExts, ext)
}

func IsDocumentFile(filename string) bool {
	docExts := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt"}
	ext := GetFileExtension(filename)
	return Contains(docExts, ext)
}

func GenerateFileName(originalName string, prefix string) string {
	ext := GetFileExtension(originalName)
	timestamp := time.Now().Format("20060102150405")
	random := GenerateRandomString(6)

	if ext != "" {
		return fmt.Sprintf("%s_%s_%s.%s", prefix, timestamp, random, ext)
	}
	return fmt.Sprintf("%s_%s_%s", prefix, timestamp, random)
}

// Pagination utilities
type PaginationInfo struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrevious  bool `json:"has_previous"`
}

func CalculatePagination(page, limit, totalItems int) PaginationInfo {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationInfo{
		Page:        page,
		Limit:       limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}

func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

// Response utilities
type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *PaginationInfo `json:"pagination,omitempty"`
	Meta       map[string]interface{} `json:"meta,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessResponseWithPagination(message string, data interface{}, pagination PaginationInfo) APIResponse {
	return APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
	}
}

func ErrorResponse(message string, err error) APIResponse {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return APIResponse{
		Success: false,
		Message: message,
		Error:   errorMsg,
	}
}

// Configuration utilities
func GetEnvOrDefault(key, defaultValue string) string {
	// This would typically use os.Getenv
	// For now, return default
	return defaultValue
}

func ParseConfigInt(value string, defaultValue int) int {
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return defaultValue
}

func ParseConfigFloat(value string, defaultValue float64) float64 {
	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return parsed
	}
	return defaultValue
}

func ParseConfigBool(value string, defaultValue bool) bool {
	if parsed, err := strconv.ParseBool(value); err == nil {
		return parsed
	}
	return defaultValue
}