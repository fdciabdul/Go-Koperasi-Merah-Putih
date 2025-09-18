package services

import (
	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type WilayahService struct {
	wilayahRepo *postgresRepo.WilayahRepository
}

func NewWilayahService(wilayahRepo *postgresRepo.WilayahRepository) *WilayahService {
	return &WilayahService{wilayahRepo: wilayahRepo}
}

func (s *WilayahService) GetProvinsiList() ([]postgres.Provinsi, error) {
	return s.wilayahRepo.GetProvinsiList()
}

func (s *WilayahService) GetKabupatenByProvinsi(provinsiID uint64) ([]postgres.Kabupaten, error) {
	return s.wilayahRepo.GetKabupatenByProvinsi(provinsiID)
}

func (s *WilayahService) GetKecamatanByKabupaten(kabupatenID uint64) ([]postgres.Kecamatan, error) {
	return s.wilayahRepo.GetKecamatanByKabupaten(kabupatenID)
}

func (s *WilayahService) GetKelurahanByKecamatan(kecamatanID uint64) ([]postgres.Kelurahan, error) {
	return s.wilayahRepo.GetKelurahanByKecamatan(kecamatanID)
}