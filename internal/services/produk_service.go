package services

import (
	"fmt"
	"strconv"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	"koperasi-merah-putih/internal/repository"
)

type ProdukService struct {
	produkRepo    *repository.ProdukRepository
	sequenceRepo  *repository.SequenceRepository
}

func NewProdukService(produkRepo *repository.ProdukRepository, sequenceRepo *repository.SequenceRepository) *ProdukService {
	return &ProdukService{
		produkRepo:   produkRepo,
		sequenceRepo: sequenceRepo,
	}
}

// Request Structs
type CreateKategoriProdukRequest struct {
	Kode      string `json:"kode" binding:"required,max=20"`
	Nama      string `json:"nama" binding:"required,max=100"`
	Deskripsi string `json:"deskripsi"`
	Icon      string `json:"icon"`
}

type CreateSatuanProdukRequest struct {
	Kode string `json:"kode" binding:"required,max=10"`
	Nama string `json:"nama" binding:"required,max=50"`
}

type CreateSupplierRequest struct {
	KoperasiID     uint64 `json:"koperasi_id" binding:"required"`
	Kode           string `json:"kode" binding:"required,max=20"`
	Nama           string `json:"nama" binding:"required,max=255"`
	KontakPerson   string `json:"kontak_person"`
	Telepon        string `json:"telepon"`
	Email          string `json:"email"`
	Alamat         string `json:"alamat"`
	ProvinsiID     uint64 `json:"provinsi_id"`
	KabupatenID    uint64 `json:"kabupaten_id"`
	NoRekening     string `json:"no_rekening"`
	NamaBank       string `json:"nama_bank"`
	AtasNamaBank   string `json:"atas_nama_bank"`
	NPWP           string `json:"npwp"`
	JenisSupplier  string `json:"jenis_supplier" binding:"oneof=individu perusahaan koperasi"`
	TermPembayaran int    `json:"term_pembayaran"`
	CreatedBy      uint64 `json:"created_by"`
}

type CreateProdukRequest struct {
	KoperasiID       uint64  `json:"koperasi_id" binding:"required"`
	KategoriProdukID uint64  `json:"kategori_produk_id" binding:"required"`
	SatuanProdukID   uint64  `json:"satuan_produk_id" binding:"required"`
	NamaProduk       string  `json:"nama_produk" binding:"required,max=255"`
	Deskripsi        string  `json:"deskripsi"`
	Brand            string  `json:"brand"`
	Varian           string  `json:"varian"`
	BeratBersih      float64 `json:"berat_bersih"`
	Dimensi          string  `json:"dimensi"`
	FotoProduk       string  `json:"foto_produk"`
	HargaBeli        float64 `json:"harga_beli"`
	HargaJual        float64 `json:"harga_jual" binding:"required,gt=0"`
	MarginPersen     float64 `json:"margin_persen"`
	StokMinimal      int     `json:"stok_minimal"`
	StokMaksimal     int     `json:"stok_maksimal"`
	IsPerishable     bool    `json:"is_perishable"`
	ShelfLife        int     `json:"shelf_life"`
	IsProduksi       bool    `json:"is_produksi"`
	CreatedBy        uint64  `json:"created_by"`
}

type CreatePurchaseOrderRequest struct {
	KoperasiID  uint64                        `json:"koperasi_id" binding:"required"`
	SupplierID  uint64                        `json:"supplier_id" binding:"required"`
	TanggalPO   time.Time                     `json:"tanggal_po" binding:"required"`
	Keterangan  string                        `json:"keterangan"`
	Items       []PurchaseOrderDetailRequest  `json:"items" binding:"required,min=1"`
	CreatedBy   uint64                        `json:"created_by"`
}

type PurchaseOrderDetailRequest struct {
	ProdukID    uint64  `json:"produk_id" binding:"required"`
	Qty         int     `json:"qty" binding:"required,gt=0"`
	HargaSatuan float64 `json:"harga_satuan" binding:"required,gt=0"`
	Keterangan  string  `json:"keterangan"`
}

