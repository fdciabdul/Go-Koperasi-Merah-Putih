package api

import (
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"koperasi-merah-putih/tests/helpers"
)

type ComprehensiveTestSuite struct {
	helpers.BaseTestSuite
}

func TestComprehensiveTestSuite(t *testing.T) {
	suite.Run(t, new(ComprehensiveTestSuite))
}

// =======================
// WILAYAH ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestWilayahEndpoints() {
	// Test GET /api/v1/wilayah/provinsi
	s.MockDB.ExpectQuery("SELECT \\* FROM \"wilayah_provinsi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama"}).
			AddRow(11, "11", "Aceh").AddRow(12, "12", "Sumatera Utara"))
	w := s.MakeRequest("GET", "/api/v1/wilayah/provinsi", nil, "")
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/wilayah/provinsi/:provinsi_id/kabupaten
	s.MockDB.ExpectQuery("SELECT \\* FROM \"wilayah_kabupaten\" WHERE provinsi_id = \\$1").
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama", "provinsi_id"}).
			AddRow(1101, "1101", "Aceh Selatan", 11))
	w = s.MakeRequest("GET", "/api/v1/wilayah/provinsi/11/kabupaten", nil, "")
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/wilayah/kabupaten/:kabupaten_id/kecamatan
	s.MockDB.ExpectQuery("SELECT \\* FROM \"wilayah_kecamatan\" WHERE kabupaten_id = \\$1").
		WithArgs(1101).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama", "kabupaten_id"}).
			AddRow(110101, "110101", "Trumon", 1101))
	w = s.MakeRequest("GET", "/api/v1/wilayah/kabupaten/1101/kecamatan", nil, "")
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/wilayah/kecamatan/:kecamatan_id/kelurahan
	s.MockDB.ExpectQuery("SELECT \\* FROM \"wilayah_kelurahan\" WHERE kecamatan_id = \\$1").
		WithArgs(110101).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama", "kecamatan_id"}).
			AddRow(1101011001, "1101011001", "Trumon Timur", 110101))
	w = s.MakeRequest("GET", "/api/v1/wilayah/kecamatan/110101/kelurahan", nil, "")
	s.Equal(http.StatusOK, w.Code)
}

// =======================
// PRODUK ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestProdukEndpoints() {
	// Test POST /api/v1/produk/kategori
	kategoriReq := map[string]interface{}{
		"kode":      "MKN",
		"nama":      "Makanan",
		"deskripsi": "Kategori produk makanan",
		"icon":      "food",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"kategori_produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w := s.MakeRequest("POST", "/api/v1/produk/kategori", kategoriReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/produk/kategori
	s.MockDB.ExpectQuery("SELECT \\* FROM \"kategori_produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama"}).
			AddRow(1, "MKN", "Makanan"))
	w = s.MakeRequest("GET", "/api/v1/produk/kategori", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test POST /api/v1/produk/satuan
	satuanReq := map[string]interface{}{
		"kode": "PCS",
		"nama": "Pieces",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"satuan_produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/produk/satuan", satuanReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/produk/supplier
	supplierReq := map[string]interface{}{
		"koperasi_id":   1,
		"kode":          "SUP001",
		"nama":          "PT Supplier Test",
		"kontak_person": "Budi Supplier",
		"telepon":       "0211234567",
		"email":         "supplier@test.com",
		"alamat":        "Jl. Supplier No. 123",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"supplier\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/produk/supplier", supplierReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/produk (Create Product)
	produkReq := map[string]interface{}{
		"koperasi_id":        1,
		"kategori_produk_id": 1,
		"satuan_produk_id":   1,
		"kode_produk":        "PRD001",
		"nama_produk":        "Test Product",
		"deskripsi":          "Produk untuk testing",
		"harga_beli":         10000,
		"harga_jual":         15000,
		"stok_minimal":       10,
		"stok_current":       50,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"produk\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/produk", produkReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/produk/:koperasi_id
	s.MockDB.ExpectQuery("SELECT \\* FROM \"produk\" WHERE koperasi_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_produk", "nama_produk"}).
			AddRow(1, "PRD001", "Test Product"))
	w = s.MakeRequest("GET", "/api/v1/produk/1", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test POST /api/v1/produk/purchase-order
	poReq := map[string]interface{}{
		"koperasi_id":   1,
		"supplier_id":   1,
		"nomor_po":      "PO-TEST-001",
		"tanggal_po":    "2024-01-01",
		"items": []map[string]interface{}{
			{
				"produk_id":    1,
				"qty":          10,
				"harga_satuan": 10000,
			},
		},
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"purchase_order\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/produk/purchase-order", poReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/produk/penjualan
	penjualanReq := map[string]interface{}{
		"koperasi_id":       1,
		"anggota_id":        1,
		"nomor_transaksi":   "TXN-TEST-001",
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
	s.Equal(http.StatusCreated, w.Code)
}

// =======================
// FINANCIAL ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestFinancialEndpoints() {
	// Test POST /api/v1/financial/coa/akun
	akunReq := map[string]interface{}{
		"koperasi_id":   1,
		"kode_akun":     "1-1001",
		"nama_akun":     "Kas",
		"kategori_id":   1,
		"saldo_normal":  "debit",
		"is_kas":        true,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"coa_akun\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w := s.MakeRequest("POST", "/api/v1/financial/coa/akun", akunReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/financial/:koperasi_id/coa/akun
	s.MockDB.ExpectQuery("SELECT \\* FROM \"coa_akun\" WHERE koperasi_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_akun", "nama_akun"}).
			AddRow(1, "1-1001", "Kas"))
	w = s.MakeRequest("GET", "/api/v1/financial/1/coa/akun", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test POST /api/v1/financial/jurnal
	jurnalReq := map[string]interface{}{
		"koperasi_id":        1,
		"nomor_jurnal":       "JU-TEST-001",
		"tanggal_transaksi":  "2024-01-01",
		"keterangan":         "Jurnal test",
		"details": []map[string]interface{}{
			{
				"akun_id":    1,
				"keterangan": "Debit kas",
				"debit":      100000,
				"kredit":     0,
			},
			{
				"akun_id":    2,
				"keterangan": "Kredit modal",
				"debit":      0,
				"kredit":     100000,
			},
		},
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"jurnal_umum\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/financial/jurnal", jurnalReq, s.Token)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/financial/:koperasi_id/neraca-saldo
	s.MockDB.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"akun", "debit", "kredit"}).
			AddRow("Kas", 100000, 0).AddRow("Modal", 0, 100000))
	w = s.MakeRequest("GET", "/api/v1/financial/1/neraca-saldo?periode=2024-01", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/financial/:koperasi_id/laba-rugi
	s.MockDB.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"akun", "nominal"}).
			AddRow("Pendapatan", 500000).AddRow("Beban", 300000))
	w = s.MakeRequest("GET", "/api/v1/financial/1/laba-rugi?periode=2024-01", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)
}

