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

type KoperasiTestSuite struct {
	helpers.BaseTestSuite
}

func TestKoperasiTestSuite(t *testing.T) {
	suite.Run(t, new(KoperasiTestSuite))
}

// Test Create Koperasi (POST /api/v1/koperasi)
func (s *KoperasiTestSuite) TestCreateKoperasi_Success() {
	koperasiReq := map[string]interface{}{
		"nomor_sk":           "SK-TEST-001",
		"nik":                1234567890123456,
		"nama_koperasi":      "Koperasi Test Baru",
		"nama_sk":            "Koperasi Test SK",
		"jenis_koperasi_id":  1,
		"bentuk_koperasi_id": 1,
		"status_koperasi_id": 1,
		"provinsi_id":        11,
		"kabupaten_id":       1101,
		"kecamatan_id":       110101,
		"kelurahan_id":       1101011001,
		"alamat":             "Jl. Test No. 123",
		"rt":                 "001",
		"rw":                 "002",
		"kode_pos":           "12345",
		"email":              "koperasi@test.com",
		"telepon":            "0812345678",
		"website":            "www.test.com",
		"tanggal_berdiri":    "2024-01-01",
		"tanggal_sk":         "2024-01-15",
		"tanggal_pengesahan": "2024-02-01",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, s.SuperAdminToken)

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Koperasi created successfully")
}

func (s *KoperasiTestSuite) TestCreateKoperasi_Unauthorized() {
	koperasiReq := map[string]interface{}{
		"nama_koperasi": "Test",
	}

	// Test without token
	w := s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, "")
	s.Equal(http.StatusUnauthorized, w.Code)

	// Test with regular user token
	w = s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, s.Token)
	s.Equal(http.StatusForbidden, w.Code)

	// Test with admin token (not super admin)
	w = s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, s.AdminToken)
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *KoperasiTestSuite) TestCreateKoperasi_DuplicateSK() {
	koperasiReq := map[string]interface{}{
		"nomor_sk":      "SK-EXISTING-001",
		"nik":           1234567890123456,
		"nama_koperasi": "Koperasi Test",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"koperasi\"").
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	s.MockDB.ExpectRollback()

	w := s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, s.SuperAdminToken)

	s.Equal(http.StatusInternalServerError, w.Code)
}

// Test Get Koperasi List (GET /api/v1/koperasi)
func (s *KoperasiTestSuite) TestGetKoperasiList_Success() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "nomor_sk", "nama_koperasi", "email",
		}).AddRow(1, 1, "SK-001", "Koperasi Test 1", "test1@koperasi.com").
			AddRow(2, 1, "SK-002", "Koperasi Test 2", "test2@koperasi.com"))

	w := s.MakeRequest("GET", "/api/v1/koperasi", nil, s.Token)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Koperasi Test 1")
	s.Contains(w.Body.String(), "Koperasi Test 2")
}

func (s *KoperasiTestSuite) TestGetKoperasiList_WithPagination() {
	w := s.MakeRequest("GET", "/api/v1/koperasi?page=1&limit=10", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)
}

func (s *KoperasiTestSuite) TestGetKoperasiList_WithFilter() {
	w := s.MakeRequest("GET", "/api/v1/koperasi?status=aktif&provinsi_id=11", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)
}

// Test Get Koperasi by ID (GET /api/v1/koperasi/:id)
func (s *KoperasiTestSuite) TestGetKoperasi_Success() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"koperasi\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "nomor_sk", "nama_koperasi", "email",
		}).AddRow(1, 1, "SK-001", "Koperasi Test", "test@koperasi.com"))

	w := s.MakeRequest("GET", "/api/v1/koperasi/1", nil, s.Token)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Koperasi Test")
}

func (s *KoperasiTestSuite) TestGetKoperasi_NotFound() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"koperasi\" WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	w := s.MakeRequest("GET", "/api/v1/koperasi/999", nil, s.Token)

	s.Equal(http.StatusNotFound, w.Code)
}

// Test Update Koperasi (PUT /api/v1/koperasi/:id)
func (s *KoperasiTestSuite) TestUpdateKoperasi_Success() {
	updateReq := map[string]interface{}{
		"nama_koperasi": "Updated Koperasi Name",
		"email":         "updated@koperasi.com",
		"telepon":       "081999999999",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"koperasi\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "nama_koperasi",
		}).AddRow(1, "Old Name"))

	s.MockDB.ExpectExec("UPDATE \"koperasi\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("PUT", "/api/v1/koperasi/1", updateReq, s.AdminToken)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Koperasi updated successfully")
}

func (s *KoperasiTestSuite) TestUpdateKoperasi_Unauthorized() {
	updateReq := map[string]interface{}{
		"nama_koperasi": "Updated Name",
	}

	// Test with regular user
	w := s.MakeRequest("PUT", "/api/v1/koperasi/1", updateReq, s.Token)
	s.Equal(http.StatusForbidden, w.Code)
}

// Test Delete Koperasi (DELETE /api/v1/koperasi/:id)
func (s *KoperasiTestSuite) TestDeleteKoperasi_Success() {
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectExec("DELETE FROM \"koperasi\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("DELETE", "/api/v1/koperasi/1", nil, s.SuperAdminToken)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Koperasi deleted successfully")
}

