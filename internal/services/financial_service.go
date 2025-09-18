package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type FinancialService struct {
	financialRepo   *postgresRepo.FinancialRepository
	sequenceService *SequenceService
}

func NewFinancialService(
	financialRepo *postgresRepo.FinancialRepository,
	sequenceService *SequenceService,
) *FinancialService {
	return &FinancialService{
		financialRepo:   financialRepo,
		sequenceService: sequenceService,
	}
}

func (s *FinancialService) CreateCOAAkun(req *CreateCOAAkunRequest) (*postgres.COAAkun, error) {
	existing, _ := s.financialRepo.GetCOAAkunByKode(req.KoperasiID, req.KodeAkun)
	if existing != nil {
		return nil, fmt.Errorf("akun with kode %s already exists", req.KodeAkun)
	}

	akun := &postgres.COAAkun{
		TenantID:    req.TenantID,
		KoperasiID:  req.KoperasiID,
		KodeAkun:    req.KodeAkun,
		NamaAkun:    req.NamaAkun,
		KategoriID:  req.KategoriID,
		ParentID:    req.ParentID,
		LevelAkun:   req.LevelAkun,
		SaldoNormal: req.SaldoNormal,
		IsKas:       req.IsKas,
		IsAktif:     true,
	}

	err := s.financialRepo.CreateCOAAkun(akun)
	if err != nil {
		return nil, fmt.Errorf("failed to create COA akun: %v", err)
	}

	return akun, nil
}

func (s *FinancialService) GetCOAAkunList(koperasiID uint64) ([]postgres.COAAkun, error) {
	return s.financialRepo.GetCOAAkunByKoperasi(koperasiID)
}

func (s *FinancialService) GetCOAKategoriList() ([]postgres.COAKategori, error) {
	return s.financialRepo.GetCOAKategoriList()
}

func (s *FinancialService) CreateJurnalUmum(req *CreateJurnalRequest) (*postgres.JurnalUmum, error) {
	if len(req.Details) == 0 {
		return nil, fmt.Errorf("journal details cannot be empty")
	}

	var totalDebit, totalKredit float64
	for _, detail := range req.Details {
		totalDebit += detail.Debit
		totalKredit += detail.Kredit
	}

	if totalDebit != totalKredit {
		return nil, fmt.Errorf("total debit (%.2f) must equal total kredit (%.2f)", totalDebit, totalKredit)
	}

	nomorJurnal, err := s.generateNomorJurnal(req.TenantID, req.KoperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nomor jurnal: %v", err)
	}

	jurnal := &postgres.JurnalUmum{
		TenantID:         req.TenantID,
		KoperasiID:       req.KoperasiID,
		NomorJurnal:      nomorJurnal,
		TanggalTransaksi: req.TanggalTransaksi,
		Referensi:        req.Referensi,
		Keterangan:       req.Keterangan,
		TotalDebit:       totalDebit,
		TotalKredit:      totalKredit,
		Status:           "draft",
		CreatedBy:        req.CreatedBy,
	}

	err = s.financialRepo.CreateJurnalUmum(jurnal)
	if err != nil {
		return nil, fmt.Errorf("failed to create jurnal umum: %v", err)
	}

	var jurnalDetails []postgres.JurnalDetail
	for _, detail := range req.Details {
		jurnalDetail := postgres.JurnalDetail{
			JurnalID:   jurnal.ID,
			AkunID:     detail.AkunID,
			Keterangan: detail.Keterangan,
			Debit:      detail.Debit,
			Kredit:     detail.Kredit,
		}
		jurnalDetails = append(jurnalDetails, jurnalDetail)
	}

	err = s.financialRepo.CreateJurnalDetail(jurnalDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create jurnal details: %v", err)
	}

	return jurnal, nil
}

func (s *FinancialService) GetJurnalUmumByID(id uint64) (*postgres.JurnalUmum, error) {
	return s.financialRepo.GetJurnalUmumByID(id)
}

