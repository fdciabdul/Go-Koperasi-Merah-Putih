package postgres

import (
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type MasterDataRepository struct {
	db *gorm.DB
}

func NewMasterDataRepository(db *gorm.DB) *MasterDataRepository {
	return &MasterDataRepository{db: db}
}

func (r *MasterDataRepository) CreateKBLI(kbli *postgres.KBLI) error {
	return r.db.Create(kbli).Error
}

func (r *MasterDataRepository) GetKBLIList(search string, limit, offset int) ([]postgres.KBLI, error) {
	var kblis []postgres.KBLI
	query := r.db.Where("is_aktif = ?", true)

	if search != "" {
		query = query.Where("kode ILIKE ? OR nama ILIKE ? OR kategori ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	err := query.Order("kode ASC").Limit(limit).Offset(offset).Find(&kblis).Error
	return kblis, err
}

func (r *MasterDataRepository) GetKBLIByID(id uint64) (*postgres.KBLI, error) {
	var kbli postgres.KBLI
	err := r.db.First(&kbli, id).Error
	if err != nil {
		return nil, err
	}
	return &kbli, nil
}

func (r *MasterDataRepository) GetKBLIByKode(kode string) (*postgres.KBLI, error) {
	var kbli postgres.KBLI
	err := r.db.Where("kode = ?", kode).First(&kbli).Error
	if err != nil {
		return nil, err
	}
	return &kbli, nil
}

func (r *MasterDataRepository) UpdateKBLI(kbli *postgres.KBLI) error {
	return r.db.Save(kbli).Error
}

func (r *MasterDataRepository) CreateJenisKoperasi(jenis *postgres.JenisKoperasi) error {
	return r.db.Create(jenis).Error
}

func (r *MasterDataRepository) GetJenisKoperasiList() ([]postgres.JenisKoperasi, error) {
	var jenisKoperasis []postgres.JenisKoperasi
	err := r.db.Where("is_aktif = ?", true).Order("kode ASC").Find(&jenisKoperasis).Error
	return jenisKoperasis, err
}

func (r *MasterDataRepository) GetJenisKoperasiByID(id uint64) (*postgres.JenisKoperasi, error) {
	var jenis postgres.JenisKoperasi
	err := r.db.First(&jenis, id).Error
	if err != nil {
		return nil, err
	}
	return &jenis, nil
}

func (r *MasterDataRepository) GetJenisKoperasiByKode(kode string) (*postgres.JenisKoperasi, error) {
	var jenis postgres.JenisKoperasi
	err := r.db.Where("kode = ?", kode).First(&jenis).Error
	if err != nil {
		return nil, err
	}
	return &jenis, nil
}

func (r *MasterDataRepository) UpdateJenisKoperasi(jenis *postgres.JenisKoperasi) error {
	return r.db.Save(jenis).Error
}

func (r *MasterDataRepository) CreateBentukKoperasi(bentuk *postgres.BentukKoperasi) error {
	return r.db.Create(bentuk).Error
}

func (r *MasterDataRepository) GetBentukKoperasiList() ([]postgres.BentukKoperasi, error) {
	var bentukKoperasis []postgres.BentukKoperasi
	err := r.db.Where("is_aktif = ?", true).Order("kode ASC").Find(&bentukKoperasis).Error
	return bentukKoperasis, err
}

func (r *MasterDataRepository) GetBentukKoperasiByID(id uint64) (*postgres.BentukKoperasi, error) {
	var bentuk postgres.BentukKoperasi
	err := r.db.First(&bentuk, id).Error
	if err != nil {
		return nil, err
	}
	return &bentuk, nil
}

func (r *MasterDataRepository) GetBentukKoperasiByKode(kode string) (*postgres.BentukKoperasi, error) {
	var bentuk postgres.BentukKoperasi
	err := r.db.Where("kode = ?", kode).First(&bentuk).Error
	if err != nil {
		return nil, err
	}
	return &bentuk, nil
}

func (r *MasterDataRepository) UpdateBentukKoperasi(bentuk *postgres.BentukKoperasi) error {
	return r.db.Save(bentuk).Error
}