func (s *KoperasiTestSuite) TestDeleteKoperasi_Unauthorized() {
	// Test with admin (not super admin)
	w := s.MakeRequest("DELETE", "/api/v1/koperasi/1", nil, s.AdminToken)
	s.Equal(http.StatusForbidden, w.Code)
}

// Test Create Anggota (POST /api/v1/koperasi/anggota)
func (s *KoperasiTestSuite) TestCreateAnggota_Success() {
	anggotaReq := map[string]interface{}{
		"koperasi_id":    1,
		"nik":            "1234567890123459",
		"nama":           "Anggota Baru",
		"jenis_kelamin":  "L",
		"tempat_lahir":   "Jakarta",
		"tanggal_lahir":  "1990-01-01",
		"alamat":         "Jl. Anggota No. 1",
		"rt":             "001",
		"rw":             "001",
		"kelurahan_id":   1101011001,
		"telepon":        "081234567890",
		"email":          "anggota@test.com",
		"posisi":         "anggota",
		"jabatan_id":     1,
		"tanggal_masuk":  "2024-01-01",
		"status_anggota": "aktif",
		"pekerjaan":      "Wiraswasta",
		"pendidikan":     "S1",
	}

	s.MockDB.ExpectBegin()

	// Mock sequence generation for NIAK
	s.MockDB.ExpectQuery("SELECT \\* FROM \"sequence_numbers\"").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "current_number",
		}).AddRow(1, 100))

	s.MockDB.ExpectExec("UPDATE \"sequence_numbers\"").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.MockDB.ExpectQuery("INSERT INTO \"anggota_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/koperasi/anggota", anggotaReq, s.AdminToken)

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "Anggota created successfully")
}

func (s *KoperasiTestSuite) TestCreateAnggota_DuplicateNIK() {
	anggotaReq := map[string]interface{}{
		"koperasi_id":   1,
		"nik":           "1234567890123456", // Existing NIK
		"nama":          "Anggota Duplicate",
		"jenis_kelamin": "P",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"anggota_koperasi\"").
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	s.MockDB.ExpectRollback()

	w := s.MakeRequest("POST", "/api/v1/koperasi/anggota", anggotaReq, s.AdminToken)

	s.Equal(http.StatusInternalServerError, w.Code)
}

// Test Get Anggota List (GET /api/v1/koperasi/:id/anggota)
func (s *KoperasiTestSuite) TestGetAnggotaList_Success() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"anggota_koperasi\" WHERE koperasi_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "koperasi_id", "niak", "nama", "status_anggota",
		}).AddRow(1, 1, "NIAK001", "Anggota 1", "aktif").
			AddRow(2, 1, "NIAK002", "Anggota 2", "aktif"))

	w := s.MakeRequest("GET", "/api/v1/koperasi/1/anggota", nil, s.Token)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Anggota 1")
	s.Contains(w.Body.String(), "Anggota 2")
}

func (s *KoperasiTestSuite) TestGetAnggotaList_WithFilter() {
	w := s.MakeRequest("GET", "/api/v1/koperasi/1/anggota?status=aktif&posisi=pengurus", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)
}

// Test Get Anggota by ID (GET /api/v1/koperasi/anggota/:id)
func (s *KoperasiTestSuite) TestGetAnggota_Success() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"anggota_koperasi\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "koperasi_id", "niak", "nama", "nik",
		}).AddRow(1, 1, "NIAK001", "Test Anggota", "1234567890123456"))

	w := s.MakeRequest("GET", "/api/v1/koperasi/anggota/1", nil, s.Token)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Test Anggota")
}

func (s *KoperasiTestSuite) TestGetAnggota_NotFound() {
	s.MockDB.ExpectQuery("SELECT \\* FROM \"anggota_koperasi\" WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	w := s.MakeRequest("GET", "/api/v1/koperasi/anggota/999", nil, s.Token)

	s.Equal(http.StatusNotFound, w.Code)
}

// Test Update Anggota Status (PUT /api/v1/koperasi/anggota/:id/status)
func (s *KoperasiTestSuite) TestUpdateAnggotaStatus_Success() {
	statusReq := map[string]interface{}{
		"status_anggota": "nonaktif",
		"tanggal_keluar": "2024-12-31",
		"alasan":         "Mengundurkan diri",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"anggota_koperasi\" WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "status_anggota",
		}).AddRow(1, "aktif"))

	s.MockDB.ExpectExec("UPDATE \"anggota_koperasi\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("PUT", "/api/v1/koperasi/anggota/1/status", statusReq, s.AdminToken)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Status updated successfully")
}

func (s *KoperasiTestSuite) TestUpdateAnggotaStatus_InvalidStatus() {
	statusReq := map[string]interface{}{
		"status_anggota": "invalid_status",
	}

	w := s.MakeRequest("PUT", "/api/v1/koperasi/anggota/1/status", statusReq, s.AdminToken)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *KoperasiTestSuite) TestUpdateAnggotaStatus_Unauthorized() {
	statusReq := map[string]interface{}{
		"status_anggota": "nonaktif",
	}

	w := s.MakeRequest("PUT", "/api/v1/koperasi/anggota/1/status", statusReq, s.Token)
	s.Equal(http.StatusForbidden, w.Code)
}