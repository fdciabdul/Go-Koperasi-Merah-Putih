package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=yourpassword dbname=koperasi_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Starting seeder...")

	// Run seeders in order
	seedProvinsi(db)
	seedKabupaten(db)
	seedKecamatan(db)
	seedKelurahan(db)
	seedKBLI(db)
	seedJenisKoperasi(db)
	seedBentukKoperasi(db)
	seedCOAKategori(db)
	seedTenants(db)
	seedUsers(db)
	seedKoperasi(db)
	seedAnggotaKoperasi(db)
	seedCOAAkun(db)
	seedSimpanPinjamProduk(db)
	seedSimpanPinjamRekening(db)
	seedSimpanPinjamTransaksi(db)
	seedKlinikTenagaMedis(db)
	seedKlinikPasien(db)
	seedKlinikObat(db)
	seedKlinikKunjungan(db)
	seedJurnalUmum(db)
	seedPPOBKategori(db)
	seedPPOBProduk(db)
	seedKategoriProduk(db)
	seedSatuanProduk(db)
	seedSupplier(db)
	seedProduk(db)
	seedSequences(db)

	fmt.Println("Seeder completed successfully!")
}

func seedProvinsi(db *gorm.DB) {
	provinsis := []postgres.Provinsi{
		{ID: 11, Kode: "11", Nama: "ACEH"},
		{ID: 12, Kode: "12", Nama: "SUMATERA UTARA"},
		{ID: 31, Kode: "31", Nama: "DKI JAKARTA"},
		{ID: 32, Kode: "32", Nama: "JAWA BARAT"},
		{ID: 33, Kode: "33", Nama: "JAWA TENGAH"},
		{ID: 34, Kode: "34", Nama: "DI YOGYAKARTA"},
		{ID: 35, Kode: "35", Nama: "JAWA TIMUR"},
	}

	for _, provinsi := range provinsis {
		db.FirstOrCreate(&provinsi, postgres.Provinsi{ID: provinsi.ID})
	}
	fmt.Println("✓ Seeded Provinsi")
}

func seedKabupaten(db *gorm.DB) {
	kabupatens := []postgres.Kabupaten{
		{ID: 3171, ProvinsiID: 31, Kode: "3171", Nama: "KEPULAUAN SERIBU"},
		{ID: 3201, ProvinsiID: 32, Kode: "3201", Nama: "BOGOR"},
		{ID: 3202, ProvinsiID: 32, Kode: "3202", Nama: "SUKABUMI"},
		{ID: 3203, ProvinsiID: 32, Kode: "3203", Nama: "CIANJUR"},
		{ID: 3204, ProvinsiID: 32, Kode: "3204", Nama: "BANDUNG"},
	}

	for _, kabupaten := range kabupatens {
		db.FirstOrCreate(&kabupaten, postgres.Kabupaten{ID: kabupaten.ID})
	}
	fmt.Println("✓ Seeded Kabupaten")
}

func seedKecamatan(db *gorm.DB) {
	kecamatans := []postgres.Kecamatan{
		{ID: 320101, KabupatenID: 3201, Kode: "320101", Nama: "NANGGUNG"},
		{ID: 320102, KabupatenID: 3201, Kode: "320102", Nama: "LEUWILIANG"},
		{ID: 320103, KabupatenID: 3201, Kode: "320103", Nama: "LEUWISADENG"},
		{ID: 320104, KabupatenID: 3201, Kode: "320104", Nama: "PAMIJAHAN"},
		{ID: 320105, KabupatenID: 3201, Kode: "320105", Nama: "CIBUNGBULANG"},
	}

	for _, kecamatan := range kecamatans {
		db.FirstOrCreate(&kecamatan, postgres.Kecamatan{ID: kecamatan.ID})
	}
	fmt.Println("✓ Seeded Kecamatan")
}

func seedKelurahan(db *gorm.DB) {
	kelurahans := []postgres.Kelurahan{
		{ID: 3201011001, KecamatanID: 320101, Kode: "3201011001", Nama: "NANGGUNG"},
		{ID: 3201011002, KecamatanID: 320101, Kode: "3201011002", Nama: "PARAKANMUNCANG"},
		{ID: 3201011003, KecamatanID: 320101, Kode: "3201011003", Nama: "CURUGBITUNG"},
		{ID: 3201011004, KecamatanID: 320101, Kode: "3201011004", Nama: "BANTARJAYA"},
		{ID: 3201011005, KecamatanID: 320101, Kode: "3201011005", Nama: "HAMBARO"},
	}

	for _, kelurahan := range kelurahans {
		db.FirstOrCreate(&kelurahan, postgres.Kelurahan{ID: kelurahan.ID})
	}
	fmt.Println("✓ Seeded Kelurahan")
}

