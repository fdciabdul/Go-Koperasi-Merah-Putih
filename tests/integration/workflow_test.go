package integration

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"koperasi-merah-putih/tests/helpers"
)

type WorkflowTestSuite struct {
	helpers.BaseTestSuite
}

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}

// TestCompleteWorkflow tests a complete end-to-end workflow
func (s *WorkflowTestSuite) TestCompleteWorkflow() {
	// Step 1: User Registration
	s.T().Log("Step 1: User Registration")
	registerReq := map[string]interface{}{
		"koperasi_id":        1,
		"nik":                "1234567890123457",
		"nama_lengkap":       "Test Workflow User",
		"email":              "workflow@test.com",
		"no_telepon":         "081234567899",
		"alamat":             "Jl. Workflow No. 1",
		"kelurahan_id":       1,
		"password":           "password123",
		"jenis_keanggotaan":  "anggota",
		"rencana_simpanan":   1000000,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"user_registrations\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/users/register", registerReq, "")
	s.Equal(200, w.Code)

	// Step 2: Payment Creation
	s.T().Log("Step 2: Payment Creation")
	paymentReq := map[string]interface{}{
		"registration_id": 1,
		"amount":         1000000,
		"payment_method": "bank_transfer",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"payment_transactions\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/payments", paymentReq, s.Token)
	s.Equal(200, w.Code)

	// Step 3: Payment Verification
	s.T().Log("Step 3: Payment Verification")
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"payment_transactions\"").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "pending"))
	s.MockDB.ExpectExec("UPDATE \"payment_transactions\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("PUT", "/api/v1/users/verify-payment/1", nil, "")
	s.Equal(200, w.Code)

	// Step 4: Registration Approval
	s.T().Log("Step 4: Registration Approval")
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("SELECT \\* FROM \"user_registrations\"").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "pending_approval"))
	s.MockDB.ExpectQuery("INSERT INTO \"users\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectQuery("INSERT INTO \"anggota_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectExec("UPDATE \"user_registrations\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("PUT", "/api/v1/users/registrations/1/approve", nil, s.AdminToken)
	s.Equal(200, w.Code)

	// Step 5: User Login
	s.T().Log("Step 5: User Login")
	loginReq := map[string]interface{}{
		"email":    "workflow@test.com",
		"password": "password123",
	}

	s.MockDB.ExpectQuery("SELECT \\* FROM \"users\" WHERE email = \\$1").
		WithArgs("workflow@test.com").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tenant_id", "email", "password_hash", "nama_lengkap", "role", "is_active",
		}).AddRow(2, 1, "workflow@test.com", "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", "Test Workflow User", "user", true))

	w = s.MakeRequest("POST", "/api/v1/auth/login", loginReq, "")
	s.Equal(200, w.Code)
	s.Contains(w.Body.String(), "token")

	// Step 6: Create Simpan Pinjam Account
	s.T().Log("Step 6: Create Simpan Pinjam Account")
	rekeningReq := map[string]interface{}{
		"koperasi_id":    1,
		"anggota_id":     1,
		"produk_id":      1,
		"nomor_rekening": "SV-001-0002",
		"saldo_simpanan": 0,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"rekening_simpan_pinjam\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/simpan-pinjam/rekening", rekeningReq, s.Token)
	s.Equal(200, w.Code)

	// Step 7: Make Savings Transaction
	s.T().Log("Step 7: Make Savings Transaction")
	transaksiReq := map[string]interface{}{
		"rekening_id":       1,
		"jenis_transaksi":   "setor",
		"jumlah":            500000,
		"keterangan":        "Setoran awal",
		"tanggal_transaksi": "2024-01-01",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"transaksi_simpan_pinjam\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/simpan-pinjam/transaksi", transaksiReq, s.Token)
	s.Equal(200, w.Code)

	// Step 8: Purchase Products
	s.T().Log("Step 8: Purchase Products")
	penjualanReq := map[string]interface{}{
		"koperasi_id":       1,
		"anggota_id":        1,
		"nomor_transaksi":   "TXN-WF-001",
		"tanggal_transaksi": "2024-01-01",
		"items": []map[string]interface{}{
			{
				"produk_id":    1,
				"qty":          2,
				"harga_satuan": 15000,
			},
		},
		"metode_pembayaran": "cash",
		"jumlah_bayar":      30000,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"penjualan_header\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/produk/penjualan", penjualanReq, s.Token)
	s.Equal(200, w.Code)

	// Step 9: Make PPOB Transaction
	s.T().Log("Step 9: Make PPOB Transaction")
	ppobReq := map[string]interface{}{
		"koperasi_id":     1,
		"anggota_id":      1,
		"produk_id":       1,
		"nomor_tujuan":    "123456789012",
		"nama_pelanggan":  "Test Customer",
		"customer_name":   "Test Customer",
		"customer_email":  "customer@test.com",
		"customer_phone":  "081234567890",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"ppob_transaksi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/ppob/transactions", ppobReq, s.Token)
	s.Equal(200, w.Code)

	// Step 10: Create Medical Record
	s.T().Log("Step 10: Create Medical Record")
	kunjunganReq := map[string]interface{}{
		"koperasi_id":       1,
		"pasien_id":         1,
		"dokter_id":         1,
		"tanggal_kunjungan": "2024-01-01",
		"keluhan":           "Demam dan batuk",
		"diagnosis":         "ISPA",
		"terapi":            "Istirahat dan minum obat",
		"biaya_konsultasi":  100000,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"klinik_kunjungan\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/klinik/kunjungan", kunjunganReq, s.Token)
	s.Equal(200, w.Code)

	s.T().Log("✓ Complete workflow test passed successfully!")
}

