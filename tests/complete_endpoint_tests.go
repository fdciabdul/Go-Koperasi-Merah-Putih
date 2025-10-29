package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRealEndpointsWithData tests all endpoints with real data
func TestRealEndpointsWithData(t *testing.T) {
	baseURL := "http://localhost:8080"

	// Step 1: Test User Registration with Real Data
	t.Run("01_UserRegistration", func(t *testing.T) {
		registerData := map[string]interface{}{
			"koperasi_id":        1,
			"nik":               "3201012345678901",
			"nama_lengkap":      "Budi Santoso",
			"email":             "budi.santoso@email.com",
			"no_telepon":        "081234567890",
			"alamat":            "Jl. Merdeka No. 123, RT 01/RW 05",
			"kelurahan_id":      1101011001,
			"password":          "BudiSecure123!",
			"jenis_keanggotaan": "anggota",
			"rencana_simpanan":  2000000,
		}

		jsonData, _ := json.Marshal(registerData)
		resp, err := http.Post(baseURL+"/api/v1/users/register", "application/json", bytes.NewBuffer(jsonData))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Contains(t, response, "message")
		assert.Contains(t, response, "registration")

		t.Logf("Registration Response: %+v", response)
	})

	// Step 2: Test Login with Real Credentials
	var authToken string
	t.Run("02_Login", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "admin@demo.local",
			"password": "admin123",
		}

		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Contains(t, response, "message")
		assert.Equal(t, "Login successful", response["message"])
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "token")
		assert.Contains(t, data, "user")
		assert.Contains(t, data, "expires_at")

		authToken = data["token"].(string)
		t.Logf("Login successful, token obtained: %s...", authToken[:20])
	})

	// Step 3: Test Wilayah Endpoints (No auth needed)
	t.Run("03_WilayahEndpoints", func(t *testing.T) {
		// Test Get Provinces
		resp, err := http.Get(baseURL + "/api/v1/wilayah/provinsi")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var provinces []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&provinces)
		assert.Greater(t, len(provinces), 0)
		t.Logf("Found %d provinces", len(provinces))

		// Test Get Cities for a province
		resp, err = http.Get(baseURL + "/api/v1/wilayah/provinsi/11/kabupaten")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var cities []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&cities)
		t.Logf("Found %d cities in Aceh province", len(cities))
	})

	// Step 4: Test Koperasi Management with Authentication
	var koperasiID float64
	t.Run("04_KoperasiManagement", func(t *testing.T) {
		// Create Koperasi (Super Admin only)
		koperasiData := map[string]interface{}{
			"nomor_sk":           "001/SK/TEST/2024",
			"nik":                3201234567890123,
			"nama_koperasi":      "Koperasi Sejahtera Bersama",
			"nama_sk":            "Koperasi Sejahtera Bersama",
			"jenis_koperasi_id":  1,
			"bentuk_koperasi_id": 1,
			"status_koperasi_id": 1,
			"provinsi_id":        11,
			"kabupaten_id":       1101,
			"kecamatan_id":       110101,
			"kelurahan_id":       1101011001,
			"alamat":             "Jl. Koperasi Raya No. 45",
			"rt":                 "003",
			"rw":                 "007",
			"kode_pos":           "23711",
			"email":              "info@koperasisejahtera.com",
			"telepon":            "0651234567",
			"website":            "www.koperasisejahtera.com",
			"tanggal_berdiri":    "2020-01-15",
			"tanggal_sk":         "2020-02-01",
			"tanggal_pengesahan": "2020-03-01",
		}

		jsonData, _ := json.Marshal(koperasiData)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/koperasi", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		client := &http.Client{}
		resp, err := client.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		koperasi := response["koperasi"].(map[string]interface{})
		koperasiID = koperasi["id"].(float64)
		t.Logf("Koperasi created with ID: %.0f", koperasiID)

		// Test Get Koperasi List
		req, _ = http.NewRequest("GET", baseURL+"/api/v1/koperasi", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var koperasiList []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&koperasiList)
		assert.Greater(t, len(koperasiList), 0)
		t.Logf("Found %d koperasi in list", len(koperasiList))
	})

	// Step 5: Test Product Management
	t.Run("05_ProductManagement", func(t *testing.T) {
		client := &http.Client{}

		// Create Product Category
		categoryData := map[string]interface{}{
			"kode":      "FOOD",
			"nama":      "Makanan & Minuman",
			"deskripsi": "Kategori produk makanan dan minuman",
			"icon":      "utensils",
		}

		jsonData, _ := json.Marshal(categoryData)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/produk/kategori", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create Product Unit
		unitData := map[string]interface{}{
			"kode": "KG",
			"nama": "Kilogram",
		}

		jsonData, _ = json.Marshal(unitData)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/produk/satuan", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create Supplier
		supplierData := map[string]interface{}{
			"koperasi_id":     koperasiID,
			"kode":            "SUP001",
			"nama":            "CV Sumber Rejeki",
			"kontak_person":   "Pak Andi",
			"telepon":         "081987654321",
			"email":           "andi@sumberrejeki.com",
			"alamat":          "Jl. Industri No. 89, Medan",
			"provinsi_id":     12,
			"kabupaten_id":    1201,
			"no_rekening":     "1234567890",
			"nama_bank":       "Bank Mandiri",
			"atas_nama_bank":  "CV Sumber Rejeki",
			"npwp":            "12.345.678.9-012.000",
			"jenis_supplier":  "perusahaan",
			"term_pembayaran": 30,
		}

		jsonData, _ = json.Marshal(supplierData)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/produk/supplier", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Log("Product management setup completed")
	})

	// Step 6: Test Financial System
	t.Run("06_FinancialSystem", func(t *testing.T) {
		client := &http.Client{}

		// Create COA Account
		coaData := map[string]interface{}{
			"koperasi_id":   koperasiID,
			"kode_akun":     "1-1100",
			"nama_akun":     "Bank BCA",
			"kategori_id":   1,
			"level_akun":    2,
			"saldo_normal":  "debit",
			"is_kas":        false,
		}

		jsonData, _ := json.Marshal(coaData)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/financial/coa/akun", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create Journal Entry
		jurnalData := map[string]interface{}{
			"koperasi_id":        koperasiID,
			"nomor_jurnal":       "JU-2024-001",
			"tanggal_transaksi":  time.Now().Format("2006-01-02"),
			"referensi":          "Setoran Modal Awal",
			"keterangan":         "Penyetoran modal awal koperasi",
			"details": []map[string]interface{}{
				{
					"akun_id":    1,
					"keterangan": "Kas masuk dari setoran modal",
					"debit":      10000000,
					"kredit":     0,
				},
				{
					"akun_id":    2,
					"keterangan": "Modal disetor anggota",
					"debit":      0,
					"kredit":     10000000,
				},
			},
		}

		jsonData, _ = json.Marshal(jurnalData)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/financial/jurnal", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Log("Financial system setup completed")
	})

	// Step 7: Test Simpan Pinjam System
	t.Run("07_SimpanPinjamSystem", func(t *testing.T) {
		client := &http.Client{}

		// Create Savings Product
		produkData := map[string]interface{}{
			"koperasi_id":       koperasiID,
			"kode_produk":       "SISUKARELA",
			"nama_produk":       "Simpanan Sukarela",
			"jenis":             "simpanan",
			"kategori":          "sukarela",
			"bunga_simpanan":    4.5,
			"minimal_saldo":     100000,
			"syarat_ketentuan":  "Minimal setoran Rp 100.000, dapat diambil sewaktu-waktu",
		}

		jsonData, _ := json.Marshal(produkData)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/simpan-pinjam/produk", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create Account
		accountData := map[string]interface{}{
			"koperasi_id":    koperasiID,
			"anggota_id":     1,
			"produk_id":      1,
			"nomor_rekening": fmt.Sprintf("%.0f-001-0001", koperasiID),
			"saldo_simpanan": 500000,
		}

		jsonData, _ = json.Marshal(accountData)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/simpan-pinjam/rekening", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Log("Simpan Pinjam system setup completed")
	})

	// Step 8: Test PPOB System
	t.Run("08_PPOBSystem", func(t *testing.T) {
		// Get PPOB Categories (Public endpoint)
		resp, err := http.Get(baseURL + "/api/v1/ppob/kategoris")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var categories []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&categories)
		t.Logf("Found %d PPOB categories", len(categories))

		if len(categories) > 0 {
			categoryID := categories[0]["id"].(float64)

			// Get Products by Category
			resp, err = http.Get(fmt.Sprintf("%s/api/v1/ppob/kategoris/%.0f/produks", baseURL, categoryID))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var products []map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&products)
			t.Logf("Found %d products in category %.0f", len(products), categoryID)
		}
	})

	// Step 9: Test Clinic System
	t.Run("09_ClinicSystem", func(t *testing.T) {
		client := &http.Client{}

		// Create Patient
		patientData := map[string]interface{}{
			"koperasi_id":      koperasiID,
			"nomor_rm":         fmt.Sprintf("RM%.0f001", koperasiID),
			"nik":              "3201012345678902",
			"nama_lengkap":     "Siti Rahayu",
			"jenis_kelamin":    "P",
			"tempat_lahir":     "Bandung",
			"tanggal_lahir":    "1985-05-15",
			"alamat":           "Jl. Kesehatan No. 234",
			"rt":               "02",
			"rw":               "03",
			"kelurahan_id":     1101011001,
			"telepon":          "081345678901",
			"email":            "siti.rahayu@email.com",
			"golongan_darah":   "B",
			"alergi":           "Tidak ada alergi yang diketahui",
			"riwayat_penyakit": "Hipertensi ringan",
			"anggota_id":       1,
		}

		jsonData, _ := json.Marshal(patientData)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/klinik/pasien", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Create Medical Staff
		medicalStaffData := map[string]interface{}{
			"koperasi_id":      koperasiID,
			"nip":              "DOK001",
			"nama_lengkap":     "Dr. Ahmad Wijaya, Sp.PD",
			"jenis_kelamin":    "L",
			"spesialisasi":     "Penyakit Dalam",
			"no_str":           "STR123456789",
			"no_sip":           "SIP987654321",
			"telepon":          "081456789012",
			"email":            "dr.ahmad@klinikkoperasi.com",
			"jadwal_praktik":   `{"senin": "08:00-12:00", "rabu": "14:00-17:00", "jumat": "08:00-12:00"}`,
			"tarif_konsultasi": 150000,
		}

		jsonData, _ = json.Marshal(medicalStaffData)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/klinik/tenaga-medis", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		t.Log("Clinic system setup completed")
	})

	// Step 10: Test Reporting System
	t.Run("10_ReportingSystem", func(t *testing.T) {
		client := &http.Client{}

		// Test Quick Summary
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reports/%.0f/quick-summary", baseURL, koperasiID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var summary map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&summary)
		t.Logf("Quick summary: %+v", summary)

		// Test Dashboard
		req, _ = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/reports/%.0f/dashboard", baseURL, koperasiID), nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var dashboard map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&dashboard)
		t.Logf("Dashboard data retrieved successfully")
	})

	t.Log("🎉 All endpoint tests completed successfully!")
}