type CreatePembelianRequest struct {
	KoperasiID       uint64                   `json:"koperasi_id" binding:"required"`
	SupplierID       uint64                   `json:"supplier_id" binding:"required"`
	PurchaseOrderID  uint64                   `json:"purchase_order_id"`
	NomorFaktur      string                   `json:"nomor_faktur" binding:"required"`
	TanggalFaktur    time.Time                `json:"tanggal_faktur" binding:"required"`
	TanggalJatuhTempo *time.Time              `json:"tanggal_jatuh_tempo"`
	PajakPersen      float64                  `json:"pajak_persen"`
	BiayaKirim       float64                  `json:"biaya_kirim"`
	Diskon           float64                  `json:"diskon"`
	Keterangan       string                   `json:"keterangan"`
	Items            []PembelianDetailRequest `json:"items" binding:"required,min=1"`
	CreatedBy        uint64                   `json:"created_by"`
}

type PembelianDetailRequest struct {
	ProdukID       uint64     `json:"produk_id" binding:"required"`
	Qty            int        `json:"qty" binding:"required,gt=0"`
	HargaSatuan    float64    `json:"harga_satuan" binding:"required,gt=0"`
	TanggalExpired *time.Time `json:"tanggal_expired"`
	BatchNumber    string     `json:"batch_number"`
	Keterangan     string     `json:"keterangan"`
}

type CreatePenjualanRequest struct {
	KoperasiID       uint64                   `json:"koperasi_id" binding:"required"`
	AnggotaID        uint64                   `json:"anggota_id"`
	TanggalTransaksi time.Time                `json:"tanggal_transaksi" binding:"required"`
	MetodePembayaran string                   `json:"metode_pembayaran" binding:"oneof=cash debit credit transfer simpanan"`
	JumlahBayar      float64                  `json:"jumlah_bayar" binding:"required,gt=0"`
	Kasir            string                   `json:"kasir"`
	Keterangan       string                   `json:"keterangan"`
	Items            []PenjualanDetailRequest `json:"items" binding:"required,min=1"`
	CreatedBy        uint64                   `json:"created_by"`
}

type PenjualanDetailRequest struct {
	ProdukID     uint64  `json:"produk_id" binding:"required"`
	Qty          int     `json:"qty" binding:"required,gt=0"`
	HargaSatuan  float64 `json:"harga_satuan" binding:"required,gt=0"`
	DiskonPersen float64 `json:"diskon_persen"`
	DiskonRupiah float64 `json:"diskon_rupiah"`
	Keterangan   string  `json:"keterangan"`
}

// Kategori Produk Services
func (s *ProdukService) CreateKategoriProduk(req *CreateKategoriProdukRequest) (*postgres.KategoriProduk, error) {
	kategori := &postgres.KategoriProduk{
		Kode:      req.Kode,
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		Icon:      req.Icon,
		IsActive:  true,
	}

	if err := s.produkRepo.CreateKategoriProduk(kategori); err != nil {
		return nil, fmt.Errorf("failed to create kategori produk: %v", err)
	}

	return kategori, nil
}

func (s *ProdukService) GetKategoriProdukByID(id uint64) (*postgres.KategoriProduk, error) {
	return s.produkRepo.GetKategoriProdukByID(id)
}

func (s *ProdukService) GetAllKategoriProduk() ([]postgres.KategoriProduk, error) {
	return s.produkRepo.GetAllKategoriProduk()
}

// Satuan Produk Services
func (s *ProdukService) CreateSatuanProduk(req *CreateSatuanProdukRequest) (*postgres.SatuanProduk, error) {
	satuan := &postgres.SatuanProduk{
		Kode:     req.Kode,
		Nama:     req.Nama,
		IsActive: true,
	}

	if err := s.produkRepo.CreateSatuanProduk(satuan); err != nil {
		return nil, fmt.Errorf("failed to create satuan produk: %v", err)
	}

	return satuan, nil
}