func seedKBLI(db *gorm.DB) {
	kblis := []postgres.KBLI{
		{Kode: "64191", Nama: "Bank Sentral", Kategori: "Jasa Keuangan", Deskripsi: "Kegiatan bank sentral", IsAktif: true},
		{Kode: "64921", Nama: "Kegiatan Koperasi Kredit", Kategori: "Jasa Keuangan", Deskripsi: "Kegiatan koperasi kredit dan simpan pinjam", IsAktif: true},
		{Kode: "64922", Nama: "Kegiatan Koperasi Non-Kredit", Kategori: "Jasa Keuangan", Deskripsi: "Kegiatan koperasi selain kredit", IsAktif: true},
		{Kode: "86101", Nama: "Kegiatan Rumah Sakit", Kategori: "Kesehatan", Deskripsi: "Kegiatan pelayanan rumah sakit", IsAktif: true},
		{Kode: "86201", Nama: "Praktik Dokter Umum", Kategori: "Kesehatan", Deskripsi: "Praktik dokter dan dokter gigi umum", IsAktif: true},
	}

	for _, kbli := range kblis {
		db.FirstOrCreate(&kbli, postgres.KBLI{Kode: kbli.Kode})
	}
	fmt.Println("✓ Seeded KBLI")
}

func seedJenisKoperasi(db *gorm.DB) {
	jenisKoperasis := []postgres.JenisKoperasi{
		{Kode: "KSP", Nama: "Koperasi Simpan Pinjam", Deskripsi: "Koperasi yang bergerak di bidang simpan pinjam", IsAktif: true},
		{Kode: "KSU", Nama: "Koperasi Serba Usaha", Deskripsi: "Koperasi dengan berbagai bidang usaha", IsAktif: true},
		{Kode: "KJKS", Nama: "Koperasi Jasa Keuangan Syariah", Deskripsi: "Koperasi jasa keuangan berbasis syariah", IsAktif: true},
	}

	for _, jenis := range jenisKoperasis {
		db.FirstOrCreate(&jenis, postgres.JenisKoperasi{Kode: jenis.Kode})
	}
	fmt.Println("✓ Seeded Jenis Koperasi")
}

func seedBentukKoperasi(db *gorm.DB) {
	bentukKoperasis := []postgres.BentukKoperasi{
		{Kode: "PRIMER", Nama: "Koperasi Primer", Deskripsi: "Koperasi yang beranggotakan orang perorangan", IsAktif: true},
		{Kode: "SEKUNDER", Nama: "Koperasi Sekunder", Deskripsi: "Koperasi yang beranggotakan koperasi primer", IsAktif: true},
	}

	for _, bentuk := range bentukKoperasis {
		db.FirstOrCreate(&bentuk, postgres.BentukKoperasi{Kode: bentuk.Kode})
	}
	fmt.Println("✓ Seeded Bentuk Koperasi")
}

func seedCOAKategori(db *gorm.DB) {
	kategoris := []postgres.COAKategori{
		{Nama: "ASET", Tipe: "aset", Urutan: 1, IsAktif: true},
		{Nama: "KEWAJIBAN", Tipe: "kewajiban", Urutan: 2, IsAktif: true},
		{Nama: "EKUITAS", Tipe: "ekuitas", Urutan: 3, IsAktif: true},
		{Nama: "PENDAPATAN", Tipe: "pendapatan", Urutan: 4, IsAktif: true},
		{Nama: "BEBAN", Tipe: "beban", Urutan: 5, IsAktif: true},
	}

	for _, kategori := range kategoris {
		db.FirstOrCreate(&kategori, postgres.COAKategori{Nama: kategori.Nama})
	}
	fmt.Println("✓ Seeded COA Kategori")
}

func seedTenants(db *gorm.DB) {
	tenants := []postgres.Tenant{
		{
			Nama:        "Koperasi Merah Putih",
			Domain:      "koperasi-merah-putih.com",
			DatabaseURL: "postgresql://localhost/koperasi_db",
			IsActive:    true,
		},
		{
			Nama:        "Koperasi Sejahtera",
			Domain:      "koperasi-sejahtera.com",
			DatabaseURL: "postgresql://localhost/koperasi_db",
			IsActive:    true,
		},
	}

	for _, tenant := range tenants {
		db.FirstOrCreate(&tenant, postgres.Tenant{Domain: tenant.Domain})
	}
	fmt.Println("✓ Seeded Tenants")
}

func seedUsers(db *gorm.DB) {
	now := time.Now()
	users := []postgres.User{
		{
			TenantID:     1,
			Email:        "admin@koperasi.com",
			PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: admin123
			Role:         "super_admin",
			IsActive:     true,
			EmailVerifiedAt: &now,
		},
		{
			TenantID:     1,
			Email:        "manager@koperasi.com",
			PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: admin123
			Role:         "admin",
			IsActive:     true,
			EmailVerifiedAt: &now,
		},
		{
			TenantID:     1,
			Email:        "staff@koperasi.com",
			PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: admin123
			Role:         "staff",
			IsActive:     true,
			EmailVerifiedAt: &now,
		},
	}

	for _, user := range users {
		db.FirstOrCreate(&user, postgres.User{Email: user.Email})
	}
	fmt.Println("✓ Seeded Users")
}

