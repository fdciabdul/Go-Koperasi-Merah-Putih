package postgres

import (
	"time"

	"koperasi-merah-putih/internal/models/postgres"

	"gorm.io/gorm"
)

type ProdukRepository struct {
	db *gorm.DB
}

func NewProdukRepository(db *gorm.DB) *ProdukRepository {
	return &ProdukRepository{db: db}
}

// Kategori Produk
func (r *ProdukRepository) CreateKategoriProduk(kategori *postgres.KategoriProduk) error {
	return r.db.Create(kategori).Error
}

func (r *ProdukRepository) GetKategoriProdukByID(id uint64) (*postgres.KategoriProduk, error) {
	var kategori postgres.KategoriProduk
	err := r.db.First(&kategori, id).Error
	if err != nil {
		return nil, err
	}
	return &kategori, nil
}

func (r *ProdukRepository) GetAllKategoriProduk() ([]postgres.KategoriProduk, error) {
	var kategori []postgres.KategoriProduk
	err := r.db.Where("is_active = ?", true).Find(&kategori).Error
	return kategori, err
}

func (r *ProdukRepository) UpdateKategoriProduk(kategori *postgres.KategoriProduk) error {
	return r.db.Save(kategori).Error
}

func (r *ProdukRepository) DeleteKategoriProduk(id uint64) error {
	return r.db.Delete(&postgres.KategoriProduk{}, id).Error
}

// Satuan Produk
func (r *ProdukRepository) CreateSatuanProduk(satuan *postgres.SatuanProduk) error {
	return r.db.Create(satuan).Error
}

func (r *ProdukRepository) GetSatuanProdukByID(id uint64) (*postgres.SatuanProduk, error) {
	var satuan postgres.SatuanProduk
	err := r.db.First(&satuan, id).Error
	if err != nil {
		return nil, err
	}
	return &satuan, nil
}

func (r *ProdukRepository) GetAllSatuanProduk() ([]postgres.SatuanProduk, error) {
	var satuan []postgres.SatuanProduk
	err := r.db.Where("is_active = ?", true).Find(&satuan).Error
	return satuan, err
}

func (r *ProdukRepository) UpdateSatuanProduk(satuan *postgres.SatuanProduk) error {
	return r.db.Save(satuan).Error
}

func (r *ProdukRepository) DeleteSatuanProduk(id uint64) error {
	return r.db.Delete(&postgres.SatuanProduk{}, id).Error
}

// Supplier
func (r *ProdukRepository) CreateSupplier(supplier *postgres.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *ProdukRepository) GetSupplierByID(id uint64) (*postgres.Supplier, error) {
	var supplier postgres.Supplier
	err := r.db.Preload("Provinsi").Preload("Kabupaten").First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *ProdukRepository) GetSuppliersByKoperasi(koperasiID uint64, limit, offset int) ([]postgres.Supplier, error) {
	var suppliers []postgres.Supplier
	err := r.db.Where("koperasi_id = ? AND is_active = ?", koperasiID, true).
		Limit(limit).Offset(offset).Find(&suppliers).Error
	return suppliers, err
}

func (r *ProdukRepository) UpdateSupplier(supplier *postgres.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *ProdukRepository) DeleteSupplier(id uint64) error {
	return r.db.Delete(&postgres.Supplier{}, id).Error
}

// Produk
func (r *ProdukRepository) CreateProduk(produk *postgres.Produk) error {
	return r.db.Create(produk).Error
}

func (r *ProdukRepository) GetProdukByID(id uint64) (*postgres.Produk, error) {
	var produk postgres.Produk
	err := r.db.Preload("KategoriProduk").Preload("SatuanProduk").First(&produk, id).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *ProdukRepository) GetProdukByBarcode(barcode string) (*postgres.Produk, error) {
	var produk postgres.Produk
	err := r.db.Where("barcode = ?", barcode).
		Preload("KategoriProduk").Preload("SatuanProduk").First(&produk).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *ProdukRepository) GetProduksByKoperasi(koperasiID uint64, filters ProdukFilters, limit, offset int) ([]postgres.Produk, error) {
	var produk []postgres.Produk

	query := r.db.Where("koperasi_id = ? AND is_active = ?", koperasiID, true)

	if filters.KategoriID != 0 {
		query = query.Where("kategori_produk_id = ?", filters.KategoriID)
	}

	if filters.Nama != "" {
		query = query.Where("nama_produk ILIKE ?", "%"+filters.Nama+"%")
	}

	if filters.Brand != "" {
		query = query.Where("brand ILIKE ?", "%"+filters.Brand+"%")
	}

	if filters.StokRendah {
		query = query.Where("stok_current <= stok_minimal")
	}

	if filters.ReadyStock {
		query = query.Where("is_ready_stock = ?", true)
	}

	err := query.Preload("KategoriProduk").Preload("SatuanProduk").
		Limit(limit).Offset(offset).Find(&produk).Error
	return produk, err
}