func (s *ProdukService) GetAllSatuanProduk() ([]postgres.SatuanProduk, error) {
	return s.produkRepo.GetAllSatuanProduk()
}

// Supplier Services
func (s *ProdukService) CreateSupplier(req *CreateSupplierRequest) (*postgres.Supplier, error) {
	supplier := &postgres.Supplier{
		KoperasiID:     req.KoperasiID,
		Kode:           req.Kode,
		Nama:           req.Nama,
		KontakPerson:   req.KontakPerson,
		Telepon:        req.Telepon,
		Email:          req.Email,
		Alamat:         req.Alamat,
		ProvinsiID:     req.ProvinsiID,
		KabupatenID:    req.KabupatenID,
		NoRekening:     req.NoRekening,
		NamaBank:       req.NamaBank,
		AtasNamaBank:   req.AtasNamaBank,
		NPWP:           req.NPWP,
		JenisSupplier:  req.JenisSupplier,
		TermPembayaran: req.TermPembayaran,
		Status:         "aktif",
		IsActive:       true,
		CreatedBy:      req.CreatedBy,
		UpdatedBy:      req.CreatedBy,
	}

	if err := s.produkRepo.CreateSupplier(supplier); err != nil {
		return nil, fmt.Errorf("failed to create supplier: %v", err)
	}

	return supplier, nil
}

func (s *ProdukService) GetSuppliersByKoperasi(koperasiID uint64, page, limit int) ([]postgres.Supplier, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.produkRepo.GetSuppliersByKoperasi(koperasiID, limit, offset)
}

// Produk Services
func (s *ProdukService) CreateProduk(req *CreateProdukRequest) (*postgres.Produk, error) {
	sequence, err := s.sequenceRepo.GetNextNumber(1, req.KoperasiID, "produk")
	if err != nil {
		return nil, fmt.Errorf("failed to generate product code: %v", err)
	}

	kodeProduk := fmt.Sprintf("PRD%04d%06d", req.KoperasiID, sequence)

	produk := &postgres.Produk{
		KoperasiID:       req.KoperasiID,
		KategoriProdukID: req.KategoriProdukID,
		SatuanProdukID:   req.SatuanProdukID,
		KodeProduk:       kodeProduk,
		NamaProduk:       req.NamaProduk,
		Deskripsi:        req.Deskripsi,
		Brand:            req.Brand,
		Varian:           req.Varian,
		BeratBersih:      req.BeratBersih,
		Dimensi:          req.Dimensi,
		FotoProduk:       req.FotoProduk,
		HargaBeli:        req.HargaBeli,
		HargaJual:        req.HargaJual,
		MarginPersen:     req.MarginPersen,
		StokMinimal:      req.StokMinimal,
		StokMaksimal:     req.StokMaksimal,
		IsPerishable:     req.IsPerishable,
		ShelfLife:        req.ShelfLife,
		IsProduksi:       req.IsProduksi,
		IsActive:         true,
		IsReadyStock:     true,
		CreatedBy:        req.CreatedBy,
		UpdatedBy:        req.CreatedBy,
	}

	if req.MarginPersen == 0 && req.HargaBeli > 0 {
		produk.MarginPersen = ((req.HargaJual - req.HargaBeli) / req.HargaBeli) * 100
	}

	if err := s.produkRepo.CreateProduk(produk); err != nil {
		return nil, fmt.Errorf("failed to create produk: %v", err)
	}

	return produk, nil
}

func (s *ProdukService) GetProdukByID(id uint64) (*postgres.Produk, error) {
	return s.produkRepo.GetProdukByID(id)
}

func (s *ProdukService) GetProdukByBarcode(barcode string) (*postgres.Produk, error) {
	return s.produkRepo.GetProdukByBarcode(barcode)
}

