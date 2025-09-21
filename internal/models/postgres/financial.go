package postgres

import (
	"time"

	"gorm.io/gorm"
)

type COAKategori struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Tipe      string    `gorm:"type:varchar(20);not null" json:"tipe"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	COAAkun []COAAkun `gorm:"foreignKey:KategoriID" json:"coa_akun,omitempty"`
}

type COAAkun struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID     uint64    `gorm:"not null" json:"tenant_id"`
	KoperasiID   uint64    `json:"koperasi_id"`
	KodeAkun     string    `gorm:"size:20;not null" json:"kode_akun"`
	NamaAkun     string    `gorm:"size:255;not null" json:"nama_akun"`
	KategoriID   uint64    `gorm:"not null;index" json:"kategori_id"`
	ParentID     uint64    `gorm:"index" json:"parent_id"`
	LevelAkun    int       `gorm:"default:1" json:"level_akun"`
	SaldoNormal  string    `gorm:"type:varchar(10);not null" json:"saldo_normal"`
	IsKas        bool      `gorm:"default:false" json:"is_kas"`
	IsAktif      bool      `gorm:"default:true" json:"is_aktif"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Tenant       Tenant        `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi     Koperasi      `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Kategori     COAKategori   `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Parent       *COAAkun      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children     []COAAkun     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	JurnalDetail []JurnalDetail `gorm:"foreignKey:AkunID" json:"jurnal_detail,omitempty"`
	SimpananPokokConfig []SimpananPokokConfig `gorm:"foreignKey:AkunSimpananPokokID" json:"simpanan_pokok_config,omitempty"`
}

type ModalKoperasi struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64     `gorm:"not null;index" json:"koperasi_id"`
	JenisModal        string     `gorm:"type:varchar(30);not null" json:"jenis_modal"`
	Jumlah            float64    `gorm:"type:decimal(15,2);default:0" json:"jumlah"`
	Keterangan        string     `gorm:"type:text" json:"keterangan"`
	TanggalPencatatan *time.Time `json:"tanggal_pencatatan"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
}

type JurnalUmum struct {
	ID               uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID         uint64      `gorm:"not null" json:"tenant_id"`
	KoperasiID       uint64      `gorm:"not null" json:"koperasi_id"`
	NomorJurnal      string      `gorm:"size:50;not null" json:"nomor_jurnal"`
	TanggalTransaksi time.Time   `gorm:"not null;index" json:"tanggal_transaksi"`
	Referensi        string      `gorm:"size:100" json:"referensi"`
	Keterangan       string      `gorm:"type:text" json:"keterangan"`
	TotalDebit       float64     `gorm:"type:decimal(15,2);default:0" json:"total_debit"`
	TotalKredit      float64     `gorm:"type:decimal(15,2);default:0" json:"total_kredit"`
	Status           string      `gorm:"type:varchar(20);default:'draft';index" json:"status"`
	CreatedAt        time.Time   `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy        uint64      `json:"created_by"`
	PostedAt         *time.Time  `json:"posted_at"`
	PostedBy         uint64      `json:"posted_by"`

	Tenant       Tenant         `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi     Koperasi       `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	JurnalDetail []JurnalDetail `gorm:"foreignKey:JurnalID" json:"jurnal_detail,omitempty"`
	PPOBTransaksi []PPOBTransaksi `gorm:"foreignKey:JurnalID" json:"ppob_transaksi,omitempty"`
	TransaksiSimpanPinjam []TransaksiSimpanPinjam `gorm:"foreignKey:JurnalID" json:"transaksi_simpan_pinjam,omitempty"`
	KlinikKunjungan []KlinikKunjungan `gorm:"foreignKey:JurnalID" json:"klinik_kunjungan,omitempty"`
	SimpananPokokTransaksi []SimpananPokokTransaksi `gorm:"foreignKey:JurnalID" json:"simpanan_pokok_transaksi,omitempty"`
	PPOBSettlement []PPOBSettlement `gorm:"foreignKey:JurnalID" json:"ppob_settlement,omitempty"`
}

type JurnalDetail struct {
	ID         uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	JurnalID   uint64  `gorm:"not null;index" json:"jurnal_id"`
	AkunID     uint64  `gorm:"not null;index" json:"akun_id"`
	Keterangan string  `gorm:"size:255" json:"keterangan"`
	Debit      float64 `gorm:"type:decimal(15,2);default:0" json:"debit"`
	Kredit     float64 `gorm:"type:decimal(15,2);default:0" json:"kredit"`

	Jurnal JurnalUmum `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
	Akun   COAAkun    `gorm:"foreignKey:AkunID" json:"akun,omitempty"`
}