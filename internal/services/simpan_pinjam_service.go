package services

import (
	"fmt"
	"math"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type SimpanPinjamService struct {
	simpanPinjamRepo *postgresRepo.SimpanPinjamRepository
	sequenceService  *SequenceService
}

func NewSimpanPinjamService(
	simpanPinjamRepo *postgresRepo.SimpanPinjamRepository,
	sequenceService *SequenceService,
) *SimpanPinjamService {
	return &SimpanPinjamService{
		simpanPinjamRepo: simpanPinjamRepo,
		sequenceService:  sequenceService,
	}
}

func (s *SimpanPinjamService) CreateProduk(req *CreateProdukRequest) (*postgres.ProdukSimpanPinjam, error) {
	produk := &postgres.ProdukSimpanPinjam{
		KoperasiID:       req.KoperasiID,
		KodeProduk:       req.KodeProduk,
		NamaProduk:       req.NamaProduk,
		Jenis:            req.Jenis,
		Kategori:         req.Kategori,
		BungaSimpanan:    req.BungaSimpanan,
		MinimalSaldo:     req.MinimalSaldo,
		BungaPinjaman:    req.BungaPinjaman,
		BungaDenda:       req.BungaDenda,
		MaksimalPinjaman: req.MaksimalPinjaman,
		JangkaWaktuMax:   req.JangkaWaktuMax,
		SyaratKetentuan:  req.SyaratKetentuan,
		IsAktif:          true,
	}

	err := s.simpanPinjamRepo.CreateProduk(produk)
	if err != nil {
		return nil, fmt.Errorf("failed to create produk: %v", err)
	}

	return produk, nil
}

func (s *SimpanPinjamService) GetProdukList(koperasiID uint64, jenis string) ([]postgres.ProdukSimpanPinjam, error) {
	return s.simpanPinjamRepo.GetProdukByKoperasi(koperasiID, jenis)
}

func (s *SimpanPinjamService) CreateRekening(req *CreateRekeningRequest) (*postgres.RekeningSimpanPinjam, error) {
	produk, err := s.simpanPinjamRepo.GetProdukByID(req.ProdukID)
	if err != nil {
		return nil, fmt.Errorf("produk not found: %v", err)
	}

	nomorRekening, err := s.generateNomorRekening(req.KoperasiID, produk.Jenis)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nomor rekening: %v", err)
	}

	rekening := &postgres.RekeningSimpanPinjam{
		KoperasiID:    req.KoperasiID,
		AnggotaID:     req.AnggotaID,
		ProdukID:      req.ProdukID,
		NomorRekening: nomorRekening,
		Status:        "aktif",
		TanggalBuka:   time.Now(),
	}

	if produk.Jenis == "pinjaman" {
		rekening.PokokPinjaman = req.PokokPinjaman
		rekening.SisaPokok = req.PokokPinjaman
		rekening.JangkaWaktu = req.JangkaWaktu
		rekening.TanggalMulai = &req.TanggalMulai

		jatuhTempo := req.TanggalMulai.AddDate(0, req.JangkaWaktu, 0)
		rekening.TanggalJatuhTempo = &jatuhTempo

		angsuranPokok, angsuranBunga := s.calculateAngsuran(
			req.PokokPinjaman, produk.BungaPinjaman, req.JangkaWaktu)
		rekening.AngsuranPokok = angsuranPokok
		rekening.AngsuranBunga = angsuranBunga
	}

	err = s.simpanPinjamRepo.CreateRekening(rekening)
	if err != nil {
		return nil, fmt.Errorf("failed to create rekening: %v", err)
	}

	return rekening, nil
}

func (s *SimpanPinjamService) GetRekeningByAnggota(anggotaID uint64) ([]postgres.RekeningSimpanPinjam, error) {
	return s.simpanPinjamRepo.GetRekeningByAnggota(anggotaID)
}

func (s *SimpanPinjamService) CreateTransaksi(req *CreateTransaksiRequest) (*postgres.TransaksiSimpanPinjam, error) {
	rekening, err := s.simpanPinjamRepo.GetRekeningByID(req.RekeningID)
	if err != nil {
		return nil, fmt.Errorf("rekening not found: %v", err)
	}

	if rekening.Status != "aktif" {
		return nil, fmt.Errorf("rekening not active")
	}

	nomorTransaksi, err := s.generateNomorTransaksi(req.KoperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nomor transaksi: %v", err)
	}

	saldoSebelum := rekening.SaldoSimpanan
	if rekening.Produk.Jenis == "pinjaman" {
		saldoSebelum = rekening.SisaPokok
	}

	saldoSesudah := saldoSebelum
	switch req.JenisTransaksi {
	case "setoran":
		saldoSesudah = saldoSebelum + req.Jumlah
		rekening.SaldoSimpanan = saldoSesudah
	case "penarikan":
		if saldoSebelum < req.Jumlah {
			return nil, fmt.Errorf("insufficient balance")
		}
		saldoSesudah = saldoSebelum - req.Jumlah
		rekening.SaldoSimpanan = saldoSesudah
	case "pencairan":
		rekening.PokokPinjaman = req.Jumlah
		rekening.SisaPokok = req.Jumlah
		saldoSesudah = req.Jumlah
	case "angsuran":
		if saldoSebelum < req.Jumlah {
			return nil, fmt.Errorf("angsuran amount exceeds remaining principal")
		}
		saldoSesudah = saldoSebelum - req.Jumlah
		rekening.SisaPokok = saldoSesudah
		if saldoSesudah == 0 {
			rekening.Status = "lunas"
		}
	}

	transaksi := &postgres.TransaksiSimpanPinjam{
		KoperasiID:       req.KoperasiID,
		RekeningID:       req.RekeningID,
		NomorTransaksi:   nomorTransaksi,
		TanggalTransaksi: time.Now(),
		JenisTransaksi:   req.JenisTransaksi,
		Jumlah:           req.Jumlah,
		SaldoSebelum:     saldoSebelum,
		SaldoSesudah:     saldoSesudah,
		Keterangan:       req.Keterangan,
		Referensi:        req.Referensi,
		CreatedBy:        req.CreatedBy,
	}

	err = s.simpanPinjamRepo.CreateTransaksi(transaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaksi: %v", err)
	}

	err = s.simpanPinjamRepo.UpdateRekening(rekening)
	if err != nil {
		return nil, fmt.Errorf("failed to update rekening: %v", err)
	}

	return transaksi, nil
}