func (s *ProdukService) GetProduksByKoperasi(koperasiID uint64, filters repository.ProdukFilters, page, limit int) ([]postgres.Produk, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.produkRepo.GetProduksByKoperasi(koperasiID, filters, limit, offset)
}

func (s *ProdukService) GenerateBarcode(produkID uint64) (string, error) {
	sequence, err := s.sequenceRepo.GetNextNumber(1, 0, "barcode")
	if err != nil {
		return "", fmt.Errorf("failed to generate barcode: %v", err)
	}

	barcode := fmt.Sprintf("890%010d%d", sequence, s.calculateCheckDigit(fmt.Sprintf("890%010d", sequence)))
	return barcode, nil
}

func (s *ProdukService) calculateCheckDigit(code string) int {
	sum := 0
	for i, char := range code {
		digit, _ := strconv.Atoi(string(char))
		if i%2 == 0 {
			sum += digit * 1
		} else {
			sum += digit * 3
		}
	}
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}

// Purchase Order Services
func (s *ProdukService) CreatePurchaseOrder(req *CreatePurchaseOrderRequest) (*postgres.PurchaseOrder, error) {
	sequence, err := s.sequenceRepo.GetNextNumber(1, req.KoperasiID, "purchase_order")
	if err != nil {
		return nil, fmt.Errorf("failed to generate PO number: %v", err)
	}

	nomorPO := fmt.Sprintf("PO%04d%06d", req.KoperasiID, sequence)

	var totalItem int
	var subTotal float64
	var details []postgres.PurchaseOrderDetail

	for _, item := range req.Items {
		detail := postgres.PurchaseOrderDetail{
			ProdukID:    item.ProdukID,
			Qty:         item.Qty,
			HargaSatuan: item.HargaSatuan,
			Subtotal:    float64(item.Qty) * item.HargaSatuan,
			Keterangan:  item.Keterangan,
		}
		details = append(details, detail)
		totalItem += item.Qty
		subTotal += detail.Subtotal
	}

	po := &postgres.PurchaseOrder{
		KoperasiID:          req.KoperasiID,
		SupplierID:          req.SupplierID,
		NomorPO:             nomorPO,
		TanggalPO:           req.TanggalPO,
		TotalItem:           totalItem,
		SubTotal:            subTotal,
		GrandTotal:          subTotal,
		Status:              "draft",
		Keterangan:          req.Keterangan,
		PurchaseOrderDetail: details,
		CreatedBy:           req.CreatedBy,
		UpdatedBy:           req.CreatedBy,
	}

	if err := s.produkRepo.CreatePurchaseOrder(po); err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %v", err)
	}

	return po, nil
}

func (s *ProdukService) GetPurchaseOrdersByKoperasi(koperasiID uint64, page, limit int) ([]postgres.PurchaseOrder, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.produkRepo.GetPurchaseOrdersByKoperasi(koperasiID, limit, offset)
}

// Pembelian Services
func (s *ProdukService) CreatePembelian(req *CreatePembelianRequest) (*postgres.PembelianHeader, error) {
	var totalItem int
	var subTotal float64
	var details []postgres.PembelianDetail

	for _, item := range req.Items {
		detail := postgres.PembelianDetail{
			ProdukID:       item.ProdukID,
			Qty:            item.Qty,
			HargaSatuan:    item.HargaSatuan,
			Subtotal:       float64(item.Qty) * item.HargaSatuan,
			TanggalExpired: item.TanggalExpired,
			BatchNumber:    item.BatchNumber,
			Keterangan:     item.Keterangan,
		}
		details = append(details, detail)
		totalItem += item.Qty
		subTotal += detail.Subtotal
	}

	totalPajak := subTotal * (req.PajakPersen / 100)
	grandTotal := subTotal + totalPajak + req.BiayaKirim - req.Diskon

	pembelian := &postgres.PembelianHeader{
		KoperasiID:        req.KoperasiID,
		SupplierID:        req.SupplierID,
		PurchaseOrderID:   req.PurchaseOrderID,
		NomorFaktur:       req.NomorFaktur,
		TanggalFaktur:     req.TanggalFaktur,
		TanggalJatuhTempo: req.TanggalJatuhTempo,
		TotalItem:         totalItem,
		SubTotal:          subTotal,
		PajakPersen:       req.PajakPersen,
		TotalPajak:        totalPajak,
		BiayaKirim:        req.BiayaKirim,
		Diskon:            req.Diskon,
		GrandTotal:        grandTotal,
		StatusPembayaran:  "unpaid",
		Keterangan:        req.Keterangan,
		PembelianDetail:   details,
		CreatedBy:         req.CreatedBy,
		UpdatedBy:         req.CreatedBy,
	}

	if err := s.produkRepo.CreatePembelian(pembelian); err != nil {
		return nil, fmt.Errorf("failed to create pembelian: %v", err)
	}

	return pembelian, nil
}