// TestCompleteWorkflowScenario tests a complete business workflow
func TestCompleteWorkflowScenario(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := &http.Client{}

	t.Log("=== COMPLETE BUSINESS WORKFLOW TEST ===")

	// Get auth token first
	loginData := map[string]interface{}{
		"email":    "admin@demo.local",
		"password": "admin123",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	var loginResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResponse)
	authToken := loginResponse["data"].(map[string]interface{})["token"].(string)

	// Scenario: New member joins, makes transactions, visits clinic
	t.Run("CompleteBusinessScenario", func(t *testing.T) {
		// 1. New user registration
		t.Log("Step 1: New member registration")
		memberData := map[string]interface{}{
			"koperasi_id":        1,
			"nik":               "3201998877665544",
			"nama_lengkap":      "Dewi Sartika",
			"email":             "dewi.sartika@email.com",
			"no_telepon":        "081567890123",
			"alamat":            "Jl. Pahlawan No. 567",
			"kelurahan_id":      1101011001,
			"password":          "DewiSecure456!",
			"jenis_keanggotaan": "anggota",
			"rencana_simpanan":  3000000,
		}

		jsonData, _ := json.Marshal(memberData)
		resp, err := http.Post(baseURL+"/api/v1/users/register", "application/json", bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 2. Create savings account
		t.Log("Step 2: Create savings account")
		savingsAccount := map[string]interface{}{
			"koperasi_id":    1,
			"anggota_id":     1,
			"produk_id":      1,
			"nomor_rekening": "1-SV-0003",
			"saldo_simpanan": 1000000,
		}

		jsonData, _ = json.Marshal(savingsAccount)
		req, _ := http.NewRequest("POST", baseURL+"/api/v1/simpan-pinjam/rekening", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 3. Make deposit transaction
		t.Log("Step 3: Make deposit transaction")
		depositTxn := map[string]interface{}{
			"rekening_id":       1,
			"jenis_transaksi":   "setor",
			"jumlah":            500000,
			"keterangan":        "Setoran rutin bulanan",
			"tanggal_transaksi": time.Now().Format("2006-01-02"),
		}

		jsonData, _ = json.Marshal(depositTxn)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/simpan-pinjam/transaksi", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 4. Buy products from cooperative store
		t.Log("Step 4: Purchase products")
		purchase := map[string]interface{}{
			"koperasi_id":       1,
			"anggota_id":        1,
			"nomor_transaksi":   "TXN-" + fmt.Sprintf("%d", time.Now().Unix()),
			"tanggal_transaksi": time.Now().Format("2006-01-02"),
			"items": []map[string]interface{}{
				{
					"produk_id":    1,
					"qty":          3,
					"harga_satuan": 50000,
				},
				{
					"produk_id":    2,
					"qty":          2,
					"harga_satuan": 15000,
				},
			},
			"metode_pembayaran": "cash",
			"jumlah_bayar":      180000,
			"jumlah_kembalian":  0,
			"kasir":             "Admin Koperasi",
		}

		jsonData, _ = json.Marshal(purchase)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/produk/penjualan", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 5. Pay utility bills through PPOB
		t.Log("Step 5: PPOB transaction")
		ppobTxn := map[string]interface{}{
			"koperasi_id":     1,
			"anggota_id":      1,
			"produk_id":       1,
			"nomor_tujuan":    "123456789012",
			"nama_pelanggan":  "Dewi Sartika",
			"customer_name":   "Dewi Sartika",
			"customer_email":  "dewi.sartika@email.com",
			"customer_phone":  "081567890123",
		}

		jsonData, _ = json.Marshal(ppobTxn)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/ppob/transactions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 6. Visit clinic
		t.Log("Step 6: Clinic visit")
		clinicVisit := map[string]interface{}{
			"koperasi_id":       1,
			"pasien_id":         1,
			"dokter_id":         1,
			"tanggal_kunjungan": time.Now().Format("2006-01-02"),
			"keluhan":           "Merasa lelah dan pusing dalam beberapa hari terakhir",
			"anamnesis":         "Pasien mengeluh lelah, pusing, dan kurang tidur",
			"pemeriksaan_fisik": "TD: 130/80 mmHg, Nadi: 88x/menit, Suhu: 36.5°C",
			"diagnosis":         "Hipertensi stage 1, stress",
			"terapi":            "Istirahat cukup, olahraga ringan, kontrol diet garam",
			"resep_obat":        "Amlodipine 5mg 1x1, Vitamin B kompleks 1x1",
			"biaya_konsultasi":  150000,
			"biaya_obat":        75000,
			"total_biaya":       225000,
		}

		jsonData, _ = json.Marshal(clinicVisit)
		req, _ = http.NewRequest("POST", baseURL+"/api/v1/klinik/kunjungan", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// 7. Generate financial report
		t.Log("Step 7: Generate reports")
		req, _ = http.NewRequest("GET", baseURL+"/api/v1/reports/1/quick-summary", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var report map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&report)
		t.Logf("Final business summary: %+v", report)

		t.Log("✅ Complete business workflow executed successfully!")
	})
}