func (s *FinancialService) GetJurnalUmumList(koperasiID uint64, dari, sampai time.Time, page, limit int) ([]postgres.JurnalUmum, error) {
	offset := (page - 1) * limit
	return s.financialRepo.GetJurnalUmumByKoperasi(koperasiID, dari, sampai, limit, offset)
}

func (s *FinancialService) PostJurnal(id uint64, postedBy uint64) error {
	jurnal, err := s.financialRepo.GetJurnalUmumByID(id)
	if err != nil {
		return fmt.Errorf("jurnal not found: %v", err)
	}

	if jurnal.Status != "draft" {
		return fmt.Errorf("only draft journals can be posted")
	}

	return s.financialRepo.UpdateJurnalStatus(id, "posted", postedBy)
}

func (s *FinancialService) CancelJurnal(id uint64, cancelledBy uint64) error {
	jurnal, err := s.financialRepo.GetJurnalUmumByID(id)
	if err != nil {
		return fmt.Errorf("jurnal not found: %v", err)
	}

	if jurnal.Status == "posted" {
		return fmt.Errorf("posted journals cannot be cancelled, create reversal journal instead")
	}

	return s.financialRepo.UpdateJurnalStatus(id, "cancelled", cancelledBy)
}

func (s *FinancialService) GetNeracaSaldo(koperasiID uint64, tanggal time.Time) ([]postgresRepo.NeracaSaldoItem, error) {
	return s.financialRepo.GetNeracaSaldo(koperasiID, tanggal)
}

func (s *FinancialService) GetLabaRugi(koperasiID uint64, dari, sampai time.Time) (*postgresRepo.LabaRugi, error) {
	return s.financialRepo.GetLabaRugi(koperasiID, dari, sampai)
}

func (s *FinancialService) GetNeraca(koperasiID uint64, tanggal time.Time) (*postgresRepo.Neraca, error) {
	return s.financialRepo.GetNeraca(koperasiID, tanggal)
}

func (s *FinancialService) GetSaldoAkun(akunID uint64, tanggal time.Time) (float64, error) {
	return s.financialRepo.GetSaldoAkun(akunID, tanggal)
}

func (s *FinancialService) generateNomorJurnal(tenantID, koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(tenantID, koperasiID, "jurnal_umum")
	if err != nil {
		return "", err
	}

	now := time.Now()
	return fmt.Sprintf("JU%04d%02d%02d%06d", now.Year(), now.Month(), now.Day(), number), nil
}

type CreateCOAAkunRequest struct {
	TenantID    uint64 `json:"tenant_id" binding:"required"`
	KoperasiID  uint64 `json:"koperasi_id" binding:"required"`
	KodeAkun    string `json:"kode_akun" binding:"required"`
	NamaAkun    string `json:"nama_akun" binding:"required"`
	KategoriID  uint64 `json:"kategori_id" binding:"required"`
	ParentID    uint64 `json:"parent_id"`
	LevelAkun   int    `json:"level_akun"`
	SaldoNormal string `json:"saldo_normal" binding:"required,oneof=debit kredit"`
	IsKas       bool   `json:"is_kas"`
}

type CreateJurnalRequest struct {
	TenantID         uint64                   `json:"tenant_id" binding:"required"`
	KoperasiID       uint64                   `json:"koperasi_id" binding:"required"`
	TanggalTransaksi time.Time                `json:"tanggal_transaksi" binding:"required"`
	Referensi        string                   `json:"referensi"`
	Keterangan       string                   `json:"keterangan" binding:"required"`
	Details          []CreateJurnalDetailRequest `json:"details" binding:"required,min=2"`
	CreatedBy        uint64                   `json:"created_by"`
}

type CreateJurnalDetailRequest struct {
	AkunID     uint64  `json:"akun_id" binding:"required"`
	Keterangan string  `json:"keterangan"`
	Debit      float64 `json:"debit"`
	Kredit     float64 `json:"kredit"`
}