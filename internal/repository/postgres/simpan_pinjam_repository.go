package postgres

import (
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type SimpanPinjamRepository struct {
	db *gorm.DB
}

func NewSimpanPinjamRepository(db *gorm.DB) *SimpanPinjamRepository {
	return &SimpanPinjamRepository{db: db}
}

func (r *SimpanPinjamRepository) CreateProduk(produk *postgres.ProdukSimpanPinjam) error {
	return r.db.Create(produk).Error
}

func (r *SimpanPinjamRepository) GetProdukByID(id uint64) (*postgres.ProdukSimpanPinjam, error) {
	var produk postgres.ProdukSimpanPinjam
	err := r.db.Preload("Koperasi").First(&produk, id).Error
	if err != nil {
		return nil, err
	}
	return &produk, nil
}

func (r *SimpanPinjamRepository) GetProdukByKoperasi(koperasiID uint64, jenis string) ([]postgres.ProdukSimpanPinjam, error) {
	var produks []postgres.ProdukSimpanPinjam
	query := r.db.Where("koperasi_id = ? AND is_aktif = ?", koperasiID, true)

	if jenis != "" {
		query = query.Where("jenis = ?", jenis)
	}

	err := query.Find(&produks).Error
	return produks, err
}

func (r *SimpanPinjamRepository) UpdateProduk(produk *postgres.ProdukSimpanPinjam) error {
	return r.db.Save(produk).Error
}

func (r *SimpanPinjamRepository) DeleteProduk(id uint64) error {
	return r.db.Model(&postgres.ProdukSimpanPinjam{}).Where("id = ?", id).Update("is_aktif", false).Error
}

func (r *SimpanPinjamRepository) CreateRekening(rekening *postgres.RekeningSimpanPinjam) error {
	return r.db.Create(rekening).Error
}

func (r *SimpanPinjamRepository) GetRekeningByID(id uint64) (*postgres.RekeningSimpanPinjam, error) {
	var rekening postgres.RekeningSimpanPinjam
	err := r.db.Preload("Koperasi").Preload("Anggota").Preload("Produk").
		First(&rekening, id).Error
	if err != nil {
		return nil, err
	}
	return &rekening, nil
}

func (r *SimpanPinjamRepository) GetRekeningByNomor(nomorRekening string) (*postgres.RekeningSimpanPinjam, error) {
	var rekening postgres.RekeningSimpanPinjam
	err := r.db.Where("nomor_rekening = ?", nomorRekening).First(&rekening).Error
	if err != nil {
		return nil, err
	}
	return &rekening, nil
}

func (r *SimpanPinjamRepository) GetRekeningByAnggota(anggotaID uint64) ([]postgres.RekeningSimpanPinjam, error) {
	var rekenings []postgres.RekeningSimpanPinjam
	err := r.db.Where("anggota_id = ? AND status = ?", anggotaID, "aktif").
		Preload("Produk").Find(&rekenings).Error
	return rekenings, err
}

func (r *SimpanPinjamRepository) UpdateRekening(rekening *postgres.RekeningSimpanPinjam) error {
	return r.db.Save(rekening).Error
}

func (r *SimpanPinjamRepository) CreateTransaksi(transaksi *postgres.TransaksiSimpanPinjam) error {
	return r.db.Create(transaksi).Error
}

func (r *SimpanPinjamRepository) GetTransaksiByID(id uint64) (*postgres.TransaksiSimpanPinjam, error) {
	var transaksi postgres.TransaksiSimpanPinjam
	err := r.db.Preload("Koperasi").Preload("Rekening").Preload("Jurnal").
		First(&transaksi, id).Error
	if err != nil {
		return nil, err
	}
	return &transaksi, nil
}

func (r *SimpanPinjamRepository) GetTransaksiByRekening(rekeningID uint64, limit, offset int) ([]postgres.TransaksiSimpanPinjam, error) {
	var transaksis []postgres.TransaksiSimpanPinjam
	err := r.db.Where("rekening_id = ?", rekeningID).
		Order("tanggal_transaksi DESC").
		Limit(limit).Offset(offset).
		Find(&transaksis).Error
	return transaksis, err
}

func (r *SimpanPinjamRepository) GetTransaksiByKoperasi(koperasiID uint64, dari, sampai time.Time) ([]postgres.TransaksiSimpanPinjam, error) {
	var transaksis []postgres.TransaksiSimpanPinjam
	err := r.db.Where("koperasi_id = ? AND tanggal_transaksi BETWEEN ? AND ?",
		koperasiID, dari, sampai).
		Preload("Rekening").Preload("Rekening.Anggota").Preload("Rekening.Produk").
		Order("tanggal_transaksi DESC").
		Find(&transaksis).Error
	return transaksis, err
}

func (r *SimpanPinjamRepository) GetRekeningPinjamanJatuhTempo(days int) ([]postgres.RekeningSimpanPinjam, error) {
	var rekenings []postgres.RekeningSimpanPinjam
	targetDate := time.Now().AddDate(0, 0, days)

	err := r.db.Where("status = ? AND tanggal_jatuh_tempo <= ?", "aktif", targetDate).
		Preload("Anggota").Preload("Produk").Find(&rekenings).Error
	return rekenings, err
}

func (r *SimpanPinjamRepository) GetStatistikSimpanPinjam(koperasiID uint64) (*SimpanPinjamStatistik, error) {
	var statistik SimpanPinjamStatistik

	err := r.db.Model(&postgres.RekeningSimpanPinjam{}).
		Select(`
			COUNT(CASE WHEN produk.jenis = 'simpanan' THEN 1 END) as total_rekening_simpanan,
			COUNT(CASE WHEN produk.jenis = 'pinjaman' THEN 1 END) as total_rekening_pinjaman,
			SUM(CASE WHEN produk.jenis = 'simpanan' THEN saldo_simpanan ELSE 0 END) as total_saldo_simpanan,
			SUM(CASE WHEN produk.jenis = 'pinjaman' THEN sisa_pokok ELSE 0 END) as total_sisa_pinjaman
		`).
		Joins("JOIN produk_simpan_pinjam produk ON rekening_simpan_pinjam.produk_id = produk.id").
		Where("rekening_simpan_pinjam.koperasi_id = ? AND rekening_simpan_pinjam.status = ?", koperasiID, "aktif").
		Scan(&statistik).Error

	return &statistik, err
}

type SimpanPinjamStatistik struct {
	TotalRekeningSimpanan uint64  `json:"total_rekening_simpanan"`
	TotalRekeningPinjaman uint64  `json:"total_rekening_pinjaman"`
	TotalSaldoSimpanan    float64 `json:"total_saldo_simpanan"`
	TotalSisaPinjaman     float64 `json:"total_sisa_pinjaman"`
}