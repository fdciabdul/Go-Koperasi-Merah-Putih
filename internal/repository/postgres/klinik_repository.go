package postgres

import (
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type KlinikRepository struct {
	db *gorm.DB
}

func NewKlinikRepository(db *gorm.DB) *KlinikRepository {
	return &KlinikRepository{db: db}
}

func (r *KlinikRepository) CreatePasien(pasien *postgres.KlinikPasien) error {
	return r.db.Create(pasien).Error
}

func (r *KlinikRepository) GetPasienByID(id uint64) (*postgres.KlinikPasien, error) {
	var pasien postgres.KlinikPasien
	err := r.db.Preload("Koperasi").Preload("Anggota").First(&pasien, id).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *KlinikRepository) GetPasienByNomorRM(nomorRM string) (*postgres.KlinikPasien, error) {
	var pasien postgres.KlinikPasien
	err := r.db.Where("nomor_rm = ?", nomorRM).First(&pasien).Error
	if err != nil {
		return nil, err
	}
	return &pasien, nil
}

func (r *KlinikRepository) GetPasienByKoperasi(koperasiID uint64, limit, offset int) ([]postgres.KlinikPasien, error) {
	var pasiens []postgres.KlinikPasien
	err := r.db.Where("koperasi_id = ?", koperasiID).
		Limit(limit).Offset(offset).
		Find(&pasiens).Error
	return pasiens, err
}

func (r *KlinikRepository) UpdatePasien(pasien *postgres.KlinikPasien) error {
	return r.db.Save(pasien).Error
}

func (r *KlinikRepository) SearchPasien(koperasiID uint64, search string) ([]postgres.KlinikPasien, error) {
	var pasiens []postgres.KlinikPasien
	err := r.db.Where("koperasi_id = ? AND (nama_lengkap ILIKE ? OR nomor_rm ILIKE ? OR nik ILIKE ?)",
		koperasiID, "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Limit(10).Find(&pasiens).Error
	return pasiens, err
}

func (r *KlinikRepository) CreateTenagaMedis(tenagaMedis *postgres.KlinikTenagaMedis) error {
	return r.db.Create(tenagaMedis).Error
}

func (r *KlinikRepository) GetTenagaMedisByID(id uint64) (*postgres.KlinikTenagaMedis, error) {
	var tenagaMedis postgres.KlinikTenagaMedis
	err := r.db.Preload("Koperasi").First(&tenagaMedis, id).Error
	if err != nil {
		return nil, err
	}
	return &tenagaMedis, nil
}

func (r *KlinikRepository) GetTenagaMedisByKoperasi(koperasiID uint64) ([]postgres.KlinikTenagaMedis, error) {
	var tenagaMedis []postgres.KlinikTenagaMedis
	err := r.db.Where("koperasi_id = ? AND status = ?", koperasiID, "aktif").Find(&tenagaMedis).Error
	return tenagaMedis, err
}

func (r *KlinikRepository) UpdateTenagaMedis(tenagaMedis *postgres.KlinikTenagaMedis) error {
	return r.db.Save(tenagaMedis).Error
}

func (r *KlinikRepository) CreateKunjungan(kunjungan *postgres.KlinikKunjungan) error {
	return r.db.Create(kunjungan).Error
}

func (r *KlinikRepository) GetKunjunganByID(id uint64) (*postgres.KlinikKunjungan, error) {
	var kunjungan postgres.KlinikKunjungan
	err := r.db.Preload("Koperasi").Preload("Pasien").Preload("Dokter").
		Preload("KlinikResep").Preload("KlinikResep.Obat").
		First(&kunjungan, id).Error
	if err != nil {
		return nil, err
	}
	return &kunjungan, nil
}

func (r *KlinikRepository) GetKunjunganByPasien(pasienID uint64, limit, offset int) ([]postgres.KlinikKunjungan, error) {
	var kunjungans []postgres.KlinikKunjungan
	err := r.db.Where("pasien_id = ?", pasienID).
		Preload("Dokter").Order("tanggal_kunjungan DESC").
		Limit(limit).Offset(offset).
		Find(&kunjungans).Error
	return kunjungans, err
}

func (r *KlinikRepository) GetKunjunganByKoperasi(koperasiID uint64, dari, sampai time.Time) ([]postgres.KlinikKunjungan, error) {
	var kunjungans []postgres.KlinikKunjungan
	err := r.db.Where("koperasi_id = ? AND tanggal_kunjungan BETWEEN ? AND ?",
		koperasiID, dari, sampai).
		Preload("Pasien").Preload("Dokter").
		Order("tanggal_kunjungan DESC").
		Find(&kunjungans).Error
	return kunjungans, err
}

func (r *KlinikRepository) UpdateKunjungan(kunjungan *postgres.KlinikKunjungan) error {
	return r.db.Save(kunjungan).Error
}

func (r *KlinikRepository) CreateObat(obat *postgres.KlinikObat) error {
	return r.db.Create(obat).Error
}

func (r *KlinikRepository) GetObatByID(id uint64) (*postgres.KlinikObat, error) {
	var obat postgres.KlinikObat
	err := r.db.Preload("Koperasi").First(&obat, id).Error
	if err != nil {
		return nil, err
	}
	return &obat, nil
}

func (r *KlinikRepository) GetObatByKoperasi(koperasiID uint64) ([]postgres.KlinikObat, error) {
	var obats []postgres.KlinikObat
	err := r.db.Where("koperasi_id = ? AND is_aktif = ?", koperasiID, true).
		Order("nama_obat ASC").Find(&obats).Error
	return obats, err
}

func (r *KlinikRepository) UpdateObat(obat *postgres.KlinikObat) error {
	return r.db.Save(obat).Error
}

func (r *KlinikRepository) SearchObat(koperasiID uint64, search string) ([]postgres.KlinikObat, error) {
	var obats []postgres.KlinikObat
	err := r.db.Where("koperasi_id = ? AND is_aktif = ? AND (nama_obat ILIKE ? OR kode_obat ILIKE ?)",
		koperasiID, true, "%"+search+"%", "%"+search+"%").
		Limit(10).Find(&obats).Error
	return obats, err
}

func (r *KlinikRepository) CreateResep(reseps []postgres.KlinikResep) error {
	return r.db.Create(&reseps).Error
}

func (r *KlinikRepository) GetResepByKunjungan(kunjunganID uint64) ([]postgres.KlinikResep, error) {
	var reseps []postgres.KlinikResep
	err := r.db.Where("kunjungan_id = ?", kunjunganID).
		Preload("Obat").Find(&reseps).Error
	return reseps, err
}

func (r *KlinikRepository) UpdateStokObat(obatID uint64, jumlah int) error {
	return r.db.Model(&postgres.KlinikObat{}).Where("id = ?", obatID).
		Update("stok_current", gorm.Expr("stok_current - ?", jumlah)).Error
}

func (r *KlinikRepository) GetStatistikKlinik(koperasiID uint64) (*KlinikStatistik, error) {
	var statistik KlinikStatistik

	err := r.db.Model(&postgres.KlinikKunjungan{}).
		Select(`
			COUNT(*) as total_kunjungan,
			SUM(total_biaya) as total_pendapatan,
			AVG(total_biaya) as rata_rata_biaya,
			COUNT(DISTINCT pasien_id) as total_pasien_aktif
		`).
		Where("koperasi_id = ? AND DATE(tanggal_kunjungan) >= DATE(NOW() - INTERVAL '30 days')", koperasiID).
		Scan(&statistik).Error

	return &statistik, err
}

func (r *KlinikRepository) GetObatStokRendah(koperasiID uint64) ([]postgres.KlinikObat, error) {
	var obats []postgres.KlinikObat
	err := r.db.Where("koperasi_id = ? AND is_aktif = ? AND stok_current <= stok_minimal",
		koperasiID, true).Find(&obats).Error
	return obats, err
}

type KlinikStatistik struct {
	TotalKunjungan     uint64  `json:"total_kunjungan"`
	TotalPendapatan    float64 `json:"total_pendapatan"`
	RataRataBiaya      float64 `json:"rata_rata_biaya"`
	TotalPasienAktif   uint64  `json:"total_pasien_aktif"`
}