func seedKoperasi(db *gorm.DB) {
	koperasis := []postgres.Koperasi{
		{
			TenantID:       1,
			Nama:           "Koperasi Sejahtera Mandiri",
			NIAK:           "1234567890123456",
			Email:          "info@sejahtera.com",
			Telepon:        "02112345678",
			Website:        "www.sejahtera.com",
			Alamat:         "Jl. Merdeka No. 123",
			ProvinsiID:     32,
			KabupatenID:    3201,
			KecamatanID:    320101,
			KelurahanID:    3201011001,
			KodePos:        "16610",
			JenisKoperasiID: 1,
			BentukKoperasiID: 1,
			KBLIID:         2,
			TanggalBerdiri: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			IsActive:       true,
		},
		{
			TenantID:       1,
			Nama:           "Koperasi Makmur Bersama",
			NIAK:           "1234567890123457",
			Email:          "info@makmur.com",
			Telepon:        "02187654321",
			Website:        "www.makmur.com",
			Alamat:         "Jl. Proklamasi No. 456",
			ProvinsiID:     32,
			KabupatenID:    3201,
			KecamatanID:    320101,
			KelurahanID:    3201011002,
			KodePos:        "16610",
			JenisKoperasiID: 2,
			BentukKoperasiID: 1,
			KBLIID:         3,
			TanggalBerdiri: time.Date(2019, 6, 10, 0, 0, 0, 0, time.UTC),
			IsActive:       true,
		},
	}

	for _, koperasi := range koperasis {
		db.FirstOrCreate(&koperasi, postgres.Koperasi{NIAK: koperasi.NIAK})
	}
	fmt.Println("✓ Seeded Koperasi")
}

func seedAnggotaKoperasi(db *gorm.DB) {
	now := time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC)
	anggotas := []postgres.AnggotaKoperasi{
		{
			KoperasiID:     1,
			NomorAnggota:   "A001",
			NIK:            "3201011505850001",
			NamaLengkap:    "Ahmad Suryadi",
			JenisKelamin:   "L",
			TempatLahir:    "Bogor",
			TanggalLahir:   &now,
			Alamat:         "Jl. Mawar No. 10",
			Telepon:        "081234567890",
			Email:          "ahmad@email.com",
			Pekerjaan:      "Pegawai Swasta",
			StatusPernikahan: "Menikah",
			Status:         "aktif",
			TanggalBergabung: time.Now(),
		},
		{
			KoperasiID:     1,
			NomorAnggota:   "A002",
			NIK:            "3201011505850002",
			NamaLengkap:    "Siti Nurhaliza",
			JenisKelamin:   "P",
			TempatLahir:    "Jakarta",
			TanggalLahir:   &now,
			Alamat:         "Jl. Melati No. 15",
			Telepon:        "081234567891",
			Email:          "siti@email.com",
			Pekerjaan:      "Guru",
			StatusPernikahan: "Menikah",
			Status:         "aktif",
			TanggalBergabung: time.Now(),
		},
	}

	for _, anggota := range anggotas {
		db.FirstOrCreate(&anggota, postgres.AnggotaKoperasi{NomorAnggota: anggota.NomorAnggota, KoperasiID: anggota.KoperasiID})
	}
	fmt.Println("✓ Seeded Anggota Koperasi")
}

func seedCOAAkun(db *gorm.DB) {
	akuns := []postgres.COAAkun{
		{TenantID: 1, KoperasiID: 1, KodeAkun: "1001", NamaAkun: "Kas", KategoriID: 1, SaldoNormal: "debit", IsKas: true, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "1101", NamaAkun: "Bank BCA", KategoriID: 1, SaldoNormal: "debit", IsKas: true, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "1201", NamaAkun: "Piutang Anggota", KategoriID: 1, SaldoNormal: "debit", IsKas: false, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "2001", NamaAkun: "Simpanan Pokok", KategoriID: 2, SaldoNormal: "kredit", IsKas: false, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "2002", NamaAkun: "Simpanan Wajib", KategoriID: 2, SaldoNormal: "kredit", IsKas: false, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "3001", NamaAkun: "Modal Koperasi", KategoriID: 3, SaldoNormal: "kredit", IsKas: false, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "4001", NamaAkun: "Pendapatan Bunga", KategoriID: 4, SaldoNormal: "kredit", IsKas: false, IsAktif: true},
		{TenantID: 1, KoperasiID: 1, KodeAkun: "5001", NamaAkun: "Beban Operasional", KategoriID: 5, SaldoNormal: "debit", IsKas: false, IsAktif: true},
	}

	for _, akun := range akuns {
		db.FirstOrCreate(&akun, postgres.COAAkun{KodeAkun: akun.KodeAkun, KoperasiID: akun.KoperasiID})
	}
	fmt.Println("✓ Seeded COA Akun")
}

