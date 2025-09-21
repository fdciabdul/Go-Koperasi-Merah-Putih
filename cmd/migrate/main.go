package main

import (
	"flag"
	"fmt"
	"log"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/database"
	"koperasi-merah-putih/internal/models/postgres"

	"gorm.io/gorm"
)

func main() {
	var (
		drop  = flag.Bool("drop", false, "Drop all tables before migration")
		seed  = flag.Bool("seed", false, "Run seeders after migration")
		fresh = flag.Bool("fresh", false, "Drop all tables, migrate, and seed")
	)
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewPostgresConnection(&cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if *fresh {
		*drop = true
		*seed = true
	}

	if *drop {
		fmt.Println("Dropping all tables...")
		dropAllTables(db.DB)
	}

	fmt.Println("Running migrations...")
	if err := runMigrations(db.DB); err != nil {
		log.Fatal("Migration failed:", err)
	}

	if *seed {
		fmt.Println("Running seeders...")
		runSeeders(db.DB)
	}

	fmt.Println("✓ Migration completed successfully!")
}

func runMigrations(db *gorm.DB) error {
	models := []interface{}{
		// System & Tenant
		&postgres.Tenant{},
		&postgres.User{},

		// Wilayah
		&postgres.WilayahProvinsi{},
		&postgres.WilayahKabupaten{},
		&postgres.WilayahKecamatan{},
		&postgres.WilayahKelurahan{},

		// Master Data
		&postgres.KBLI{},
		&postgres.JenisKoperasi{},
		&postgres.BentukKoperasi{},
		&postgres.StatusKoperasi{},

		// Koperasi
		&postgres.Koperasi{},
		&postgres.KoperasiAktivitasUsaha{},
		&postgres.AnggotaKoperasi{},
		&postgres.ModalKoperasi{},

		// Financial
		&postgres.COAKategori{},
		&postgres.COAAkun{},
		&postgres.JurnalUmum{},
		&postgres.JurnalDetail{},

		// Simpan Pinjam
		&postgres.ProdukSimpanPinjam{},
		&postgres.RekeningSimpanPinjam{},
		&postgres.TransaksiSimpanPinjam{},

		// Klinik
		&postgres.KlinikTenagaMedis{},
		&postgres.KlinikPasien{},
		&postgres.KlinikObat{},
		&postgres.KlinikKunjungan{},
		&postgres.KlinikResep{},

		// Product Management
		&postgres.KategoriProduk{},
		&postgres.SatuanProduk{},
		&postgres.Supplier{},
		&postgres.Produk{},
		&postgres.SupplierProduk{},
		&postgres.PurchaseOrder{},
		&postgres.PurchaseOrderDetail{},
		&postgres.PembelianHeader{},
		&postgres.PembelianDetail{},
		&postgres.PembayaranPembelian{},
		&postgres.PenjualanHeader{},
		&postgres.PenjualanDetail{},
		&postgres.StokMovement{},
		&postgres.ProdukDiskon{},

		// PPOB
		&postgres.PPOBKategori{},
		&postgres.PPOBProduk{},
		&postgres.PPOBTransaksi{},
		&postgres.PPOBSettlement{},

		// Payment
		&postgres.PaymentTransaction{},
		&postgres.SimpananPokokConfig{},

		// System
		&postgres.SequenceNumber{},
		&postgres.AuditLog{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
		fmt.Printf("✓ Migrated %T\n", model)
	}

	createCustomIndexes(db)
	createCustomConstraints(db)

	return nil
}

func dropAllTables(db *gorm.DB) {
	tables := []string{
		"audit_logs",
		"sequences",
		"simpanan_pokok_configs",
		"payment_transactions",
		"ppob_settlements",
		"ppob_transaksis",
		"ppob_produks",
		"ppob_kategoris",
		"produk_diskons",
		"stok_movements",
		"penjualan_details",
		"penjualan_headers",
		"pembayaran_pembelians",
		"pembelian_details",
		"pembelian_headers",
		"purchase_order_details",
		"purchase_orders",
		"supplier_produks",
		"produks",
		"suppliers",
		"satuan_produks",
		"kategori_produks",
		"reseps",
		"kunjungans",
		"obats",
		"pasiens",
		"tenaga_medis",
		"angsuran_pinjamen",
		"pinjamen",
		"transaksi_simpan_pinjams",
		"rekening_simpan_pinjams",
		"produk_simpan_pinjams",
		"jurnal_details",
		"jurnal_umums",
		"coa_akuns",
		"coa_kategoris",
		"modal_koperasis",
		"anggota_koperasis",
		"koperasi_aktivitas_usahas",
		"koperasis",
		"status_koperasis",
		"bentuk_koperasis",
		"jenis_koperasis",
		"kblis",
		"kelurahans",
		"kecamatans",
		"kabupatens",
		"provinsis",
		"users",
		"tenants",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		} else {
			fmt.Printf("✓ Dropped table %s\n", table)
		}
	}
}

func createCustomIndexes(db *gorm.DB) {
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
		"CREATE INDEX IF NOT EXISTS idx_koperasi_tenant_id ON koperasis(tenant_id)",
		"CREATE INDEX IF NOT EXISTS idx_anggota_koperasi_id ON anggota_koperasis(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_anggota_nik ON anggota_koperasis(nik)",
		"CREATE INDEX IF NOT EXISTS idx_coa_akun_koperasi_id ON coa_akuns(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_jurnal_umum_koperasi_id ON jurnal_umums(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_jurnal_detail_jurnal_id ON jurnal_details(jurnal_id)",
		"CREATE INDEX IF NOT EXISTS idx_jurnal_detail_akun_id ON jurnal_details(akun_id)",
		"CREATE INDEX IF NOT EXISTS idx_produk_koperasi_id ON produks(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_produk_kategori_id ON produks(kategori_produk_id)",
		"CREATE INDEX IF NOT EXISTS idx_produk_barcode ON produks(barcode)",
		"CREATE INDEX IF NOT EXISTS idx_stok_movement_produk_id ON stok_movements(produk_id)",
		"CREATE INDEX IF NOT EXISTS idx_penjualan_header_koperasi_id ON penjualan_headers(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_penjualan_header_tanggal ON penjualan_headers(tanggal_transaksi)",
		"CREATE INDEX IF NOT EXISTS idx_pembelian_header_koperasi_id ON pembelian_headers(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_rekening_simpan_pinjam_koperasi_id ON rekening_simpan_pinjams(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_rekening_simpan_pinjam_anggota_id ON rekening_simpan_pinjams(anggota_id)",
		"CREATE INDEX IF NOT EXISTS idx_pinjaman_koperasi_id ON pinjamen(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_pinjaman_anggota_id ON pinjamen(anggota_id)",
		"CREATE INDEX IF NOT EXISTS idx_pasien_koperasi_id ON pasiens(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_kunjungan_pasien_id ON kunjungans(pasien_id)",
		"CREATE INDEX IF NOT EXISTS idx_kunjungan_tanggal ON kunjungans(tanggal_kunjungan)",
		"CREATE INDEX IF NOT EXISTS idx_ppob_transaksi_koperasi_id ON ppob_transaksis(koperasi_id)",
		"CREATE INDEX IF NOT EXISTS idx_ppob_transaksi_tanggal ON ppob_transaksis(tanggal_transaksi)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_table_name ON audit_logs(table_name)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			log.Printf("Failed to create index: %v", err)
		}
	}
	fmt.Println("✓ Created custom indexes")
}