// =======================
// SIMPAN PINJAM ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestSimpanPinjamEndpoints() {
	// Test POST /api/v1/simpan-pinjam/produk
	produkReq := map[string]interface{}{
		"koperasi_id":       1,
		"kode_produk":       "SP001",
		"nama_produk":       "Simpanan Sukarela",
		"jenis":             "simpanan",
		"kategori":          "sukarela",
		"bunga_simpanan":    3.5,
		"minimal_saldo":     50000,
		"syarat_ketentuan":  "Minimal setoran 50.000",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"produk_simpan_pinjam\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w := s.MakeRequest("POST", "/api/v1/simpan-pinjam/produk", produkReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/simpan-pinjam/rekening
	rekeningReq := map[string]interface{}{
		"koperasi_id":    1,
		"anggota_id":     1,
		"produk_id":      1,
		"nomor_rekening": "SV-001-0001",
		"saldo_simpanan": 100000,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"rekening_simpan_pinjam\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/simpan-pinjam/rekening", rekeningReq, s.Token)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/simpan-pinjam/transaksi
	transaksiReq := map[string]interface{}{
		"rekening_id":       1,
		"jenis_transaksi":   "setor",
		"jumlah":            50000,
		"keterangan":        "Setoran bulanan",
		"tanggal_transaksi": "2024-01-01",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"transaksi_simpan_pinjam\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/simpan-pinjam/transaksi", transaksiReq, s.Token)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/simpan-pinjam/:koperasi_id/statistik
	s.MockDB.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_simpanan", "total_pinjaman"}).
			AddRow(1000000, 500000))
	w = s.MakeRequest("GET", "/api/v1/simpan-pinjam/1/statistik", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)
}