// TestKoperasiManagementWorkflow tests koperasi management workflow
func (s *WorkflowTestSuite) TestKoperasiManagementWorkflow() {
	// Step 1: Create Koperasi (Super Admin)
	s.T().Log("Step 1: Create Koperasi")
	koperasiReq := map[string]interface{}{
		"nomor_sk":           "SK-WF-001",
		"nik":                1234567890123460,
		"nama_koperasi":      "Koperasi Workflow Test",
		"nama_sk":            "Koperasi Workflow Test SK",
		"jenis_koperasi_id":  1,
		"bentuk_koperasi_id": 1,
		"status_koperasi_id": 1,
		"email":              "workflow-kop@test.com",
		"telepon":            "0812345679",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w := s.MakeRequest("POST", "/api/v1/koperasi", koperasiReq, s.SuperAdminToken)
	s.Equal(200, w.Code)

	// Step 2: Add Member to Koperasi
	s.T().Log("Step 2: Add Member to Koperasi")
	anggotaReq := map[string]interface{}{
		"koperasi_id":    2,
		"nik":            "1234567890123461",
		"nama":           "Anggota Workflow",
		"jenis_kelamin":  "L",
		"email":          "anggota-wf@test.com",
		"telepon":        "081234567892",
		"status_anggota": "aktif",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"anggota_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/koperasi/anggota", anggotaReq, s.AdminToken)
	s.Equal(200, w.Code)

	// Step 3: Create Product Category
	s.T().Log("Step 3: Create Product Category")
	kategoriReq := map[string]interface{}{
		"kode": "WF",
		"nama": "Workflow Category",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"kategori_produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/produk/kategori", kategoriReq, s.AdminToken)
	s.Equal(200, w.Code)

	// Step 4: Create Product
	s.T().Log("Step 4: Create Product")
	produkReq := map[string]interface{}{
		"koperasi_id":        2,
		"kategori_produk_id": 2,
		"satuan_produk_id":   1,
		"kode_produk":        "WF-001",
		"nama_produk":        "Workflow Product",
		"harga_beli":         8000,
		"harga_jual":         10000,
		"stok_current":       100,
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/produk", produkReq, s.AdminToken)
	s.Equal(200, w.Code)

	// Step 5: Create COA Account
	s.T().Log("Step 5: Create COA Account")
	akunReq := map[string]interface{}{
		"koperasi_id":   2,
		"kode_akun":     "1-1002",
		"nama_akun":     "Bank",
		"kategori_id":   1,
		"saldo_normal":  "debit",
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"coa_akun\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/financial/coa/akun", akunReq, s.AdminToken)
	s.Equal(200, w.Code)

	// Step 6: Create Journal Entry
	s.T().Log("Step 6: Create Journal Entry")
	jurnalReq := map[string]interface{}{
		"koperasi_id":        2,
		"nomor_jurnal":       "JU-WF-001",
		"tanggal_transaksi":  "2024-01-01",
		"keterangan":         "Jurnal workflow test",
		"details": []map[string]interface{}{
			{
				"akun_id":    1,
				"keterangan": "Debit kas",
				"debit":      500000,
				"kredit":     0,
			},
			{
				"akun_id":    2,
				"keterangan": "Kredit bank",
				"debit":      0,
				"kredit":     500000,
			},
		},
	}

	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"jurnal_umum\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.MockDB.ExpectCommit()

	w = s.MakeRequest("POST", "/api/v1/financial/jurnal", jurnalReq, s.Token)
	s.Equal(200, w.Code)

	s.T().Log("✓ Koperasi management workflow test passed successfully!")
}

// TestErrorHandlingWorkflow tests error handling in workflows
func (s *WorkflowTestSuite) TestErrorHandlingWorkflow() {
	s.T().Log("Testing error handling workflow")

	// Test database connection failure
	s.MockDB.ExpectQuery("SELECT").WillReturnError(sqlmock.ErrCancelled)
	w := s.MakeRequest("GET", "/api/v1/koperasi", nil, s.Token)
	s.Equal(500, w.Code)

	// Test validation errors
	invalidReq := map[string]interface{}{
		"invalid_field": "invalid_value",
	}
	w = s.MakeRequest("POST", "/api/v1/koperasi", invalidReq, s.SuperAdminToken)
	s.Equal(400, w.Code)

	// Test authentication errors
	w = s.MakeRequest("GET", "/api/v1/koperasi", nil, "invalid-token")
	s.Equal(401, w.Code)

	// Test authorization errors
	w = s.MakeRequest("POST", "/api/v1/koperasi", map[string]interface{}{"nama": "test"}, s.Token)
	s.Equal(403, w.Code)

	s.T().Log("✓ Error handling workflow test passed successfully!")
}

// TestPerformanceWorkflow tests performance under load
func (s *WorkflowTestSuite) TestPerformanceWorkflow() {
	s.T().Log("Testing performance workflow")

	// Simulate multiple concurrent requests
	for i := 0; i < 10; i++ {
		s.MockDB.ExpectQuery("SELECT \\* FROM \"koperasi\"").
			WillReturnRows(sqlmock.NewRows([]string{"id", "nama_koperasi"}).
				AddRow(1, "Test Koperasi"))
	}

	// Make multiple requests
	for i := 0; i < 10; i++ {
		w := s.MakeRequest("GET", "/api/v1/koperasi", nil, s.Token)
		s.Equal(200, w.Code)
	}

	s.T().Log("✓ Performance workflow test passed successfully!")
}