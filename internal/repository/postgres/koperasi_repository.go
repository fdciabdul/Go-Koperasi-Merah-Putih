package postgres

import (
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type KoperasiRepository struct {
	db *gorm.DB
}

func NewKoperasiRepository(db *gorm.DB) *KoperasiRepository {
	return &KoperasiRepository{db: db}
}

func (r *KoperasiRepository) Create(koperasi *postgres.Koperasi) error {
	return r.db.Create(koperasi).Error
}

func (r *KoperasiRepository) GetByID(id uint64) (*postgres.Koperasi, error) {
	var koperasi postgres.Koperasi
	err := r.db.Preload("Tenant").Preload("JenisKoperasi").Preload("BentukKoperasi").
		Preload("StatusKoperasi").Preload("Provinsi").Preload("Kabupaten").
		Preload("Kecamatan").Preload("Kelurahan").First(&koperasi, id).Error
	if err != nil {
		return nil, err
	}
	return &koperasi, nil
}

func (r *KoperasiRepository) GetByTenantID(tenantID uint64) ([]postgres.Koperasi, error) {
	var koperasis []postgres.Koperasi
	err := r.db.Where("tenant_id = ?", tenantID).
		Preload("JenisKoperasi").Preload("BentukKoperasi").Preload("StatusKoperasi").
		Find(&koperasis).Error
	return koperasis, err
}

func (r *KoperasiRepository) Update(koperasi *postgres.Koperasi) error {
	return r.db.Save(koperasi).Error
}

func (r *KoperasiRepository) Delete(id uint64) error {
	return r.db.Delete(&postgres.Koperasi{}, id).Error
}

func (r *KoperasiRepository) GetByNIK(nik uint64) (*postgres.Koperasi, error) {
	var koperasi postgres.Koperasi
	err := r.db.Where("nik = ?", nik).First(&koperasi).Error
	if err != nil {
		return nil, err
	}
	return &koperasi, nil
}

func (r *KoperasiRepository) GetByNomorSK(nomorSK string) (*postgres.Koperasi, error) {
	var koperasi postgres.Koperasi
	err := r.db.Where("nomor_sk = ?", nomorSK).First(&koperasi).Error
	if err != nil {
		return nil, err
	}
	return &koperasi, nil
}

type AnggotaKoperasiRepository struct {
	db *gorm.DB
}

func NewAnggotaKoperasiRepository(db *gorm.DB) *AnggotaKoperasiRepository {
	return &AnggotaKoperasiRepository{db: db}
}

func (r *AnggotaKoperasiRepository) Create(anggota *postgres.AnggotaKoperasi) error {
	return r.db.Create(anggota).Error
}

func (r *AnggotaKoperasiRepository) GetByID(id uint64) (*postgres.AnggotaKoperasi, error) {
	var anggota postgres.AnggotaKoperasi
	err := r.db.Preload("Koperasi").Preload("Jabatan").Preload("Kelurahan").
		First(&anggota, id).Error
	if err != nil {
		return nil, err
	}
	return &anggota, nil
}

func (r *AnggotaKoperasiRepository) GetByKoperasiID(koperasiID uint64, limit, offset int) ([]postgres.AnggotaKoperasi, error) {
	var anggotas []postgres.AnggotaKoperasi
	err := r.db.Where("koperasi_id = ?", koperasiID).
		Preload("Jabatan").Preload("Kelurahan").
		Limit(limit).Offset(offset).
		Find(&anggotas).Error
	return anggotas, err
}

func (r *AnggotaKoperasiRepository) GetByNIAK(niak string) (*postgres.AnggotaKoperasi, error) {
	var anggota postgres.AnggotaKoperasi
	err := r.db.Where("niak = ?", niak).First(&anggota).Error
	if err != nil {
		return nil, err
	}
	return &anggota, nil
}

func (r *AnggotaKoperasiRepository) Update(anggota *postgres.AnggotaKoperasi) error {
	return r.db.Save(anggota).Error
}

func (r *AnggotaKoperasiRepository) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&postgres.AnggotaKoperasi{}).Where("id = ?", id).Update("status_anggota", status).Error
}

func (r *AnggotaKoperasiRepository) CountByKoperasiID(koperasiID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&postgres.AnggotaKoperasi{}).Where("koperasi_id = ?", koperasiID).Count(&count).Error
	return count, err
}

func (r *AnggotaKoperasiRepository) GetActiveByKoperasiID(koperasiID uint64) ([]postgres.AnggotaKoperasi, error) {
	var anggotas []postgres.AnggotaKoperasi
	err := r.db.Where("koperasi_id = ? AND status_anggota = ?", koperasiID, "aktif").
		Find(&anggotas).Error
	return anggotas, err
}

type WilayahRepository struct {
	db *gorm.DB
}

func NewWilayahRepository(db *gorm.DB) *WilayahRepository {
	return &WilayahRepository{db: db}
}

func (r *WilayahRepository) GetProvinsiList() ([]postgres.WilayahProvinsi, error) {
	var provinsis []postgres.WilayahProvinsi
	err := r.db.Order("nama ASC").Find(&provinsis).Error
	return provinsis, err
}

func (r *WilayahRepository) GetKabupatenByProvinsiID(provinsiID uint64) ([]postgres.WilayahKabupaten, error) {
	var kabupatens []postgres.WilayahKabupaten
	err := r.db.Where("provinsi_id = ?", provinsiID).Order("nama ASC").Find(&kabupatens).Error
	return kabupatens, err
}

func (r *WilayahRepository) GetKecamatanByKabupatenID(kabupatenID uint64) ([]postgres.WilayahKecamatan, error) {
	var kecamatans []postgres.WilayahKecamatan
	err := r.db.Where("kabupaten_id = ?", kabupatenID).Order("nama ASC").Find(&kecamatans).Error
	return kecamatans, err
}

func (r *WilayahRepository) GetKelurahanByKecamatanID(kecamatanID uint64) ([]postgres.WilayahKelurahan, error) {
	var kelurahans []postgres.WilayahKelurahan
	err := r.db.Where("kecamatan_id = ?", kecamatanID).Order("nama ASC").Find(&kelurahans).Error
	return kelurahans, err
}