// =======================
// PPOB ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestPPOBEndpoints() {
	// Test GET /api/v1/ppob/kategoris
	s.MockDB.ExpectQuery("SELECT \\* FROM \"ppob_kategori\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama"}).
			AddRow(1, "PLN", "PLN Token"))
	w := s.MakeRequest("GET", "/api/v1/ppob/kategoris", nil, "")
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/ppob/kategoris/:kategori_id/produks
	s.MockDB.ExpectQuery("SELECT \\* FROM \"ppob_produk\" WHERE kategori_id = \\$1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode_produk", "nama_produk"}).
			AddRow(1, "PLN20", "PLN Token 20.000"))
	w = s.MakeRequest("GET", "/api/v1/ppob/kategoris/1/produks", nil, "")
	s.Equal(http.StatusOK, w.Code)

	// Test POST /api/v1/ppob/transactions
	transaksiReq := map[string]interface{}{
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
	w = s.MakeRequest("POST", "/api/v1/ppob/transactions", transaksiReq, s.Token)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/ppob/settlements
	settlementReq := map[string]interface{}{
		"koperasi_id":        1,
		"periode_dari":       "2024-01-01",
		"periode_sampai":     "2024-01-31",
		"transaction_ids":    []int{1, 2, 3},
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"ppob_settlement\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/ppob/settlements", settlementReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)
}

// =======================
// KLINIK ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestKlinikEndpoints() {
	// Test POST /api/v1/klinik/pasien
	pasienReq := map[string]interface{}{
		"koperasi_id":      1,
		"nomor_rm":         "RM001",
		"nik":              "1234567890123458",
		"nama_lengkap":     "Pasien Test",
		"jenis_kelamin":    "P",
		"tempat_lahir":     "Jakarta",
		"tanggal_lahir":    "1990-01-01",
		"alamat":           "Jl. Pasien No. 1",
		"telepon":          "081234567890",
		"email":            "pasien@test.com",
		"golongan_darah":   "A",
		"anggota_id":       1,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"klinik_pasien\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w := s.MakeRequest("POST", "/api/v1/klinik/pasien", pasienReq, s.Token)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/klinik/tenaga-medis
	tenagaMedisReq := map[string]interface{}{
		"koperasi_id":      1,
		"nip":              "DOK001",
		"nama_lengkap":     "Dr. Test",
		"jenis_kelamin":    "L",
		"spesialisasi":     "Umum",
		"no_str":           "STR123456",
		"telepon":          "081234567891",
		"email":            "dokter@test.com",
		"tarif_konsultasi": 100000,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"klinik_tenaga_medis\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/klinik/tenaga-medis", tenagaMedisReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/klinik/kunjungan
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
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/klinik/obat
	obatReq := map[string]interface{}{
		"koperasi_id":   1,
		"kode_obat":     "OBT001",
		"nama_obat":     "Paracetamol",
		"jenis_obat":    "tablet",
		"satuan":        "strip",
		"harga_beli":    5000,
		"harga_jual":    7500,
		"stok_current":  100,
		"stok_minimal":  10,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"klinik_obat\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/klinik/obat", obatReq, s.AdminToken)
	s.Equal(http.StatusCreated, w.Code)
}

// =======================
// MASTER DATA ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestMasterDataEndpoints() {
	// Test POST /api/v1/master-data/kbli
	kbliReq := map[string]interface{}{
		"kode":      "47911",
		"nama":      "Perdagangan Eceran Online",
		"kategori":  "Perdagangan",
		"deskripsi": "Kegiatan perdagangan eceran melalui internet",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"kbli\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w := s.MakeRequest("POST", "/api/v1/master-data/kbli", kbliReq, s.SuperAdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test GET /api/v1/master-data/kbli
	s.MockDB.ExpectQuery("SELECT \\* FROM \"kbli\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "kode", "nama"}).
			AddRow(1, "47911", "Perdagangan Eceran Online"))
	w = s.MakeRequest("GET", "/api/v1/master-data/kbli", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test POST /api/v1/master-data/jenis-koperasi
	jenisReq := map[string]interface{}{
		"kode":      "KP",
		"nama":      "Koperasi Primer",
		"deskripsi": "Koperasi yang beranggotakan orang-seorang",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"jenis_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/master-data/jenis-koperasi", jenisReq, s.SuperAdminToken)
	s.Equal(http.StatusCreated, w.Code)

	// Test POST /api/v1/master-data/bentuk-koperasi
	bentukReq := map[string]interface{}{
		"kode":      "UNIT",
		"nama":      "Unit Koperasi",
		"deskripsi": "Unit koperasi",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectQuery("INSERT INTO \"bentuk_koperasi\"").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("POST", "/api/v1/master-data/bentuk-koperasi", bentukReq, s.SuperAdminToken)
	s.Equal(http.StatusCreated, w.Code)
}

// =======================
// ADMIN ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestAdminEndpoints() {
	// Test GET /api/v1/admin/sequences
	s.MockDB.ExpectQuery("SELECT \\* FROM \"sequence_numbers\"").
		WillReturnRows(sqlmock.NewRows([]string{"id", "sequence_name", "current_number"}).
			AddRow(1, "global", 100).AddRow(2, "koperasi", 5))
	w := s.MakeRequest("GET", "/api/v1/admin/sequences", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)

	// Test PUT /api/v1/admin/sequences/update-value
	updateReq := map[string]interface{}{
		"sequence_name":  "global",
		"current_number": 200,
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectExec("UPDATE \"sequence_numbers\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("PUT", "/api/v1/admin/sequences/update-value", updateReq, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)

	// Test PUT /api/v1/admin/sequences/reset
	resetReq := map[string]interface{}{
		"sequence_name": "global",
	}
	s.MockDB.ExpectBegin()
	s.MockDB.ExpectExec("UPDATE \"sequence_numbers\"").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.MockDB.ExpectCommit()
	w = s.MakeRequest("PUT", "/api/v1/admin/sequences/reset", resetReq, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)
}

// =======================
// REPORTING ENDPOINTS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestReportingEndpoints() {
	// Test GET /api/v1/reports/:koperasi_id/dashboard
	s.MockDB.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"total_anggota", "total_transaksi"}).
			AddRow(100, 500))
	w := s.MakeRequest("GET", "/api/v1/reports/1/dashboard", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/reports/:koperasi_id/quick-summary
	s.MockDB.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"summary"}).AddRow("OK"))
	w = s.MakeRequest("GET", "/api/v1/reports/1/quick-summary", nil, s.Token)
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/reports/:koperasi_id/analytics/revenue
	w = s.MakeRequest("GET", "/api/v1/reports/1/analytics/revenue?periode=2024-01", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/reports/:koperasi_id/sales
	w = s.MakeRequest("GET", "/api/v1/reports/1/sales?start_date=2024-01-01&end_date=2024-01-31", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)

	// Test GET /api/v1/reports/:koperasi_id/export/sales
	w = s.MakeRequest("GET", "/api/v1/reports/1/export/sales?start_date=2024-01-01&end_date=2024-01-31&format=csv", nil, s.AdminToken)
	s.Equal(http.StatusOK, w.Code)
}