func seedSimpanPinjamProduk(db *gorm.DB) {
	produks := []postgres.SimpanPinjamProduk{
		{
			KoperasiID:       1,
			KodeProduk:       "SP001",
			NamaProduk:       "Simpanan Berjangka 12 Bulan",
			JenisProduk:      "simpanan",
			Deskripsi:        "Simpanan berjangka dengan tenor 12 bulan",
			BungaTahun:       6.5,
			TenorBulan:       12,
			MinimalSetoran:   1000000,
			MaksimalSetoran:  50000000,
			BiayaAdmin:       10000,
			IsAktif:          true,
		},
		{
			KoperasiID:       1,
			KodeProduk:       "SP002",
			NamaProduk:       "Pinjaman Konsumtif",
			JenisProduk:      "pinjaman",
			Deskripsi:        "Pinjaman untuk kebutuhan konsumtif",
			BungaTahun:       18.0,
			TenorBulan:       24,
			MinimalSetoran:   5000000,
			MaksimalSetoran:  100000000,
			BiayaAdmin:       100000,
			IsAktif:          true,
		},
	}

	for _, produk := range produks {
		db.FirstOrCreate(&produk, postgres.SimpanPinjamProduk{KodeProduk: produk.KodeProduk})
	}
	fmt.Println("✓ Seeded Simpan Pinjam Produk")
}

func seedSimpanPinjamRekening(db *gorm.DB) {
	rekenings := []postgres.SimpanPinjamRekening{
		{
			ProdukID:       1,
			AnggotaID:      1,
			NomorRekening:  "SP001001",
			SaldoPokok:     5000000,
			SaldoBunga:     0,
			Status:         "aktif",
			TanggalBuka:    time.Now().AddDate(0, -6, 0),
		},
		{
			ProdukID:       2,
			AnggotaID:      2,
			NomorRekening:  "SP002001",
			SaldoPokok:     25000000,
			SaldoBunga:     0,
			Status:         "aktif",
			TanggalBuka:    time.Now().AddDate(0, -3, 0),
		},
	}

	for _, rekening := range rekenings {
		db.FirstOrCreate(&rekening, postgres.SimpanPinjamRekening{NomorRekening: rekening.NomorRekening})
	}
	fmt.Println("✓ Seeded Simpan Pinjam Rekening")
}

func seedSimpanPinjamTransaksi(db *gorm.DB) {
	transaksis := []postgres.SimpanPinjamTransaksi{
		{
			RekeningID:     1,
			NomorTransaksi: "T001",
			JenisTransaksi: "setoran",
			Nominal:        5000000,
			Deskripsi:      "Setoran awal simpanan berjangka",
			TanggalTransaksi: time.Now().AddDate(0, -6, 0),
			Status:         "berhasil",
			UserID:         2,
		},
		{
			RekeningID:     2,
			NomorTransaksi: "T002",
			JenisTransaksi: "pencairan",
			Nominal:        25000000,
			Deskripsi:      "Pencairan pinjaman konsumtif",
			TanggalTransaksi: time.Now().AddDate(0, -3, 0),
			Status:         "berhasil",
			UserID:         2,
		},
	}

	for _, transaksi := range transaksis {
		db.FirstOrCreate(&transaksi, postgres.SimpanPinjamTransaksi{NomorTransaksi: transaksi.NomorTransaksi})
	}
	fmt.Println("✓ Seeded Simpan Pinjam Transaksi")
}

func seedKlinikTenagaMedis(db *gorm.DB) {
	tenagaMedis := []postgres.KlinikTenagaMedis{
		{
			KoperasiID:       1,
			NIP:              "DOK001",
			NamaLengkap:      "Dr. Budi Santoso",
			JenisKelamin:     "L",
			Spesialisasi:     "Dokter Umum",
			NoSTR:            "STR001234567890",
			NoSIP:            "SIP001234567890",
			Telepon:          "081234567892",
			Email:            "dr.budi@klinik.com",
			JadwalPraktik:    "Senin-Jumat 08:00-17:00",
			TarifKonsultasi:  100000,
			Status:           "aktif",
		},
		{
			KoperasiID:       1,
			NIP:              "NUR001",
			NamaLengkap:      "Ns. Dewi Sartika",
			JenisKelamin:     "P",
			Spesialisasi:     "Perawat",
			NoSTR:            "STR001234567891",
			NoSIP:            "",
			Telepon:          "081234567893",
			Email:            "ns.dewi@klinik.com",
			JadwalPraktik:    "Senin-Jumat 07:00-15:00",
			TarifKonsultasi:  0,
			Status:           "aktif",
		},
	}

	for _, tm := range tenagaMedis {
		db.FirstOrCreate(&tm, postgres.KlinikTenagaMedis{NIP: tm.NIP})
	}
	fmt.Println("✓ Seeded Klinik Tenaga Medis")
}