func (s *SimpanPinjamService) GetTransaksiByRekening(rekeningID uint64, page, limit int) ([]postgres.TransaksiSimpanPinjam, error) {
	offset := (page - 1) * limit
	return s.simpanPinjamRepo.GetTransaksiByRekening(rekeningID, limit, offset)
}

func (s *SimpanPinjamService) GetStatistik(koperasiID uint64) (*postgresRepo.SimpanPinjamStatistik, error) {
	return s.simpanPinjamRepo.GetStatistikSimpanPinjam(koperasiID)
}

func (s *SimpanPinjamService) GetPinjamanJatuhTempo(days int) ([]postgres.RekeningSimpanPinjam, error) {
	return s.simpanPinjamRepo.GetRekeningPinjamanJatuhTempo(days)
}

func (s *SimpanPinjamService) generateNomorRekening(koperasiID uint64, jenis string) (string, error) {
	var prefix string
	switch jenis {
	case "simpanan":
		prefix = "SIM"
	case "pinjaman":
		prefix = "PIN"
	default:
		prefix = "REK"
	}

	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "rekening_"+jenis)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d%08d", prefix, koperasiID, number), nil
}

func (s *SimpanPinjamService) generateNomorTransaksi(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "transaksi_simpan_pinjam")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("TRX%04d%010d", koperasiID, number), nil
}

func (s *SimpanPinjamService) calculateAngsuran(pokok, bunga float64, jangkaWaktu int) (float64, float64) {
	bungaBulanan := bunga / 12 / 100

	if bungaBulanan == 0 {
		return pokok / float64(jangkaWaktu), 0
	}

	angsuranPokok := pokok * (bungaBulanan * math.Pow(1+bungaBulanan, float64(jangkaWaktu))) /
		(math.Pow(1+bungaBulanan, float64(jangkaWaktu)) - 1)

	angsuranBunga := pokok * bungaBulanan

	return angsuranPokok, angsuranBunga
}

type CreateProdukRequest struct {
	KoperasiID       uint64  `json:"koperasi_id" binding:"required"`
	KodeProduk       string  `json:"kode_produk" binding:"required"`
	NamaProduk       string  `json:"nama_produk" binding:"required"`
	Jenis            string  `json:"jenis" binding:"required,oneof=simpanan pinjaman"`
	Kategori         string  `json:"kategori"`
	BungaSimpanan    float64 `json:"bunga_simpanan"`
	MinimalSaldo     float64 `json:"minimal_saldo"`
	BungaPinjaman    float64 `json:"bunga_pinjaman"`
	BungaDenda       float64 `json:"bunga_denda"`
	MaksimalPinjaman float64 `json:"maksimal_pinjaman"`
	JangkaWaktuMax   int     `json:"jangka_waktu_max"`
	SyaratKetentuan  string  `json:"syarat_ketentuan"`
}

type CreateRekeningRequest struct {
	KoperasiID     uint64    `json:"koperasi_id" binding:"required"`
	AnggotaID      uint64    `json:"anggota_id" binding:"required"`
	ProdukID       uint64    `json:"produk_id" binding:"required"`
	PokokPinjaman  float64   `json:"pokok_pinjaman"`
	JangkaWaktu    int       `json:"jangka_waktu"`
	TanggalMulai   time.Time `json:"tanggal_mulai"`
}

type CreateTransaksiRequest struct {
	KoperasiID      uint64  `json:"koperasi_id" binding:"required"`
	RekeningID      uint64  `json:"rekening_id" binding:"required"`
	JenisTransaksi  string  `json:"jenis_transaksi" binding:"required,oneof=setoran penarikan pencairan angsuran bunga denda"`
	Jumlah          float64 `json:"jumlah" binding:"required,gt=0"`
	Keterangan      string  `json:"keterangan"`
	Referensi       string  `json:"referensi"`
	CreatedBy       uint64  `json:"created_by"`
}