func (r *ProdukRepository) UpdateProduk(produk *postgres.Produk) error {
	return r.db.Save(produk).Error
}

func (r *ProdukRepository) UpdateStokProduk(produkID uint64, newStok int) error {
	return r.db.Model(&postgres.Produk{}).Where("id = ?", produkID).Update("stok_current", newStok).Error
}

func (r *ProdukRepository) DeleteProduk(id uint64) error {
	return r.db.Delete(&postgres.Produk{}, id).Error
}

// Purchase Order
func (r *ProdukRepository) CreatePurchaseOrder(po *postgres.PurchaseOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(po).Error; err != nil {
			return err
		}

		for i := range po.PurchaseOrderDetail {
			po.PurchaseOrderDetail[i].PurchaseOrderID = po.ID
		}

		if err := tx.CreateInBatches(po.PurchaseOrderDetail, 50).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *ProdukRepository) GetPurchaseOrderByID(id uint64) (*postgres.PurchaseOrder, error) {
	var po postgres.PurchaseOrder
	err := r.db.Preload("Supplier").Preload("PurchaseOrderDetail.Produk").First(&po, id).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (r *ProdukRepository) GetPurchaseOrdersByKoperasi(koperasiID uint64, limit, offset int) ([]postgres.PurchaseOrder, error) {
	var pos []postgres.PurchaseOrder
	err := r.db.Where("koperasi_id = ?", koperasiID).
		Preload("Supplier").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&pos).Error
	return pos, err
}

func (r *ProdukRepository) UpdatePurchaseOrder(po *postgres.PurchaseOrder) error {
	return r.db.Save(po).Error
}