func seedKlinikPasien(db *gorm.DB) {
	lahir := time.Date(1990, 3, 15, 0, 0, 0, 0, time.UTC)
	pasiens := []postgres.KlinikPasien{
		{
			KoperasiID:      1,
			NomorRM:         "RM0001000001",
			NIK:             "3201011503900001",
			NamaLengkap:     "Rina Sari",
			JenisKelamin:    "P",
			TempatLahir:     "Bandung",
			TanggalLahir:    &lahir,
			Alamat:          "Jl. Dahlia No. 20",
			Telepon:         "081234567894",
			Email:           "rina@email.com",
			GolonganDarah:   "A",
			Alergi:          "Tidak ada",
			RiwayatPenyakit: "Tidak ada",
			AnggotaID:       1,
		},
		{
			KoperasiID:      1,
			NomorRM:         "RM0001000002",
			NIK:             "3201011503900002",
			NamaLengkap:     "Joko Susilo",
			JenisKelamin:    "L",
			TempatLahir:     "Bogor",
			TanggalLahir:    &lahir,
			Alamat:          "Jl. Anggrek No. 25",
			Telepon:         "081234567895",
			Email:           "joko@email.com",
			GolonganDarah:   "B",
			Alergi:          "Seafood",
			RiwayatPenyakit: "Hipertensi",
			AnggotaID:       2,
		},
	}

	for _, pasien := range pasiens {
		db.FirstOrCreate(&pasien, postgres.KlinikPasien{NomorRM: pasien.NomorRM})
	}
	fmt.Println("✓ Seeded Klinik Pasien")
}

func seedKlinikObat(db *gorm.DB) {
	obats := []postgres.KlinikObat{
		{
			KoperasiID:    1,
			KodeObat:      "OBT001",
			NamaObat:      "Paracetamol 500mg",
			Kategori:      "Analgesik",
			BentukSediaan: "Tablet",
			Kekuatan:      "500mg",
			Satuan:        "Tablet",
			StokMinimal:   50,
			StokCurrent:   200,
			HargaBeli:     500,
			HargaJual:     1000,
			IsAktif:       true,
		},
		{
			KoperasiID:    1,
			KodeObat:      "OBT002",
			NamaObat:      "Amoxicillin 500mg",
			Kategori:      "Antibiotik",
			BentukSediaan: "Kapsul",
			Kekuatan:      "500mg",
			Satuan:        "Kapsul",
			StokMinimal:   30,
			StokCurrent:   100,
			HargaBeli:     2000,
			HargaJual:     3500,
			IsAktif:       true,
		},
	}

	for _, obat := range obats {
		db.FirstOrCreate(&obat, postgres.KlinikObat{KodeObat: obat.KodeObat})
	}
	fmt.Println("✓ Seeded Klinik Obat")
}

func seedKlinikKunjungan(db *gorm.DB) {
	kunjungans := []postgres.KlinikKunjungan{
		{
			KoperasiID:       1,
			PasienID:         1,
			DokterID:         1,
			NomorKunjungan:   "KUN0001000001",
			TanggalKunjungan: time.Now().AddDate(0, 0, -7),
			KeluhanUtama:     "Demam dan sakit kepala",
			Anamnesis:        "Pasien mengeluh demam sejak 2 hari yang lalu",
			PemeriksaanFisik: "TD: 120/80 mmHg, Nadi: 80x/menit, Suhu: 38°C",
			Diagnosis:        "Demam tifoid suspek",
			TerapiPengobatan: "Istirahat, minum obat teratur",
			BiayaKonsultasi:  100000,
			BiayaTindakan:    0,
			BiayaObat:        7000,
			TotalBiaya:       107000,
			StatusPembayaran: "lunas",
		},
		{
			KoperasiID:       1,
			PasienID:         2,
			DokterID:         1,
			NomorKunjungan:   "KUN0001000002",
			TanggalKunjungan: time.Now().AddDate(0, 0, -3),
			KeluhanUtama:     "Batuk dan pilek",
			Anamnesis:        "Pasien batuk berdahak sejak 1 minggu",
			PemeriksaanFisik: "TD: 130/90 mmHg, Nadi: 85x/menit, Suhu: 36.5°C",
			Diagnosis:        "ISPA",
			TerapiPengobatan: "Antibiotik dan ekspektoran",
			BiayaKonsultasi:  100000,
			BiayaTindakan:    25000,
			BiayaObat:        14000,
			TotalBiaya:       139000,
			StatusPembayaran: "lunas",
		},
	}

	for _, kunjungan := range kunjungans {
		db.FirstOrCreate(&kunjungan, postgres.KlinikKunjungan{NomorKunjungan: kunjungan.NomorKunjungan})
	}
	fmt.Println("✓ Seeded Klinik Kunjungan")
}