// Penjualan Services
func (s *ProdukService) CreatePenjualan(req *CreatePenjualanRequest) (*postgres.PenjualanHeader, error) {
	sequence, err := s.sequenceRepo.GetNextNumber(1, req.KoperasiID, "penjualan")
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction number: %v", err)
	}

	nomorTransaksi := fmt.Sprintf("TRX%04d%06d", req.KoperasiID, sequence)

	var totalItem int
	var subTotal float64
	var details []postgres.PenjualanDetail

	for _, item := range req.Items {
		subtotal := (float64(item.Qty) * item.HargaSatuan) - item.DiskonRupiah
		if item.DiskonPersen > 0 {
			subtotal = subtotal * (1 - item.DiskonPersen/100)
		}

		detail := postgres.PenjualanDetail{
			ProdukID:     item.ProdukID,
			Qty:          item.Qty,
			HargaSatuan:  item.HargaSatuan,
			DiskonPersen: item.DiskonPersen,
			DiskonRupiah: item.DiskonRupiah,
			Subtotal:     subtotal,
			Keterangan:   item.Keterangan,
		}
		details = append(details, detail)
		totalItem += item.Qty
		subTotal += detail.Subtotal
	}

	kembalian := req.JumlahBayar - subTotal
	if kembalian < 0 {
		return nil, fmt.Errorf("jumlah bayar tidak mencukupi")
	}

	penjualan := &postgres.PenjualanHeader{
		KoperasiID:       req.KoperasiID,
		AnggotaID:        req.AnggotaID,
		NomorTransaksi:   nomorTransaksi,
		TanggalTransaksi: req.TanggalTransaksi,
		TotalItem:        totalItem,
		SubTotal:         subTotal,
		GrandTotal:       subTotal,
		MetodePembayaran: req.MetodePembayaran,
		StatusPembayaran: "paid",
		JumlahBayar:      req.JumlahBayar,
		JumlahKembalian:  kembalian,
		Kasir:            req.Kasir,
		Keterangan:       req.Keterangan,
		PenjualanDetail:  details,
		CreatedBy:        req.CreatedBy,
		UpdatedBy:        req.CreatedBy,
	}

	if err := s.produkRepo.CreatePenjualan(penjualan); err != nil {
		return nil, fmt.Errorf("failed to create penjualan: %v", err)
	}

	return penjualan, nil
}

// Report Services
func (s *ProdukService) GetStokReport(koperasiID uint64) ([]postgres.Produk, error) {
	return s.produkRepo.GetStokReport(koperasiID)
}

func (s *ProdukService) GetProdukStokRendah(koperasiID uint64) ([]postgres.Produk, error) {
	return s.produkRepo.GetProdukStokRendah(koperasiID)
}

func (s *ProdukService) GetProdukExpiringSoon(koperasiID uint64, days int) ([]postgres.Produk, error) {
	return s.produkRepo.GetProdukExpiringSoon(koperasiID, days)
}