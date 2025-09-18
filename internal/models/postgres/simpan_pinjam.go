package postgres

import (
	"time"

	"gorm.io/gorm"
)

type ProdukSimpanPinjam struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64         `gorm:"not null" json:"koperasi_id"`
	KodeProduk        string         `gorm:"size:20;not null" json:"kode_produk"`
	NamaProduk        string         `gorm:"size:255;not null" json:"nama_produk"`
	Jenis             string         `gorm:"type:enum('simpanan','pinjaman');not null" json:"jenis"`
	Kategori          string         `gorm:"size:100" json:"kategori"`
	BungaSimpanan     float64        `gorm:"type:decimal(5,2);default:0" json:"bunga_simpanan"`
	MinimalSaldo      float64        `gorm:"type:decimal(15,2);default:0" json:"minimal_saldo"`
	BungaPinjaman     float64        `gorm:"type:decimal(5,2);default:0" json:"bunga_pinjaman"`
	BungaDenda        float64        `gorm:"type:decimal(5,2);default:0" json:"bunga_denda"`
	MaksimalPinjaman  float64        `gorm:"type:decimal(15,2);default:0" json:"maksimal_pinjaman"`
	JangkaWaktuMax    int            `gorm:"default:0" json:"jangka_waktu_max"`
	SyaratKetentuan   string         `gorm:"type:text" json:"syarat_ketentuan"`
	IsAktif           bool           `gorm:"default:true" json:"is_aktif"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi              Koperasi                `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	RekeningSimPan        []RekeningSimpanPinjam  `gorm:"foreignKey:ProdukID" json:"rekening_simpan_pinjam,omitempty"`
}

type RekeningSimpanPinjam struct {
	ID                    uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID            uint64     `gorm:"not null" json:"koperasi_id"`
	AnggotaID             uint64     `gorm:"not null;index" json:"anggota_id"`
	ProdukID              uint64     `gorm:"not null" json:"produk_id"`
	NomorRekening         string     `gorm:"size:50;not null" json:"nomor_rekening"`
	SaldoSimpanan         float64    `gorm:"type:decimal(15,2);default:0" json:"saldo_simpanan"`
	PokokPinjaman         float64    `gorm:"type:decimal(15,2);default:0" json:"pokok_pinjaman"`
	SisaPokok             float64    `gorm:"type:decimal(15,2);default:0" json:"sisa_pokok"`
	BungaBerjalan         float64    `gorm:"type:decimal(15,2);default:0" json:"bunga_berjalan"`
	DendaKeterlambatan    float64    `gorm:"type:decimal(15,2);default:0" json:"denda_keterlambatan"`
	TanggalMulai          *time.Time `json:"tanggal_mulai"`
	TanggalJatuhTempo     *time.Time `json:"tanggal_jatuh_tempo"`
	JangkaWaktu           int        `json:"jangka_waktu"`
	AngsuranPokok         float64    `gorm:"type:decimal(15,2);default:0" json:"angsuran_pokok"`
	AngsuranBunga         float64    `gorm:"type:decimal(15,2);default:0" json:"angsuran_bunga"`
	Status                string     `gorm:"type:enum('aktif','lunas','macet','tutup');default:'aktif';index" json:"status"`
	TanggalBuka           time.Time  `gorm:"default:CURRENT_DATE" json:"tanggal_buka"`
	TanggalTutup          *time.Time `json:"tanggal_tutup"`
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Koperasi            Koperasi                  `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota             AnggotaKoperasi           `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	Produk              ProdukSimpanPinjam        `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	TransaksiSimpanPinjam []TransaksiSimpanPinjam `gorm:"foreignKey:RekeningID" json:"transaksi_simpan_pinjam,omitempty"`
}

type TransaksiSimpanPinjam struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64    `gorm:"not null" json:"koperasi_id"`
	RekeningID        uint64    `gorm:"not null;index" json:"rekening_id"`
	NomorTransaksi    string    `gorm:"size:50;not null" json:"nomor_transaksi"`
	TanggalTransaksi  time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"tanggal_transaksi"`
	JenisTransaksi    string    `gorm:"type:enum('setoran','penarikan','pencairan','angsuran','bunga','denda');not null;index" json:"jenis_transaksi"`
	Jumlah            float64   `gorm:"type:decimal(15,2);not null" json:"jumlah"`
	SaldoSebelum      float64   `gorm:"type:decimal(15,2);default:0" json:"saldo_sebelum"`
	SaldoSesudah      float64   `gorm:"type:decimal(15,2);default:0" json:"saldo_sesudah"`
	Keterangan        string    `gorm:"size:255" json:"keterangan"`
	Referensi         string    `gorm:"size:100" json:"referensi"`
	JurnalID          uint64    `json:"jurnal_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy         uint64    `json:"created_by"`

	Koperasi Koperasi             `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Rekening RekeningSimpanPinjam `gorm:"foreignKey:RekeningID" json:"rekening,omitempty"`
	Jurnal   JurnalUmum           `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
}