func seedJurnalUmum(db *gorm.DB) {
	jurnals := []postgres.JurnalUmum{
		{
			TenantID:         1,
			KoperasiID:       1,
			NomorJurnal:      "JU20240101000001",
			TanggalTransaksi: time.Now().AddDate(0, -1, 0),
			Referensi:        "SETORAN-001",
			Keterangan:       "Penerimaan simpanan pokok anggota A001",
			TotalDebit:       1000000,
			TotalKredit:      1000000,
			Status:           "posted",
			CreatedBy:        2,
			PostedBy:         2,
			PostedAt:         timePtr(time.Now().AddDate(0, -1, 0)),
		},
		{
			TenantID:         1,
			KoperasiID:       1,
			NomorJurnal:      "JU20240102000001",
			TanggalTransaksi: time.Now().AddDate(0, 0, -15),
			Referensi:        "PINJAMAN-001",
			Keterangan:       "Pencairan pinjaman anggota A002",
			TotalDebit:       25000000,
			TotalKredit:      25000000,
			Status:           "posted",
			CreatedBy:        2,
			PostedBy:         2,
			PostedAt:         timePtr(time.Now().AddDate(0, 0, -15)),
		},
	}

	for _, jurnal := range jurnals {
		db.FirstOrCreate(&jurnal, postgres.JurnalUmum{NomorJurnal: jurnal.NomorJurnal})
	}

	// Seed Jurnal Details
	details := []postgres.JurnalDetail{
		// Jurnal 1: Kas Debit, Simpanan Pokok Kredit
		{JurnalID: 1, AkunID: 1, Keterangan: "Penerimaan kas dari simpanan pokok", Debit: 1000000, Kredit: 0},
		{JurnalID: 1, AkunID: 4, Keterangan: "Simpanan pokok anggota A001", Debit: 0, Kredit: 1000000},

		// Jurnal 2: Piutang Debit, Kas Kredit
		{JurnalID: 2, AkunID: 3, Keterangan: "Pencairan pinjaman anggota A002", Debit: 25000000, Kredit: 0},
		{JurnalID: 2, AkunID: 1, Keterangan: "Pengeluaran kas untuk pinjaman", Debit: 0, Kredit: 25000000},
	}

	for _, detail := range details {
		db.FirstOrCreate(&detail, postgres.JurnalDetail{JurnalID: detail.JurnalID, AkunID: detail.AkunID})
	}

	fmt.Println("✓ Seeded Jurnal Umum & Details")
}

func seedPPOBKategori(db *gorm.DB) {
	kategoris := []postgres.PPOBKategori{
		{Nama: "Pulsa & Paket Data", Kode: "PULSA", IsAktif: true},
		{Nama: "PLN", Kode: "PLN", IsAktif: true},
		{Nama: "PDAM", Kode: "PDAM", IsAktif: true},
		{Nama: "Internet & TV Kabel", Kode: "INTERNET", IsAktif: true},
	}

	for _, kategori := range kategoris {
		db.FirstOrCreate(&kategori, postgres.PPOBKategori{Kode: kategori.Kode})
	}
	fmt.Println("✓ Seeded PPOB Kategori")
}

func seedPPOBProduk(db *gorm.DB) {
	produks := []postgres.PPOBProduk{
		{KategoriID: 1, Nama: "Telkomsel 10.000", Kode: "TSEL10", Harga: 10500, Deskripsi: "Pulsa Telkomsel 10rb", IsAktif: true},
		{KategoriID: 1, Nama: "Indosat 25.000", Kode: "ISAT25", Harga: 25200, Deskripsi: "Pulsa Indosat 25rb", IsAktif: true},
		{KategoriID: 2, Nama: "PLN Token 20.000", Kode: "PLN20", Harga: 20500, Deskripsi: "Token listrik PLN 20rb", IsAktif: true},
		{KategoriID: 2, Nama: "PLN Token 50.000", Kode: "PLN50", Harga: 50500, Deskripsi: "Token listrik PLN 50rb", IsAktif: true},
	}

	for _, produk := range produks {
		db.FirstOrCreate(&produk, postgres.PPOBProduk{Kode: produk.Kode})
	}
	fmt.Println("✓ Seeded PPOB Produk")
}

func seedSequences(db *gorm.DB) {
	sequences := []postgres.Sequence{
		{TenantID: 1, KoperasiID: 1, SequenceType: "jurnal_umum", CurrentValue: 2},
		{TenantID: 1, KoperasiID: 1, SequenceType: "nomor_rm", CurrentValue: 2},
		{TenantID: 1, KoperasiID: 1, SequenceType: "kunjungan", CurrentValue: 2},
		{TenantID: 1, KoperasiID: 2, SequenceType: "jurnal_umum", CurrentValue: 0},
		{TenantID: 1, KoperasiID: 2, SequenceType: "nomor_rm", CurrentValue: 0},
	}

	for _, seq := range sequences {
		db.FirstOrCreate(&seq, postgres.Sequence{
			TenantID: seq.TenantID,
			KoperasiID: seq.KoperasiID,
			SequenceType: seq.SequenceType,
		})
	}
	fmt.Println("✓ Seeded Sequences")
}

