package postgres

import (
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type FinancialRepository struct {
	db *gorm.DB
}

func NewFinancialRepository(db *gorm.DB) *FinancialRepository {
	return &FinancialRepository{db: db}
}

func (r *FinancialRepository) CreateCOAKategori(kategori *postgres.COAKategori) error {
	return r.db.Create(kategori).Error
}

func (r *FinancialRepository) GetCOAKategoriList() ([]postgres.COAKategori, error) {
	var kategoris []postgres.COAKategori
	err := r.db.Order("urutan ASC, nama ASC").Find(&kategoris).Error
	return kategoris, err
}

func (r *FinancialRepository) CreateCOAAkun(akun *postgres.COAAkun) error {
	return r.db.Create(akun).Error
}

func (r *FinancialRepository) GetCOAAkunByKoperasi(koperasiID uint64) ([]postgres.COAAkun, error) {
	var akuns []postgres.COAAkun
	err := r.db.Where("koperasi_id = ? AND is_aktif = ?", koperasiID, true).
		Preload("Kategori").Preload("Parent").
		Order("kode_akun ASC").Find(&akuns).Error
	return akuns, err
}

func (r *FinancialRepository) GetCOAAkunByID(id uint64) (*postgres.COAAkun, error) {
	var akun postgres.COAAkun
	err := r.db.Preload("Kategori").Preload("Parent").Preload("Children").
		First(&akun, id).Error
	if err != nil {
		return nil, err
	}
	return &akun, nil
}

func (r *FinancialRepository) GetCOAAkunByKode(koperasiID uint64, kode string) (*postgres.COAAkun, error) {
	var akun postgres.COAAkun
	err := r.db.Where("koperasi_id = ? AND kode_akun = ?", koperasiID, kode).
		First(&akun).Error
	if err != nil {
		return nil, err
	}
	return &akun, nil
}

func (r *FinancialRepository) UpdateCOAAkun(akun *postgres.COAAkun) error {
	return r.db.Save(akun).Error
}

func (r *FinancialRepository) CreateJurnalUmum(jurnal *postgres.JurnalUmum) error {
	return r.db.Create(jurnal).Error
}

func (r *FinancialRepository) CreateJurnalDetail(details []postgres.JurnalDetail) error {
	return r.db.Create(&details).Error
}

func (r *FinancialRepository) GetJurnalUmumByID(id uint64) (*postgres.JurnalUmum, error) {
	var jurnal postgres.JurnalUmum
	err := r.db.Preload("JurnalDetail").Preload("JurnalDetail.Akun").
		First(&jurnal, id).Error
	if err != nil {
		return nil, err
	}
	return &jurnal, nil
}

func (r *FinancialRepository) GetJurnalUmumByKoperasi(koperasiID uint64, dari, sampai time.Time, limit, offset int) ([]postgres.JurnalUmum, error) {
	var jurnals []postgres.JurnalUmum
	err := r.db.Where("koperasi_id = ? AND tanggal_transaksi BETWEEN ? AND ?",
		koperasiID, dari, sampai).
		Order("tanggal_transaksi DESC, nomor_jurnal DESC").
		Limit(limit).Offset(offset).
		Find(&jurnals).Error
	return jurnals, err
}

func (r *FinancialRepository) UpdateJurnalStatus(id uint64, status string, postedBy uint64) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":    status,
		"posted_by": postedBy,
		"posted_at": &now,
	}
	return r.db.Model(&postgres.JurnalUmum{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinancialRepository) GetSaldoAkun(akunID uint64, sampaiTanggal time.Time) (float64, error) {
	var result struct {
		Saldo float64
	}

	err := r.db.Table("jurnal_detail jd").
		Select("SUM(jd.debit - jd.kredit) as saldo").
		Joins("JOIN jurnal_umum ju ON jd.jurnal_id = ju.id").
		Where("jd.akun_id = ? AND ju.status = 'posted' AND ju.tanggal_transaksi <= ?",
			akunID, sampaiTanggal).
		Scan(&result).Error

	return result.Saldo, err
}

func (r *FinancialRepository) GetNeracaSaldo(koperasiID uint64, tanggal time.Time) ([]NeracaSaldoItem, error) {
	var items []NeracaSaldoItem

	err := r.db.Table("coa_akun ca").
		Select(`
			ca.id as akun_id,
			ca.kode_akun,
			ca.nama_akun,
			cat.tipe as kategori_tipe,
			ca.saldo_normal,
			COALESCE(SUM(jd.debit), 0) as total_debit,
			COALESCE(SUM(jd.kredit), 0) as total_kredit,
			CASE
				WHEN ca.saldo_normal = 'debit' THEN COALESCE(SUM(jd.debit - jd.kredit), 0)
				ELSE COALESCE(SUM(jd.kredit - jd.debit), 0)
			END as saldo
		`).
		Joins("JOIN coa_kategori cat ON ca.kategori_id = cat.id").
		Joins("LEFT JOIN jurnal_detail jd ON ca.id = jd.akun_id").
		Joins("LEFT JOIN jurnal_umum ju ON jd.jurnal_id = ju.id AND ju.status = 'posted' AND ju.tanggal_transaksi <= ?", tanggal).
		Where("ca.koperasi_id = ? AND ca.is_aktif = ?", koperasiID, true).
		Group("ca.id, ca.kode_akun, ca.nama_akun, cat.tipe, ca.saldo_normal").
		Order("ca.kode_akun").
		Scan(&items).Error

	return items, err
}

func (r *FinancialRepository) GetLabaRugi(koperasiID uint64, dari, sampai time.Time) (*LabaRugi, error) {
	var labaRugi LabaRugi

	err := r.db.Table("jurnal_detail jd").
		Select(`
			SUM(CASE WHEN cat.tipe = 'pendapatan' THEN jd.kredit - jd.debit ELSE 0 END) as total_pendapatan,
			SUM(CASE WHEN cat.tipe = 'beban' THEN jd.debit - jd.kredit ELSE 0 END) as total_beban
		`).
		Joins("JOIN jurnal_umum ju ON jd.jurnal_id = ju.id").
		Joins("JOIN coa_akun ca ON jd.akun_id = ca.id").
		Joins("JOIN coa_kategori cat ON ca.kategori_id = cat.id").
		Where("ju.koperasi_id = ? AND ju.status = 'posted' AND ju.tanggal_transaksi BETWEEN ? AND ?",
			koperasiID, dari, sampai).
		Scan(&labaRugi).Error

	labaRugi.LabaRugi = labaRugi.TotalPendapatan - labaRugi.TotalBeban

	return &labaRugi, err
}

func (r *FinancialRepository) GetNeraca(koperasiID uint64, tanggal time.Time) (*Neraca, error) {
	var neraca Neraca

	err := r.db.Table("jurnal_detail jd").
		Select(`
			SUM(CASE WHEN cat.tipe = 'aset' THEN jd.debit - jd.kredit ELSE 0 END) as total_aset,
			SUM(CASE WHEN cat.tipe = 'kewajiban' THEN jd.kredit - jd.debit ELSE 0 END) as total_kewajiban,
			SUM(CASE WHEN cat.tipe = 'ekuitas' THEN jd.kredit - jd.debit ELSE 0 END) as total_ekuitas
		`).
		Joins("JOIN jurnal_umum ju ON jd.jurnal_id = ju.id").
		Joins("JOIN coa_akun ca ON jd.akun_id = ca.id").
		Joins("JOIN coa_kategori cat ON ca.kategori_id = cat.id").
		Where("ju.koperasi_id = ? AND ju.status = 'posted' AND ju.tanggal_transaksi <= ?",
			koperasiID, tanggal).
		Scan(&neraca).Error

	return &neraca, err
}

type NeracaSaldoItem struct {
	AkunID        uint64  `json:"akun_id"`
	KodeAkun      string  `json:"kode_akun"`
	NamaAkun      string  `json:"nama_akun"`
	KategoriTipe  string  `json:"kategori_tipe"`
	SaldoNormal   string  `json:"saldo_normal"`
	TotalDebit    float64 `json:"total_debit"`
	TotalKredit   float64 `json:"total_kredit"`
	Saldo         float64 `json:"saldo"`
}

type LabaRugi struct {
	TotalPendapatan float64 `json:"total_pendapatan"`
	TotalBeban      float64 `json:"total_beban"`
	LabaRugi        float64 `json:"laba_rugi"`
}

type Neraca struct {
	TotalAset      float64 `json:"total_aset"`
	TotalKewajiban float64 `json:"total_kewajiban"`
	TotalEkuitas   float64 `json:"total_ekuitas"`
}