// =======================
// UNAUTHORIZED ACCESS TESTS
// =======================

func (s *ComprehensiveTestSuite) TestUnauthorizedAccess() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/koperasi"},
		{"PUT", "/api/v1/koperasi/1"},
		{"DELETE", "/api/v1/koperasi/1"},
		{"POST", "/api/v1/produk/kategori"},
		{"POST", "/api/v1/financial/coa/akun"},
		{"POST", "/api/v1/simpan-pinjam/produk"},
		{"POST", "/api/v1/klinik/tenaga-medis"},
		{"POST", "/api/v1/master-data/kbli"},
		{"GET", "/api/v1/admin/sequences"},
		{"GET", "/api/v1/reports/1/dashboard"},
	}

	for _, endpoint := range endpoints {
		w := s.MakeRequest(endpoint.method, endpoint.path, nil, "")
		s.Equal(http.StatusUnauthorized, w.Code, "Expected unauthorized for %s %s", endpoint.method, endpoint.path)
	}
}

// =======================
// INVALID INPUT TESTS
// =======================

func (s *ComprehensiveTestSuite) TestInvalidInputValidation() {
	// Test invalid JSON
	w := s.MakeRequest("POST", "/api/v1/koperasi", "invalid-json", s.SuperAdminToken)
	s.Equal(http.StatusBadRequest, w.Code)

	// Test missing required fields
	emptyReq := map[string]interface{}{}
	w = s.MakeRequest("POST", "/api/v1/koperasi", emptyReq, s.SuperAdminToken)
	s.Equal(http.StatusBadRequest, w.Code)

	// Test invalid field types
	invalidReq := map[string]interface{}{
		"nik": "not-a-number",
	}
	w = s.MakeRequest("POST", "/api/v1/koperasi", invalidReq, s.SuperAdminToken)
	s.Equal(http.StatusBadRequest, w.Code)
}

// =======================
// RBAC TESTS
// =======================

func (s *ComprehensiveTestSuite) TestRoleBasedAccessControl() {
	// Test super admin only endpoints with admin token
	w := s.MakeRequest("POST", "/api/v1/koperasi", map[string]interface{}{"nama": "test"}, s.AdminToken)
	s.Equal(http.StatusForbidden, w.Code)

	// Test admin only endpoints with user token
	w = s.MakeRequest("POST", "/api/v1/produk/kategori", map[string]interface{}{"nama": "test"}, s.Token)
	s.Equal(http.StatusForbidden, w.Code)

	// Test financial access endpoints
	w = s.MakeRequest("POST", "/api/v1/financial/jurnal", map[string]interface{}{"keterangan": "test"}, s.Token)
	s.Equal(http.StatusForbidden, w.Code)
}