func seedKategoriProduk(db *gorm.DB) {
	kategoris := []postgres.KategoriProduk{
		{Kode: "FOOD", Nama: "Makanan", Deskripsi: "Produk makanan dan cemilan", Icon: "🍽️", IsActive: true},
		{Kode: "BEVERAGE", Nama: "Minuman", Deskripsi: "Minuman segar dan sehat", Icon: "🥤", IsActive: true},
		{Kode: "LIVESTOCK", Nama: "Ternak", Deskripsi: "Hewan ternak seperti sapi, kambing, ayam", Icon: "🐄", IsActive: true},
		{Kode: "VEGETABLE", Nama: "Sayuran", Deskripsi: "Sayuran segar dan organik", Icon: "🥬", IsActive: true},
		{Kode: "FRUIT", Nama: "Buah-buahan", Deskripsi: "Buah segar dan vitamin", Icon: "🍎", IsActive: true},
		{Kode: "GRAIN", Nama: "Biji-bijian", Deskripsi: "Beras, gandum, jagung, kedelai", Icon: "🌾", IsActive: true},
		{Kode: "DAIRY", Nama: "Produk Susu", Deskripsi: "Susu dan produk olahannya", Icon: "🥛", IsActive: true},
		{Kode: "MEAT", Nama: "Daging", Deskripsi: "Daging segar dan olahan", Icon: "🥩", IsActive: true},
		{Kode: "FISH", Nama: "Ikan", Deskripsi: "Ikan segar dan hasil laut", Icon: "🐟", IsActive: true},
		{Kode: "SPICE", Nama: "Rempah", Deskripsi: "Rempah-rempah dan bumbu dapur", Icon: "🌶️", IsActive: true},
		{Kode: "HOUSEHOLD", Nama: "Keperluan Rumah Tangga", Deskripsi: "Sabun, deterjen, dan kebutuhan sehari-hari", Icon: "🧽", IsActive: true},
		{Kode: "HEALTH", Nama: "Kesehatan", Deskripsi: "Obat-obatan dan suplemen", Icon: "💊", IsActive: true},
	}

	for _, kategori := range kategoris {
		db.FirstOrCreate(&kategori, postgres.KategoriProduk{Kode: kategori.Kode})
	}
	fmt.Println("✓ Seeded Kategori Produk")
}

func seedSatuanProduk(db *gorm.DB) {
	satuans := []postgres.SatuanProduk{
		{Kode: "KG", Nama: "Kilogram", IsActive: true},
		{Kode: "GR", Nama: "Gram", IsActive: true},
		{Kode: "LTR", Nama: "Liter", IsActive: true},
		{Kode: "ML", Nama: "Mililiter", IsActive: true},
		{Kode: "PCS", Nama: "Pieces", IsActive: true},
		{Kode: "BTL", Nama: "Botol", IsActive: true},
		{Kode: "PACK", Nama: "Pack", IsActive: true},
		{Kode: "KTN", Nama: "Karton", IsActive: true},
		{Kode: "SACK", Nama: "Sak", IsActive: true},
		{Kode: "EKOR", Nama: "Ekor", IsActive: true},
		{Kode: "IKAT", Nama: "Ikat", IsActive: true},
		{Kode: "METER", Nama: "Meter", IsActive: true},
		{Kode: "TON", Nama: "Ton", IsActive: true},
		{Kode: "DOZEN", Nama: "Lusin", IsActive: true},
	}

	for _, satuan := range satuans {
		db.FirstOrCreate(&satuan, postgres.SatuanProduk{Kode: satuan.Kode})
	}
	fmt.Println("✓ Seeded Satuan Produk")
}

func seedSupplier(db *gorm.DB) {
	suppliers := []postgres.Supplier{
		{
			KoperasiID:     1,
			Kode:           "SUP001",
			Nama:           "CV Beras Sejahtera",
			KontakPerson:   "Budi Santoso",
			Telepon:        "081234567890",
			Email:          "budi@berassejahtera.com",
			Alamat:         "Jl. Raya Karawang No. 123",
			ProvinsiID:     32,
			KabupatenID:    3201,
			JenisSupplier:  "perusahaan",
			Status:         "aktif",
			TermPembayaran: 30,
			IsActive:       true,
			CreatedBy:      1,
			UpdatedBy:      1,
		},
		{
			KoperasiID:     1,
			Kode:           "SUP002",
			Nama:           "Toko Sayur Segar",
			KontakPerson:   "Ibu Siti",
			Telepon:        "081298765432",
			Email:          "siti@sayursegar.com",
			Alamat:         "Pasar Induk Kramat Jati Blok C",
			ProvinsiID:     31,
			KabupatenID:    3101,
			JenisSupplier:  "individu",
			Status:         "aktif",
			TermPembayaran: 14,
			IsActive:       true,
			CreatedBy:      1,
			UpdatedBy:      1,
		},
		{
			KoperasiID:     2,
			Kode:           "SUP003",
			Nama:           "PT Susu Murni Indonesia",
			KontakPerson:   "Ahmad Wijaya",
			Telepon:        "081345678901",
			Email:          "ahmad@susumurni.co.id",
			Alamat:         "Jl. Industri Raya No. 45, Bekasi",
			ProvinsiID:     32,
			KabupatenID:    3275,
			JenisSupplier:  "perusahaan",
			Status:         "aktif",
			TermPembayaran: 45,
			IsActive:       true,
			CreatedBy:      2,
			UpdatedBy:      2,
		},
	}

	for _, supplier := range suppliers {
		db.FirstOrCreate(&supplier, postgres.Supplier{Kode: supplier.Kode, KoperasiID: supplier.KoperasiID})
	}
	fmt.Println("✓ Seeded Supplier")
}

