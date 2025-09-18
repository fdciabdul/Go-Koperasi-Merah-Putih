package postgres

import (
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type WilayahRepository struct {
	db *gorm.DB
}

func NewWilayahRepository(db *gorm.DB) *WilayahRepository {
	return &WilayahRepository{db: db}
}

func (r *WilayahRepository) GetProvinsiList() ([]postgres.Provinsi, error) {
	var provinsis []postgres.Provinsi
	err := r.db.Order("nama ASC").Find(&provinsis).Error
	return provinsis, err
}

func (r *WilayahRepository) GetKabupatenByProvinsi(provinsiID uint64) ([]postgres.Kabupaten, error) {
	var kabupatens []postgres.Kabupaten
	err := r.db.Where("provinsi_id = ?", provinsiID).
		Order("nama ASC").Find(&kabupatens).Error
	return kabupatens, err
}

func (r *WilayahRepository) GetKecamatanByKabupaten(kabupatenID uint64) ([]postgres.Kecamatan, error) {
	var kecamatans []postgres.Kecamatan
	err := r.db.Where("kabupaten_id = ?", kabupatenID).
		Order("nama ASC").Find(&kecamatans).Error
	return kecamatans, err
}

func (r *WilayahRepository) GetKelurahanByKecamatan(kecamatanID uint64) ([]postgres.Kelurahan, error) {
	var kelurahans []postgres.Kelurahan
	err := r.db.Where("kecamatan_id = ?", kecamatanID).
		Order("nama ASC").Find(&kelurahans).Error
	return kelurahans, err
}