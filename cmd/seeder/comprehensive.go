package main

import (
	"fmt"
	"log"
	"time"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/database"
	"koperasi-merah-putih/internal/models/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	dbManager, err := database.NewPostgresConnection(&cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db := dbManager.DB

	fmt.Println("Starting comprehensive seeder...")

	// 1. Create basic tenant
	tenant := &postgres.Tenant{
		TenantCode: "DEMO",
		TenantName: "Demo Tenant",
		Domain:     "demo.local",
		IsActive:   true,
	}
	db.FirstOrCreate(tenant, postgres.Tenant{TenantCode: "DEMO"})
	fmt.Println("✓ Created tenant")

	// 2. Create admin user
	user := &postgres.User{
		TenantID:     tenant.ID,
		Username:     "admin",
		Email:        "admin@demo.local",
		PasswordHash: "$2a$10$rWjnkQr3.yUnqBv0kWpGzO7K4BF4vMEKgBqVmWxvQnRjTzJF5oUXa", // password: admin123
		NamaLengkap:  "Administrator",
		Role:         "super_admin",
		IsActive:     true,
	}
	db.FirstOrCreate(user, postgres.User{Email: "admin@demo.local"})
	fmt.Println("✓ Created admin user")

	// 3. Create sequence numbers
	sequences := []postgres.SequenceNumber{
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "global", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "koperasi", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "anggota", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "transaksi", CurrentNumber: 1000},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "jurnal", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "purchase_order", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "pembelian", CurrentNumber: 1},
		{TenantID: tenant.ID, KoperasiID: 0, SequenceName: "penjualan", CurrentNumber: 1},
	}
	for _, seq := range sequences {
		db.FirstOrCreate(&seq, postgres.SequenceNumber{
			TenantID:     seq.TenantID,
			KoperasiID:   seq.KoperasiID,
			SequenceName: seq.SequenceName,
		})
	}
	fmt.Println("✓ Created sequences")

	// 4. Create wilayah data (basic Indonesia data)
	provinsi := &postgres.WilayahProvinsi{
		ID:   11,
		Kode: "11",
		Nama: "Aceh",
	}
	db.FirstOrCreate(provinsi, postgres.WilayahProvinsi{Kode: "11"})

	kabupaten := &postgres.WilayahKabupaten{
		ID:         1101,
		Kode:       "1101",
		Nama:       "Aceh Selatan",
		ProvinsiID: provinsi.ID,
	}
	db.FirstOrCreate(kabupaten, postgres.WilayahKabupaten{Kode: "1101"})

	kecamatan := &postgres.WilayahKecamatan{
		ID:          110101,
		Kode:        "110101",
		Nama:        "Trumon",
		KabupatenID: kabupaten.ID,
	}
	db.FirstOrCreate(kecamatan, postgres.WilayahKecamatan{Kode: "110101"})

	kelurahan := &postgres.WilayahKelurahan{
		ID:          1101011001,
		Kode:        "1101011001",
		Nama:        "Trumon Timur",
		KecamatanID: kecamatan.ID,
		Jenis:       "desa",
	}
	db.FirstOrCreate(kelurahan, postgres.WilayahKelurahan{Kode: "1101011001"})
	fmt.Println("✓ Created wilayah data")

	// 5. Create koperasi master data
	jenisKoperasi := &postgres.JenisKoperasi{
		Kode:      "KP",
		Nama:      "Koperasi Primer",
		Deskripsi: "Koperasi yang dibentuk oleh orang-seorang",
		IsActive:  true,
	}
	db.FirstOrCreate(jenisKoperasi, postgres.JenisKoperasi{Kode: "KP"})

	bentukKoperasi := &postgres.BentukKoperasi{
		Kode:      "UNIT",
		Nama:      "Unit Koperasi",
		Deskripsi: "Unit koperasi terbentuk",
		IsActive:  true,
	}
	db.FirstOrCreate(bentukKoperasi, postgres.BentukKoperasi{Kode: "UNIT"})

	statusKoperasi := &postgres.StatusKoperasi{
		Kode:      "AKTIF",
		Nama:      "Aktif",
		Deskripsi: "Koperasi dalam status aktif",
		IsActive:  true,
	}
	db.FirstOrCreate(statusKoperasi, postgres.StatusKoperasi{Kode: "AKTIF"})

	kbli := &postgres.KBLI{
		Kode:      "47911",
		Nama:      "Perdagangan Eceran Melalui Pesanan Pos atau Internet",
		Kategori:  "Perdagangan",
		Deskripsi: "Kegiatan perdagangan eceran",
		IsActive:  true,
	}
	db.FirstOrCreate(kbli, postgres.KBLI{Kode: "47911"})
	fmt.Println("✓ Created koperasi master data")

	// 6. Create main koperasi
	now := time.Now()
	koperasi := &postgres.Koperasi{
		TenantID:         tenant.ID,
		NomorSK:          "001/SK/DEMO/2024",
		NIK:              1234567890123456,
		NamaKoperasi:     "Koperasi Demo Utama",
		NamaSK:           "Koperasi Demo Utama",
		JenisKoperasiID:  jenisKoperasi.ID,
		BentukKoperasiID: bentukKoperasi.ID,
		StatusKoperasiID: statusKoperasi.ID,
		ProvinsiID:       provinsi.ID,
		KabupatenID:      kabupaten.ID,
		KecamatanID:      kecamatan.ID,
		KelurahanID:      kelurahan.ID,
		Alamat:           "Jl. Demo No. 123",
		RT:               "001",
		RW:               "002",
		KodePos:          "23711",
		Email:            "koperasi@demo.local",
		Telepon:          "0812345678",
		Website:          "https://demo.local",
		TanggalBerdiri:   &now,
		TanggalSK:        &now,
		TanggalPengesahan: &now,
		CreatedBy:        user.ID,
		UpdatedBy:        user.ID,
	}
	db.FirstOrCreate(koperasi, postgres.Koperasi{NomorSK: "001/SK/DEMO/2024"})
	fmt.Println("✓ Created koperasi")

	// 7. Create koperasi aktivitas usaha
	aktivitasUsaha := &postgres.KoperasiAktivitasUsaha{
		KoperasiID: koperasi.ID,
		KBLIID:     kbli.ID,
		JenisUsaha: "utama",
		Keterangan: "Aktivitas utama koperasi",
		IsActive:   true,
	}
	db.FirstOrCreate(aktivitasUsaha, postgres.KoperasiAktivitasUsaha{
		KoperasiID: koperasi.ID,
		KBLIID:     kbli.ID,
	})
	fmt.Println("✓ Created aktivitas usaha")

	// 8. Create jabatan koperasi
	jabatan := &postgres.JabatanKoperasi{
		Kode:     "KETUA",
		Nama:     "Ketua Koperasi",
		Tingkat:  "pengurus",
		Urutan:   1,
		IsActive: true,
	}
	db.FirstOrCreate(jabatan, postgres.JabatanKoperasi{Kode: "KETUA"})

	jabatanAnggota := &postgres.JabatanKoperasi{
		Kode:     "ANGGOTA",
		Nama:     "Anggota Biasa",
		Tingkat:  "anggota",
		Urutan:   99,
		IsActive: true,
	}
	db.FirstOrCreate(jabatanAnggota, postgres.JabatanKoperasi{Kode: "ANGGOTA"})
	fmt.Println("✓ Created jabatan koperasi")

	// 9. Create anggota koperasi
	tanggalLahir := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	anggota := &postgres.AnggotaKoperasi{
		KoperasiID:    koperasi.ID,
		NIAK:          "DEMO001",
		NIK:           "1234567890123456",
		Nama:          "Budi Santoso",
		JenisKelamin:  "L",
		TempatLahir:   "Jakarta",
		TanggalLahir:  &tanggalLahir,
		Alamat:        "Jl. Anggota No. 1",
		RT:            "001",
		RW:            "001",
		KelurahanID:   kelurahan.ID,
		Telepon:       "081234567890",
		Email:         "budi@demo.local",
		Posisi:        "ketua",
		JabatanID:     jabatan.ID,
		TanggalMasuk:  &now,
		StatusAnggota: "aktif",
		Pekerjaan:     "Wiraswasta",
		Pendidikan:    "S1",
	}
	db.FirstOrCreate(anggota, postgres.AnggotaKoperasi{NIAK: "DEMO001"})

	anggota2 := &postgres.AnggotaKoperasi{
		KoperasiID:    koperasi.ID,
		NIAK:          "DEMO002",
		NIK:           "1234567890123457",
		Nama:          "Siti Aminah",
		JenisKelamin:  "P",
		TempatLahir:   "Bandung",
		TanggalLahir:  &tanggalLahir,
		Alamat:        "Jl. Anggota No. 2",
		RT:            "002",
		RW:            "001",
		KelurahanID:   kelurahan.ID,
		Telepon:       "081234567891",
		Email:         "siti@demo.local",
		Posisi:        "anggota",
		JabatanID:     jabatanAnggota.ID,
		TanggalMasuk:  &now,
		StatusAnggota: "aktif",
		Pekerjaan:     "Pedagang",
		Pendidikan:    "SMA",
	}
	db.FirstOrCreate(anggota2, postgres.AnggotaKoperasi{NIAK: "DEMO002"})
	fmt.Println("✓ Created anggota koperasi")

	// 10. Create COA categories and accounts
	coaKategori := &postgres.COAKategori{
		Kode:  "1",
		Nama:  "Aktiva",
		Tipe:  "asset",
		Urutan: 1,
	}
	db.FirstOrCreate(coaKategori, postgres.COAKategori{Kode: "1"})

	coaKategori2 := &postgres.COAKategori{
		Kode:  "4",
		Nama:  "Pendapatan",
		Tipe:  "revenue",
		Urutan: 4,
	}
	db.FirstOrCreate(coaKategori2, postgres.COAKategori{Kode: "4"})

	coaAkun := &postgres.COAAkun{
		TenantID:    tenant.ID,
		KoperasiID:  koperasi.ID,
		KodeAkun:    "1-1001",
		NamaAkun:    "Kas",
		KategoriID:  coaKategori.ID,
		LevelAkun:   1,
		SaldoNormal: "debit",
		IsKas:       true,
		IsAktif:     true,
	}
	db.FirstOrCreate(coaAkun, postgres.COAAkun{KodeAkun: "1-1001", KoperasiID: koperasi.ID})

	coaAkun2 := &postgres.COAAkun{
		TenantID:    tenant.ID,
		KoperasiID:  koperasi.ID,
		KodeAkun:    "4-4001",
		NamaAkun:    "Pendapatan Penjualan",
		KategoriID:  coaKategori2.ID,
		LevelAkun:   1,
		SaldoNormal: "kredit",
		IsKas:       false,
		IsAktif:     true,
	}
	db.FirstOrCreate(coaAkun2, postgres.COAAkun{KodeAkun: "4-4001", KoperasiID: koperasi.ID})
	fmt.Println("✓ Created COA data")

	// 11. Create modal koperasi
	modal := &postgres.ModalKoperasi{
		KoperasiID:        koperasi.ID,
		JenisModal:        "simpanan_pokok",
		Jumlah:            10000000,
		Keterangan:        "Modal awal dari simpanan pokok anggota",
		TanggalPencatatan: &now,
	}
	db.FirstOrCreate(modal, postgres.ModalKoperasi{
		KoperasiID: koperasi.ID,
		JenisModal: "simpanan_pokok",
	})
	fmt.Println("✓ Created modal koperasi")

	// 12. Create jurnal umum
	jurnal := &postgres.JurnalUmum{
		TenantID:         tenant.ID,
		KoperasiID:       koperasi.ID,
		NomorJurnal:      "JU-001",
		TanggalTransaksi: now,
		Referensi:        "Jurnal Opening Balance",
		Keterangan:       "Saldo awal kas",
		TotalDebit:       10000000,
		TotalKredit:      10000000,
		Status:           "posted",
		CreatedBy:        user.ID,
		PostedAt:         &now,
		PostedBy:         user.ID,
	}
	db.FirstOrCreate(jurnal, postgres.JurnalUmum{NomorJurnal: "JU-001", KoperasiID: koperasi.ID})

	jurnalDetail1 := &postgres.JurnalDetail{
		JurnalID:   jurnal.ID,
		AkunID:     coaAkun.ID,
		Keterangan: "Saldo awal kas",
		Debit:      10000000,
		Kredit:     0,
	}
	db.FirstOrCreate(jurnalDetail1, postgres.JurnalDetail{
		JurnalID: jurnal.ID,
		AkunID:   coaAkun.ID,
	})

	jurnalDetail2 := &postgres.JurnalDetail{
		JurnalID:   jurnal.ID,
		AkunID:     coaAkun2.ID,
		Keterangan: "Modal awal",
		Debit:      0,
		Kredit:     10000000,
	}
	db.FirstOrCreate(jurnalDetail2, postgres.JurnalDetail{
		JurnalID: jurnal.ID,
		AkunID:   coaAkun2.ID,
	})
	fmt.Println("✓ Created jurnal umum")

	// 13. Create produk simpan pinjam
	produkSimpan := &postgres.ProdukSimpanPinjam{
		KoperasiID:       koperasi.ID,
		KodeProduk:       "SP001",
		NamaProduk:       "Simpanan Sukarela",
		Jenis:            "simpanan",
		Kategori:         "sukarela",
		BungaSimpanan:    3.5,
		MinimalSaldo:     50000,
		SyaratKetentuan:  "Minimal setoran 50.000",
		IsAktif:          true,
	}
	db.FirstOrCreate(produkSimpan, postgres.ProdukSimpanPinjam{KodeProduk: "SP001", KoperasiID: koperasi.ID})

	produkPinjam := &postgres.ProdukSimpanPinjam{
		KoperasiID:       koperasi.ID,
		KodeProduk:       "PJ001",
		NamaProduk:       "Pinjaman Produktif",
		Jenis:            "pinjaman",
		Kategori:         "produktif",
		BungaPinjaman:    12.0,
		BungaDenda:       2.0,
		MaksimalPinjaman: 50000000,
		JangkaWaktuMax:   36,
		SyaratKetentuan:  "Untuk kegiatan produktif",
		IsAktif:          true,
	}
	db.FirstOrCreate(produkPinjam, postgres.ProdukSimpanPinjam{KodeProduk: "PJ001", KoperasiID: koperasi.ID})
	fmt.Println("✓ Created produk simpan pinjam")

	// 14. Create rekening simpan pinjam
	rekening := &postgres.RekeningSimpanPinjam{
		KoperasiID:    koperasi.ID,
		AnggotaID:     anggota.ID,
		ProdukID:      produkSimpan.ID,
		NomorRekening: "SV-001-0001",
		SaldoSimpanan: 1000000,
		Status:        "aktif",
		TanggalBuka:   now,
	}
	db.FirstOrCreate(rekening, postgres.RekeningSimpanPinjam{
		NomorRekening: "SV-001-0001",
		KoperasiID:    koperasi.ID,
	})
	fmt.Println("✓ Created rekening simpan pinjam")

	// 15. Create kategori produk dan satuan
	kategoriProduk := &postgres.KategoriProduk{
		Kode:      "MKN",
		Nama:      "Makanan",
		Deskripsi: "Kategori produk makanan",
		Icon:      "food",
		IsActive:  true,
	}
	db.FirstOrCreate(kategoriProduk, postgres.KategoriProduk{Kode: "MKN"})

	satuanProduk := &postgres.SatuanProduk{
		Kode:     "PCS",
		Nama:     "Pieces",
		IsActive: true,
	}
	db.FirstOrCreate(satuanProduk, postgres.SatuanProduk{Kode: "PCS"})
	fmt.Println("✓ Created kategori dan satuan produk")

	// 16. Create supplier
	supplier := &postgres.Supplier{
		KoperasiID:     koperasi.ID,
		Kode:           "SUP001",
		Nama:           "PT Supplier Demo",
		KontakPerson:   "Andi Supplier",
		Telepon:        "0211234567",
		Email:          "supplier@demo.local",
		Alamat:         "Jl. Supplier No. 123",
		ProvinsiID:     provinsi.ID,
		KabupatenID:    kabupaten.ID,
		NoRekening:     "1234567890",
		NamaBank:       "Bank Demo",
		AtasNamaBank:   "PT Supplier Demo",
		JenisSupplier:  "perusahaan",
		Status:         "aktif",
		TermPembayaran: 30,
		IsActive:       true,
		CreatedBy:      user.ID,
		UpdatedBy:      user.ID,
	}
	db.FirstOrCreate(supplier, postgres.Supplier{Kode: "SUP001", KoperasiID: koperasi.ID})
	fmt.Println("✓ Created supplier")

	// 17. Create produk
	produk := &postgres.Produk{
		KoperasiID:       koperasi.ID,
		KategoriProdukID: kategoriProduk.ID,
		SatuanProdukID:   satuanProduk.ID,
		KodeProduk:       "PRD001",
		Barcode:          "1234567890123",
		NamaProduk:       "Beras Premium 5kg",
		Deskripsi:        "Beras premium kualitas terbaik",
		Brand:            "Demo Brand",
		Varian:           "5kg",
		BeratBersih:      5.0,
		HargaBeli:        45000,
		HargaJual:        50000,
		MarginPersen:     11.11,
		StokMinimal:      10,
		StokMaksimal:     100,
		StokCurrent:      50,
		IsPerishable:     false,
		IsProduksi:       false,
		IsActive:         true,
		IsReadyStock:     true,
		CreatedBy:        user.ID,
		UpdatedBy:        user.ID,
	}
	db.FirstOrCreate(produk, postgres.Produk{KodeProduk: "PRD001", KoperasiID: koperasi.ID})

	produk2 := &postgres.Produk{
		KoperasiID:       koperasi.ID,
		KategoriProdukID: kategoriProduk.ID,
		SatuanProdukID:   satuanProduk.ID,
		KodeProduk:       "PRD002",
		Barcode:          "1234567890124",
		NamaProduk:       "Minyak Goreng 1L",
		Deskripsi:        "Minyak goreng berkualitas",
		Brand:            "Demo Oil",
		Varian:           "1L",
		BeratBersih:      1.0,
		HargaBeli:        12000,
		HargaJual:        15000,
		MarginPersen:     25.0,
		StokMinimal:      20,
		StokMaksimal:     200,
		StokCurrent:      100,
		IsPerishable:     false,
		IsProduksi:       false,
		IsActive:         true,
		IsReadyStock:     true,
		CreatedBy:        user.ID,
		UpdatedBy:        user.ID,
	}
	db.FirstOrCreate(produk2, postgres.Produk{KodeProduk: "PRD002", KoperasiID: koperasi.ID})
	fmt.Println("✓ Created produk")

	// 18. Create supplier produk
	supplierProduk := &postgres.SupplierProduk{
		SupplierID:    supplier.ID,
		ProdukID:      produk.ID,
		KodeSupplier:  "SUP-PRD001",
		HargaSupplier: 45000,
		MinOrder:      10,
		LeadTime:      7,
		IsPreferred:   true,
		IsActive:      true,
	}
	db.FirstOrCreate(supplierProduk, postgres.SupplierProduk{
		SupplierID: supplier.ID,
		ProdukID:   produk.ID,
	})
	fmt.Println("✓ Created supplier produk")

	// 19. Create purchase order
	purchaseOrder := &postgres.PurchaseOrder{
		KoperasiID:   koperasi.ID,
		SupplierID:   supplier.ID,
		NomorPO:      "PO-001",
		TanggalPO:    now,
		TotalItem:    2,
		SubTotal:     570000,
		PajakPersen:  0,
		TotalPajak:   0,
		BiayaKirim:   0,
		Diskon:       0,
		GrandTotal:   570000,
		Status:       "approved",
		Keterangan:   "Purchase order demo",
		ApprovedBy:   user.ID,
		ApprovedAt:   &now,
		CreatedBy:    user.ID,
		UpdatedBy:    user.ID,
	}
	db.FirstOrCreate(purchaseOrder, postgres.PurchaseOrder{NomorPO: "PO-001", KoperasiID: koperasi.ID})

	poDetail1 := &postgres.PurchaseOrderDetail{
		PurchaseOrderID: purchaseOrder.ID,
		ProdukID:        produk.ID,
		Qty:             10,
		HargaSatuan:     45000,
		Subtotal:        450000,
		QtyReceived:     10,
		Keterangan:      "Beras premium",
	}
	db.FirstOrCreate(poDetail1, postgres.PurchaseOrderDetail{
		PurchaseOrderID: purchaseOrder.ID,
		ProdukID:        produk.ID,
	})

	poDetail2 := &postgres.PurchaseOrderDetail{
		PurchaseOrderID: purchaseOrder.ID,
		ProdukID:        produk2.ID,
		Qty:             10,
		HargaSatuan:     12000,
		Subtotal:        120000,
		QtyReceived:     10,
		Keterangan:      "Minyak goreng",
	}
	db.FirstOrCreate(poDetail2, postgres.PurchaseOrderDetail{
		PurchaseOrderID: purchaseOrder.ID,
		ProdukID:        produk2.ID,
	})
	fmt.Println("✓ Created purchase order")

	// 20. Create pembelian
	pembelian := &postgres.PembelianHeader{
		KoperasiID:        koperasi.ID,
		SupplierID:        supplier.ID,
		PurchaseOrderID:   purchaseOrder.ID,
		NomorFaktur:       "FB-001",
		TanggalFaktur:     now,
		TotalItem:         2,
		SubTotal:          570000,
		PajakPersen:       0,
		TotalPajak:        0,
		BiayaKirim:        0,
		Diskon:            0,
		GrandTotal:        570000,
		StatusPembayaran:  "paid",
		TotalBayar:        570000,
		Keterangan:        "Pembelian barang demo",
		CreatedBy:         user.ID,
		UpdatedBy:         user.ID,
	}
	db.FirstOrCreate(pembelian, postgres.PembelianHeader{NomorFaktur: "FB-001", KoperasiID: koperasi.ID})

	pembelianDetail1 := &postgres.PembelianDetail{
		PembelianHeaderID: pembelian.ID,
		ProdukID:          produk.ID,
		Qty:               10,
		HargaSatuan:       45000,
		Subtotal:          450000,
		Keterangan:        "Beras premium",
	}
	db.FirstOrCreate(pembelianDetail1, postgres.PembelianDetail{
		PembelianHeaderID: pembelian.ID,
		ProdukID:          produk.ID,
	})

	pembelianDetail2 := &postgres.PembelianDetail{
		PembelianHeaderID: pembelian.ID,
		ProdukID:          produk2.ID,
		Qty:               10,
		HargaSatuan:       12000,
		Subtotal:          120000,
		Keterangan:        "Minyak goreng",
	}
	db.FirstOrCreate(pembelianDetail2, postgres.PembelianDetail{
		PembelianHeaderID: pembelian.ID,
		ProdukID:          produk2.ID,
	})
	fmt.Println("✓ Created pembelian")

	// 21. Create penjualan
	penjualan := &postgres.PenjualanHeader{
		KoperasiID:       koperasi.ID,
		AnggotaID:        anggota.ID,
		NomorTransaksi:   "TXN-001",
		TanggalTransaksi: now,
		TotalItem:        2,
		SubTotal:         65000,
		PajakPersen:      0,
		TotalPajak:       0,
		Diskon:           0,
		GrandTotal:       65000,
		MetodePembayaran: "cash",
		StatusPembayaran: "paid",
		JumlahBayar:      70000,
		JumlahKembalian:  5000,
		Kasir:            "admin",
		Keterangan:       "Penjualan demo",
		CreatedBy:        user.ID,
		UpdatedBy:        user.ID,
	}
	db.FirstOrCreate(penjualan, postgres.PenjualanHeader{NomorTransaksi: "TXN-001", KoperasiID: koperasi.ID})

	penjualanDetail1 := &postgres.PenjualanDetail{
		PenjualanHeaderID: penjualan.ID,
		ProdukID:          produk.ID,
		Qty:               1,
		HargaSatuan:       50000,
		DiskonPersen:      0,
		DiskonRupiah:      0,
		Subtotal:          50000,
		Keterangan:        "Beras premium",
	}
	db.FirstOrCreate(penjualanDetail1, postgres.PenjualanDetail{
		PenjualanHeaderID: penjualan.ID,
		ProdukID:          produk.ID,
	})

	penjualanDetail2 := &postgres.PenjualanDetail{
		PenjualanHeaderID: penjualan.ID,
		ProdukID:          produk2.ID,
		Qty:               1,
		HargaSatuan:       15000,
		DiskonPersen:      0,
		DiskonRupiah:      0,
		Subtotal:          15000,
		Keterangan:        "Minyak goreng",
	}
	db.FirstOrCreate(penjualanDetail2, postgres.PenjualanDetail{
		PenjualanHeaderID: penjualan.ID,
		ProdukID:          produk2.ID,
	})
	fmt.Println("✓ Created penjualan")

	// 22. Create stok movement
	stokMovement1 := &postgres.StokMovement{
		KoperasiID:       koperasi.ID,
		ProdukID:         produk.ID,
		TipeMovement:     "in",
		ReferensiTipe:    "pembelian",
		ReferensiID:      pembelian.ID,
		TanggalMovement:  now,
		QtyBefore:        40,
		QtyMovement:      10,
		QtyAfter:         50,
		HargaSatuan:      45000,
		TotalNilai:       450000,
		Keterangan:       "Pembelian dari supplier",
		CreatedBy:        user.ID,
	}
	db.FirstOrCreate(stokMovement1, postgres.StokMovement{
		ProdukID:      produk.ID,
		ReferensiTipe: "pembelian",
		ReferensiID:   pembelian.ID,
	})

	stokMovement2 := &postgres.StokMovement{
		KoperasiID:       koperasi.ID,
		ProdukID:         produk.ID,
		TipeMovement:     "out",
		ReferensiTipe:    "penjualan",
		ReferensiID:      penjualan.ID,
		TanggalMovement:  now,
		QtyBefore:        50,
		QtyMovement:      1,
		QtyAfter:         49,
		HargaSatuan:      50000,
		TotalNilai:       50000,
		Keterangan:       "Penjualan ke anggota",
		CreatedBy:        user.ID,
	}
	db.FirstOrCreate(stokMovement2, postgres.StokMovement{
		ProdukID:      produk.ID,
		ReferensiTipe: "penjualan",
		ReferensiID:   penjualan.ID,
	})
	fmt.Println("✓ Created stok movement")

	// 23. Create PPOB provider data
	ppobProvider := &postgres.PPOBProvider{
		Kode:      "DEMO",
		Nama:      "Demo PPOB Provider",
		BaseURL:   "https://api.demo-ppob.com",
		APIKey:    "demo-api-key",
		SecretKey: "demo-secret-key",
		IsAktif:   true,
	}
	db.FirstOrCreate(ppobProvider, postgres.PPOBProvider{Kode: "DEMO"})

	ppobKategori := &postgres.PPOBKategori{
		Kode:    "PLN",
		Nama:    "PLN Token",
		Icon:    "electricity",
		Urutan:  1,
		IsAktif: true,
	}
	db.FirstOrCreate(ppobKategori, postgres.PPOBKategori{Kode: "PLN"})

	ppobProduk := &postgres.PPOBProduk{
		ProviderID:     ppobProvider.ID,
		KategoriID:     ppobKategori.ID,
		KodeProduk:     "PLN20",
		NamaProduk:     "PLN Token 20.000",
		Deskripsi:      "Token PLN senilai 20.000",
		HargaBeli:      19500,
		HargaJual:      20000,
		FeeAgen:        500,
		IsAktif:        true,
		ValidasiFormat: "numeric",
	}
	db.FirstOrCreate(ppobProduk, postgres.PPOBProduk{KodeProduk: "PLN20", ProviderID: ppobProvider.ID})
	fmt.Println("✓ Created PPOB data")

	// 24. Create clinic data
	dokter := &postgres.KlinikTenagaMedis{
		KoperasiID:      koperasi.ID,
		NIP:             "DOK001",
		NamaLengkap:     "Dr. Ahmad Sutrisno",
		JenisKelamin:    "L",
		Spesialisasi:    "Umum",
		NoSTR:           "STR123456",
		NoSIP:           "SIP789012",
		Telepon:         "081234567892",
		Email:           "dr.ahmad@demo.local",
		JadwalPraktik:   `{"senin": "08:00-12:00", "selasa": "08:00-12:00", "rabu": "08:00-12:00"}`,
		TarifKonsultasi: 50000,
		Status:          "aktif",
	}
	db.FirstOrCreate(dokter, postgres.KlinikTenagaMedis{NIP: "DOK001", KoperasiID: koperasi.ID})

	pasien := &postgres.KlinikPasien{
		KoperasiID:      koperasi.ID,
		NomorRM:         "RM001",
		NIK:             "1234567890123458",
		NamaLengkap:     "Ibu Sari",
		JenisKelamin:    "P",
		TempatLahir:     "Surabaya",
		TanggalLahir:    &tanggalLahir,
		Alamat:          "Jl. Pasien No. 123",
		Telepon:         "081234567893",
		Email:           "sari@demo.local",
		GolonganDarah:   "A",
		Alergi:          "Tidak ada",
		RiwayatPenyakit: "Hipertensi",
		AnggotaID:       anggota2.ID,
	}
	db.FirstOrCreate(pasien, postgres.KlinikPasien{NomorRM: "RM001", KoperasiID: koperasi.ID})
	fmt.Println("✓ Created clinic data")

	// 25. Create system settings
	systemSettings := []postgres.SystemSetting{
		{
			TenantID:    tenant.ID,
			KoperasiID:  koperasi.ID,
			KeyName:     "app_name",
			KeyValue:    "Koperasi Demo",
			DataType:    "string",
			Category:    "general",
			Description: "Nama aplikasi",
		},
		{
			TenantID:    tenant.ID,
			KoperasiID:  koperasi.ID,
			KeyName:     "timezone",
			KeyValue:    "Asia/Jakarta",
			DataType:    "string",
			Category:    "general",
			Description: "Timezone aplikasi",
		},
		{
			TenantID:    tenant.ID,
			KoperasiID:  koperasi.ID,
			KeyName:     "currency",
			KeyValue:    "IDR",
			DataType:    "string",
			Category:    "financial",
			Description: "Mata uang default",
		},
	}

	for _, setting := range systemSettings {
		db.FirstOrCreate(&setting, postgres.SystemSetting{
			TenantID:   setting.TenantID,
			KoperasiID: setting.KoperasiID,
			KeyName:    setting.KeyName,
		})
	}
	fmt.Println("✓ Created system settings")

	fmt.Println("Comprehensive seeder completed successfully!")
	fmt.Println("=== Login Credentials ===")
	fmt.Println("Email: admin@demo.local")
	fmt.Println("Password: admin123")
	fmt.Println("")
	fmt.Println("=== Demo Data Created ===")
	fmt.Println("- 1 Tenant: Demo Tenant")
	fmt.Println("- 1 Koperasi: Koperasi Demo Utama")
	fmt.Println("- 2 Anggota: Budi Santoso (Ketua), Siti Aminah (Anggota)")
	fmt.Println("- 2 Produk: Beras Premium 5kg, Minyak Goreng 1L")
	fmt.Println("- 1 Supplier: PT Supplier Demo")
	fmt.Println("- Sample transactions: Purchase Order, Pembelian, Penjualan")
	fmt.Println("- Financial data: COA, Jurnal, Modal Koperasi")
	fmt.Println("- PPOB data: PLN Token products")
	fmt.Println("- Simpan Pinjam: Savings & Loan products with accounts")
	fmt.Println("- Clinic data: Doctor and Patient records")
}