func createCustomConstraints(db *gorm.DB) {
	constraints := []string{
		"ALTER TABLE koperasis ADD CONSTRAINT check_nik_length CHECK (LENGTH(CAST(nik AS TEXT)) = 16)",
		"ALTER TABLE anggota_koperasis ADD CONSTRAINT check_nik_length CHECK (LENGTH(nik) = 16)",
		"ALTER TABLE users ADD CONSTRAINT check_nik_length CHECK (nik IS NULL OR LENGTH(nik) = 16)",
		"ALTER TABLE wilayah_kelurahans ADD CONSTRAINT check_jenis CHECK (jenis IN ('kelurahan', 'desa'))",
		"ALTER TABLE jabatan_koperasis ADD CONSTRAINT check_tingkat CHECK (tingkat IN ('pengurus', 'pengawas', 'anggota'))",
		"ALTER TABLE anggota_koperasis ADD CONSTRAINT check_jenis_kelamin CHECK (jenis_kelamin IN ('L', 'P'))",
		"ALTER TABLE anggota_koperasis ADD CONSTRAINT check_posisi CHECK (posisi IN ('pengurus', 'pengawas', 'anggota'))",
		"ALTER TABLE anggota_koperasis ADD CONSTRAINT check_status_anggota CHECK (status_anggota IN ('aktif', 'non_aktif', 'keluar'))",
		"ALTER TABLE users ADD CONSTRAINT check_role CHECK (role IN ('super_admin', 'admin_koperasi', 'bendahara', 'sekretaris', 'operator', 'anggota'))",
		"ALTER TABLE role_permissions ADD CONSTRAINT check_role CHECK (role IN ('super_admin', 'admin_koperasi', 'bendahara', 'sekretaris', 'operator', 'anggota'))",
		"ALTER TABLE koperasi_aktivitas_usahas ADD CONSTRAINT check_jenis_usaha CHECK (jenis_usaha IN ('utama', 'sampingan'))",
		"ALTER TABLE modal_koperasis ADD CONSTRAINT check_jenis_modal CHECK (jenis_modal IN ('simpanan_pokok', 'simpanan_wajib', 'dana_cadangan', 'dana_hibah', 'modal_penyertaan'))",
		"ALTER TABLE coa_kategoris ADD CONSTRAINT check_tipe CHECK (tipe IN ('aset', 'kewajiban', 'ekuitas', 'pendapatan', 'beban'))",
		"ALTER TABLE coa_akuns ADD CONSTRAINT check_saldo_normal CHECK (saldo_normal IN ('debit', 'kredit'))",
		"ALTER TABLE jurnal_umums ADD CONSTRAINT check_status CHECK (status IN ('draft', 'posted', 'cancelled'))",
		"ALTER TABLE produk_simpan_pinjams ADD CONSTRAINT check_jenis CHECK (jenis IN ('simpanan', 'pinjaman'))",
		"ALTER TABLE rekening_simpan_pinjams ADD CONSTRAINT check_status CHECK (status IN ('aktif', 'lunas', 'macet', 'tutup'))",
		"ALTER TABLE transaksi_simpan_pinjams ADD CONSTRAINT check_jenis_transaksi CHECK (jenis_transaksi IN ('setoran', 'penarikan', 'pencairan', 'angsuran', 'bunga', 'denda'))",
		"ALTER TABLE klinik_tenaga_medis ADD CONSTRAINT check_jenis_kelamin CHECK (jenis_kelamin IN ('L', 'P'))",
		"ALTER TABLE klinik_tenaga_medis ADD CONSTRAINT check_status CHECK (status IN ('aktif', 'non_aktif', 'cuti'))",
		"ALTER TABLE klinik_pasiens ADD CONSTRAINT check_jenis_kelamin CHECK (jenis_kelamin IN ('L', 'P'))",
		"ALTER TABLE klinik_pasiens ADD CONSTRAINT check_golongan_darah CHECK (golongan_darah IN ('A', 'B', 'AB', 'O', '-'))",
		"ALTER TABLE klinik_kunjungans ADD CONSTRAINT check_status_pembayaran CHECK (status_pembayaran IN ('belum_bayar', 'lunas', 'cicil'))",
		"ALTER TABLE jurnal_umums ADD CONSTRAINT check_balance CHECK (total_debit = total_kredit)",
		"ALTER TABLE angsuran_pinjamen ADD CONSTRAINT check_angsuran_positive CHECK (angsuran_pokok >= 0 AND angsuran_bunga >= 0)",
		"ALTER TABLE produk_simpan_pinjams ADD CONSTRAINT check_suku_bunga CHECK (suku_bunga >= 0 AND suku_bunga <= 100)",
		"ALTER TABLE produks ADD CONSTRAINT check_harga_positive CHECK (harga_jual > 0)",
		"ALTER TABLE produks ADD CONSTRAINT check_stok_non_negative CHECK (stok_current >= 0)",
	}

	for _, constraint := range constraints {
		if err := db.Exec(constraint).Error; err != nil {
			log.Printf("Failed to create constraint: %v", err)
		}
	}
	fmt.Println("✓ Created custom constraints")
}

func runSeeders(db *gorm.DB) {
	fmt.Println("Running seeders...")
	fmt.Println("Please run: go run cmd/seeder/main.go")
}