package services

import (
	"fmt"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type MasterDataService struct {
	masterDataRepo *postgresRepo.MasterDataRepository
}

func NewMasterDataService(masterDataRepo *postgresRepo.MasterDataRepository) *MasterDataService {
	return &MasterDataService{masterDataRepo: masterDataRepo}
}

func (s *MasterDataService) CreateKBLI(req *CreateKBLIRequest) (*postgres.KBLI, error) {
	existing, _ := s.masterDataRepo.GetKBLIByKode(req.Kode)
	if existing != nil {
		return nil, fmt.Errorf("KBLI with kode %s already exists", req.Kode)
	}

	kbli := &postgres.KBLI{
		Kode:        req.Kode,
		Nama:        req.Nama,
		Kategori:    req.Kategori,
		Deskripsi:   req.Deskripsi,
		IsAktif:     true,
	}

	err := s.masterDataRepo.CreateKBLI(kbli)
	if err != nil {
		return nil, fmt.Errorf("failed to create KBLI: %v", err)
	}

	return kbli, nil
}

func (s *MasterDataService) GetKBLIList(search string, page, limit int) ([]postgres.KBLI, error) {
	offset := (page - 1) * limit
	return s.masterDataRepo.GetKBLIList(search, limit, offset)
}

func (s *MasterDataService) GetKBLIByID(id uint64) (*postgres.KBLI, error) {
	return s.masterDataRepo.GetKBLIByID(id)
}

func (s *MasterDataService) UpdateKBLI(id uint64, req *UpdateKBLIRequest) (*postgres.KBLI, error) {
	kbli, err := s.masterDataRepo.GetKBLIByID(id)
	if err != nil {
		return nil, fmt.Errorf("KBLI not found: %v", err)
	}

	kbli.Nama = req.Nama
	kbli.Kategori = req.Kategori
	kbli.Deskripsi = req.Deskripsi
	kbli.IsAktif = req.IsAktif

	err = s.masterDataRepo.UpdateKBLI(kbli)
	if err != nil {
		return nil, fmt.Errorf("failed to update KBLI: %v", err)
	}

	return kbli, nil
}

func (s *MasterDataService) CreateJenisKoperasi(req *CreateJenisKoperasiRequest) (*postgres.JenisKoperasi, error) {
	existing, _ := s.masterDataRepo.GetJenisKoperasiByKode(req.Kode)
	if existing != nil {
		return nil, fmt.Errorf("jenis koperasi with kode %s already exists", req.Kode)
	}

	jenis := &postgres.JenisKoperasi{
		Kode:      req.Kode,
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		IsAktif:   true,
	}

	err := s.masterDataRepo.CreateJenisKoperasi(jenis)
	if err != nil {
		return nil, fmt.Errorf("failed to create jenis koperasi: %v", err)
	}

	return jenis, nil
}

func (s *MasterDataService) GetJenisKoperasiList() ([]postgres.JenisKoperasi, error) {
	return s.masterDataRepo.GetJenisKoperasiList()
}

func (s *MasterDataService) GetJenisKoperasiByID(id uint64) (*postgres.JenisKoperasi, error) {
	return s.masterDataRepo.GetJenisKoperasiByID(id)
}

func (s *MasterDataService) UpdateJenisKoperasi(id uint64, req *UpdateJenisKoperasiRequest) (*postgres.JenisKoperasi, error) {
	jenis, err := s.masterDataRepo.GetJenisKoperasiByID(id)
	if err != nil {
		return nil, fmt.Errorf("jenis koperasi not found: %v", err)
	}

	jenis.Nama = req.Nama
	jenis.Deskripsi = req.Deskripsi
	jenis.IsAktif = req.IsAktif

	err = s.masterDataRepo.UpdateJenisKoperasi(jenis)
	if err != nil {
		return nil, fmt.Errorf("failed to update jenis koperasi: %v", err)
	}

	return jenis, nil
}

func (s *MasterDataService) CreateBentukKoperasi(req *CreateBentukKoperasiRequest) (*postgres.BentukKoperasi, error) {
	existing, _ := s.masterDataRepo.GetBentukKoperasiByKode(req.Kode)
	if existing != nil {
		return nil, fmt.Errorf("bentuk koperasi with kode %s already exists", req.Kode)
	}

	bentuk := &postgres.BentukKoperasi{
		Kode:      req.Kode,
		Nama:      req.Nama,
		Deskripsi: req.Deskripsi,
		IsAktif:   true,
	}

	err := s.masterDataRepo.CreateBentukKoperasi(bentuk)
	if err != nil {
		return nil, fmt.Errorf("failed to create bentuk koperasi: %v", err)
	}

	return bentuk, nil
}

func (s *MasterDataService) GetBentukKoperasiList() ([]postgres.BentukKoperasi, error) {
	return s.masterDataRepo.GetBentukKoperasiList()
}

func (s *MasterDataService) GetBentukKoperasiByID(id uint64) (*postgres.BentukKoperasi, error) {
	return s.masterDataRepo.GetBentukKoperasiByID(id)
}

func (s *MasterDataService) UpdateBentukKoperasi(id uint64, req *UpdateBentukKoperasiRequest) (*postgres.BentukKoperasi, error) {
	bentuk, err := s.masterDataRepo.GetBentukKoperasiByID(id)
	if err != nil {
		return nil, fmt.Errorf("bentuk koperasi not found: %v", err)
	}

	bentuk.Nama = req.Nama
	bentuk.Deskripsi = req.Deskripsi
	bentuk.IsAktif = req.IsAktif

	err = s.masterDataRepo.UpdateBentukKoperasi(bentuk)
	if err != nil {
		return nil, fmt.Errorf("failed to update bentuk koperasi: %v", err)
	}

	return bentuk, nil
}

type CreateKBLIRequest struct {
	Kode      string `json:"kode" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

type UpdateKBLIRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	IsAktif   bool   `json:"is_aktif"`
}

type CreateJenisKoperasiRequest struct {
	Kode      string `json:"kode" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

type UpdateJenisKoperasiRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	IsAktif   bool   `json:"is_aktif"`
}

type CreateBentukKoperasiRequest struct {
	Kode      string `json:"kode" binding:"required"`
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
}

type UpdateBentukKoperasiRequest struct {
	Nama      string `json:"nama" binding:"required"`
	Deskripsi string `json:"deskripsi"`
	IsAktif   bool   `json:"is_aktif"`
}