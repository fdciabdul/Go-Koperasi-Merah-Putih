package postgres

import (
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type PPOBRepository struct {
	db *gorm.DB
}

func NewPPOBRepository(db *gorm.DB) *PPOBRepository {
	return &PPOBRepository{db: db}
}

func (r *PPOBRepository) GetKategoriList() ([]postgres.PPOBKategori, error) {
	var kategoris []postgres.PPOBKategori
	err := r.db.Where("is_aktif = ?", true).Order("urutan ASC").Find(&kategoris).Error
	return kategoris, err
}

func (r *PPOBRepository) GetProdukByKategori(kategoriID uint64) ([]postgres.PPOBProduk, error) {
	var produks []postgres.PPOBProduk
	err := r.db.Where("kategori_id = ? AND is_aktif = ?", kategoriID, true).
		Preload("Provider").Preload("Kategori").Find(&produks).Error
	return produks, err
}

func (r *PPOBRepository) GetProdukByID(id uint64) (*postgres.PPOBProduk, error) {
	var produk postgres.PPOBProduk
	err := r.db.Preload("Provider").Preload("Kategori").First(&produk, id).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *PPOBRepository) CreateTransaksi(transaksi *postgres.PPOBTransaksi) error {
	return r.db.Create(transaksi).Error
}

func (r *PPOBRepository) GetTransaksiByID(id uint64) (*postgres.PPOBTransaksi, error) {
	var transaksi postgres.PPOBTransaksi
	err := r.db.Preload("Koperasi").Preload("Anggota").Preload("Produk").
		Preload("Payment").First(&transaksi, id).Error
	if err != nil {
		return nil, err
	}
	return &transaksi, nil
}

func (r *PPOBRepository) GetTransaksiByNomor(nomor string) (*postgres.PPOBTransaksi, error) {
	var transaksi postgres.PPOBTransaksi
	err := r.db.Where("nomor_transaksi = ?", nomor).First(&transaksi).Error
	if err != nil {
		return nil, err
	}
	return &transaksi, nil
}

func (r *PPOBRepository) UpdateTransaksiStatus(id uint64, status string, pesanResponse string) error {
	return r.db.Model(&postgres.PPOBTransaksi{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         status,
		"pesan_response": pesanResponse,
	}).Error
}

func (r *PPOBRepository) UpdatePaymentStatus(id uint64, paymentStatus string) error {
	return r.db.Model(&postgres.PPOBTransaksi{}).Where("id = ?", id).Update("payment_status", paymentStatus).Error
}

func (r *PPOBRepository) GetTransaksiByKoperasi(koperasiID uint64, limit, offset int) ([]postgres.PPOBTransaksi, error) {
	var transaksis []postgres.PPOBTransaksi
	err := r.db.Where("koperasi_id = ?", koperasiID).
		Preload("Produk").Preload("Anggota").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transaksis).Error
	return transaksis, err
}

func (r *PPOBRepository) GetTransaksiForSettlement(koperasiID uint64, dari, sampai time.Time) ([]postgres.PPOBTransaksi, error) {
	var transaksis []postgres.PPOBTransaksi
	err := r.db.Where("koperasi_id = ? AND status = ? AND tanggal_transaksi BETWEEN ? AND ? AND tanggal_settlement IS NULL",
		koperasiID, "success", dari, sampai).
		Preload("Produk").Find(&transaksis).Error
	return transaksis, err
}

func (r *PPOBRepository) CreateSettlement(settlement *postgres.PPOBSettlement) error {
	return r.db.Create(settlement).Error
}

func (r *PPOBRepository) CreateSettlementDetails(details []postgres.PPOBSettlementDetail) error {
	return r.db.Create(&details).Error
}

func (r *PPOBRepository) UpdateSettlementStatus(id uint64, status string, processedBy uint64) error {
	now := time.Now()
	return r.db.Model(&postgres.PPOBSettlement{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"processed_by": processedBy,
		"processed_at": &now,
	}).Error
}

func (r *PPOBRepository) MarkTransaksiSettled(transaksiIDs []uint64) error {
	now := time.Now()
	return r.db.Model(&postgres.PPOBTransaksi{}).Where("id IN ?", transaksiIDs).
		Update("tanggal_settlement", &now).Error
}

func (r *PPOBRepository) GetTransaksiByPaymentID(paymentID uint64) (*postgres.PPOBTransaksi, error) {
	var transaksi postgres.PPOBTransaksi
	err := r.db.Where("payment_id = ?", paymentID).
		Preload("Koperasi").Preload("Anggota").Preload("Produk").
		Preload("Payment").First(&transaksi).Error
	if err != nil {
		return nil, err
	}
	return &transaksi, nil
}

func (r *PPOBRepository) GetPaymentConfig(koperasiID uint64) (*postgres.PPOBPaymentConfig, error) {
	var config postgres.PPOBPaymentConfig
	err := r.db.Where("koperasi_id = ?", koperasiID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}