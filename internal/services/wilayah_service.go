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

func (s *WilayahService) GetProvinsiList() ([]postgres.WilayahProvinsi, error) {
	return s.wilayahRepo.GetProvinsiList()
}

func (s *WilayahService) GetKabupatenByProvinsi(provinsiID uint64) ([]postgres.WilayahKabupaten, error) {
	return s.wilayahRepo.GetKabupatenByProvinsiID(provinsiID)
}

func (s *WilayahService) GetKecamatanByKabupaten(kabupatenID uint64) ([]postgres.WilayahKecamatan, error) {
	return s.wilayahRepo.GetKecamatanByKabupatenID(kabupatenID)
}

func (s *WilayahService) GetKelurahanByKecamatan(kecamatanID uint64) ([]postgres.WilayahKelurahan, error) {
	return s.wilayahRepo.GetKelurahanByKecamatanID(kecamatanID)
}