// Pembelian
func (r *ProdukRepository) CreatePembelian(pembelian *postgres.PembelianHeader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pembelian).Error; err != nil {
			return err
		}

		for i := range pembelian.PembelianDetail {
			pembelian.PembelianDetail[i].PembelianHeaderID = pembelian.ID
		}

		if err := tx.CreateInBatches(pembelian.PembelianDetail, 50).Error; err != nil {
			return err
		}

		for _, detail := range pembelian.PembelianDetail {
			var produk postgres.Produk
			if err := tx.First(&produk, detail.ProdukID).Error; err == nil {
				newStok := produk.StokCurrent + detail.Qty
				if err := tx.Model(&produk).Update("stok_current", newStok).Error; err != nil {
					return err
				}

				stokMovement := postgres.StokMovement{
					KoperasiID:      pembelian.KoperasiID,
					ProdukID:        detail.ProdukID,
					TipeMovement:    "in",
					ReferensiTipe:   "pembelian",
					ReferensiID:     pembelian.ID,
					TanggalMovement: pembelian.TanggalFaktur,
					QtyBefore:       produk.StokCurrent,
					QtyMovement:     detail.Qty,
					QtyAfter:        newStok,
					HargaSatuan:     detail.HargaSatuan,
					TotalNilai:      detail.Subtotal,
					CreatedBy:       pembelian.CreatedBy,
				}
				if err := tx.Create(&stokMovement).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *ProdukRepository) GetPembelianByID(id uint64) (*postgres.PembelianHeader, error) {
	var pembelian postgres.PembelianHeader
	err := r.db.Preload("Supplier").Preload("PembelianDetail.Produk").First(&pembelian, id).Error
	if err != nil {
		return nil, err
	}
	return &pembelian, nil
}

// Penjualan
func (r *ProdukRepository) CreatePenjualan(penjualan *postgres.PenjualanHeader) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(penjualan).Error; err != nil {
			return err
		}

		for i := range penjualan.PenjualanDetail {
			penjualan.PenjualanDetail[i].PenjualanHeaderID = penjualan.ID
		}

		if err := tx.CreateInBatches(penjualan.PenjualanDetail, 50).Error; err != nil {
			return err
		}

		for _, detail := range penjualan.PenjualanDetail {
			var produk postgres.Produk
			if err := tx.First(&produk, detail.ProdukID).Error; err == nil {
				if produk.StokCurrent < detail.Qty {
					return gorm.ErrInvalidValue
				}

				newStok := produk.StokCurrent - detail.Qty
				if err := tx.Model(&produk).Update("stok_current", newStok).Error; err != nil {
					return err
				}

				stokMovement := postgres.StokMovement{
					KoperasiID:      penjualan.KoperasiID,
					ProdukID:        detail.ProdukID,
					TipeMovement:    "out",
					ReferensiTipe:   "penjualan",
					ReferensiID:     penjualan.ID,
					TanggalMovement: penjualan.TanggalTransaksi,
					QtyBefore:       produk.StokCurrent,
					QtyMovement:     detail.Qty,
					QtyAfter:        newStok,
					HargaSatuan:     detail.HargaSatuan,
					TotalNilai:      detail.Subtotal,
					CreatedBy:       penjualan.CreatedBy,
				}
				if err := tx.Create(&stokMovement).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *ProdukRepository) GetPenjualanByID(id uint64) (*postgres.PenjualanHeader, error) {
	var penjualan postgres.PenjualanHeader
	err := r.db.Preload("Anggota").Preload("PenjualanDetail.Produk").First(&penjualan, id).Error
	if err != nil {
		return nil, err
	}
	return &penjualan, nil
}

func (r *ProdukRepository) GetPenjualansByKoperasi(koperasiID uint64, startDate, endDate time.Time, limit, offset int) ([]postgres.PenjualanHeader, error) {
	var penjualan []postgres.PenjualanHeader
	err := r.db.Where("koperasi_id = ? AND tanggal_transaksi BETWEEN ? AND ?", koperasiID, startDate, endDate).
		Order("tanggal_transaksi DESC").
		Limit(limit).Offset(offset).Find(&penjualan).Error
	return penjualan, err
}

// Stok Movement
func (r *ProdukRepository) CreateStokMovement(movement *postgres.StokMovement) error {
	return r.db.Create(movement).Error
}

func (r *ProdukRepository) GetStokMovementByProduk(produkID uint64, limit, offset int) ([]postgres.StokMovement, error) {
	var movements []postgres.StokMovement
	err := r.db.Where("produk_id = ?", produkID).
		Order("tanggal_movement DESC").
		Limit(limit).Offset(offset).Find(&movements).Error
	return movements, err
}

// Product Filters
type ProdukFilters struct {
	KategoriID  uint64
	Nama        string
	Brand       string
	StokRendah  bool
	ReadyStock  bool
}

// Reports
func (r *ProdukRepository) GetStokReport(koperasiID uint64) ([]postgres.Produk, error) {
	var produk []postgres.Produk
	err := r.db.Where("koperasi_id = ? AND is_active = ?", koperasiID, true).
		Preload("KategoriProduk").Preload("SatuanProduk").
		Order("nama_produk").Find(&produk).Error
	return produk, err
}

func (r *ProdukRepository) GetProdukStokRendah(koperasiID uint64) ([]postgres.Produk, error) {
	var produk []postgres.Produk
	err := r.db.Where("koperasi_id = ? AND stok_current <= stok_minimal AND is_active = ?", koperasiID, true).
		Preload("KategoriProduk").Preload("SatuanProduk").
		Order("stok_current").Find(&produk).Error
	return produk, err
}

func (r *ProdukRepository) GetProdukExpiringSoon(koperasiID uint64, days int) ([]postgres.Produk, error) {
	var produk []postgres.Produk
	expiredDate := time.Now().AddDate(0, 0, days)

	err := r.db.Where("koperasi_id = ? AND tanggal_expired <= ? AND tanggal_expired IS NOT NULL AND is_active = ?",
		koperasiID, expiredDate, true).
		Preload("KategoriProduk").Preload("SatuanProduk").
		Order("tanggal_expired").Find(&produk).Error
	return produk, err
}