func seedProduk(db *gorm.DB) {
	produk := []postgres.Produk{
		{
			KoperasiID:       1,
			KategoriProdukID: 6,
			SatuanProdukID:   9,
			KodeProduk:       "PRD000100001",
			NamaProduk:       "Beras Premium 5kg",
			Deskripsi:        "Beras putih premium kualitas terbaik",
			Brand:            "Sania",
			HargaBeli:        45000,
			HargaJual:        52000,
			MarginPersen:     15.56,
			StokMinimal:      10,
			StokMaksimal:     100,
			StokCurrent:      50,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
		{
			KoperasiID:       1,
			KategoriProdukID: 4,
			SatuanProdukID:   1,
			KodeProduk:       "PRD000100002",
			NamaProduk:       "Sayur Kangkung",
			Deskripsi:        "Kangkung segar organik",
			Brand:            "Organik Nusantara",
			HargaBeli:        3000,
			HargaJual:        4500,
			MarginPersen:     50,
			StokMinimal:      5,
			StokMaksimal:     50,
			StokCurrent:      25,
			IsPerishable:     true,
			ShelfLife:        3,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
		{
			KoperasiID:       1,
			KategoriProdukID: 2,
			SatuanProdukID:   6,
			KodeProduk:       "PRD000100003",
			NamaProduk:       "Air Mineral Botol 600ml",
			Deskripsi:        "Air mineral murni dalam kemasan botol",
			Brand:            "Aqua",
			HargaBeli:        2500,
			HargaJual:        3500,
			MarginPersen:     40,
			StokMinimal:      20,
			StokMaksimal:     200,
			StokCurrent:      150,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
		{
			KoperasiID:       2,
			KategoriProdukID: 7,
			SatuanProdukID:   3,
			KodeProduk:       "PRD000200001",
			NamaProduk:       "Susu Sapi Murni 1L",
			Deskripsi:        "Susu sapi segar langsung dari peternakan",
			Brand:            "Fresh Milk",
			HargaBeli:        15000,
			HargaJual:        18000,
			MarginPersen:     20,
			StokMinimal:      10,
			StokMaksimal:     50,
			StokCurrent:      30,
			IsPerishable:     true,
			ShelfLife:        7,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        2,
			UpdatedBy:        2,
		},
		{
			KoperasiID:       1,
			KategoriProdukID: 3,
			SatuanProdukID:   10,
			KodeProduk:       "PRD000100004",
			NamaProduk:       "Ayam Kampung",
			Deskripsi:        "Ayam kampung segar ukuran 1-1.5kg",
			Brand:            "Ternak Lokal",
			BeratBersih:      1.2,
			HargaBeli:        35000,
			HargaJual:        45000,
			MarginPersen:     28.57,
			StokMinimal:      5,
			StokMaksimal:     30,
			StokCurrent:      15,
			IsPerishable:     true,
			ShelfLife:        2,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
		{
			KoperasiID:       1,
			KategoriProdukID: 5,
			SatuanProdukID:   1,
			KodeProduk:       "PRD000100005",
			NamaProduk:       "Mangga Harum Manis",
			Deskripsi:        "Mangga harum manis matang pohon",
			Brand:            "Buah Lokal",
			HargaBeli:        8000,
			HargaJual:        12000,
			MarginPersen:     50,
			StokMinimal:      10,
			StokMaksimal:     100,
			StokCurrent:      60,
			IsPerishable:     true,
			ShelfLife:        5,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
		{
			KoperasiID:       1,
			KategoriProdukID: 10,
			SatuanProdukID:   2,
			KodeProduk:       "PRD000100006",
			NamaProduk:       "Cabai Merah Keriting",
			Deskripsi:        "Cabai merah keriting segar dan pedas",
			Brand:            "Rempah Nusantara",
			HargaBeli:        25000,
			HargaJual:        35000,
			MarginPersen:     40,
			StokMinimal:      2,
			StokMaksimal:     20,
			StokCurrent:      10,
			IsPerishable:     true,
			ShelfLife:        7,
			IsActive:         true,
			IsReadyStock:     true,
			CreatedBy:        1,
			UpdatedBy:        1,
		},
	}

	for _, p := range produk {
		db.FirstOrCreate(&p, postgres.Produk{KodeProduk: p.KodeProduk})
	}
	fmt.Println("✓ Seeded Produk")
}

func timePtr(t time.Time) *time.Time {
	return &t
}