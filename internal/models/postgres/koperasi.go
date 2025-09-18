package postgres

import (
	"time"

	"gorm.io/gorm"
)

type JenisKoperasi struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama      string         `gorm:"size:255;not null" json:"nama"`
	Deskripsi string         `gorm:"type:text" json:"deskripsi"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi []Koperasi `gorm:"foreignKey:JenisKoperasiID" json:"koperasi,omitempty"`
}

type BentukKoperasi struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Deskripsi string         `gorm:"type:text" json:"deskripsi"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi []Koperasi `gorm:"foreignKey:BentukKoperasiID" json:"koperasi,omitempty"`
}

type StatusKoperasi struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Deskripsi string         `gorm:"type:text" json:"deskripsi"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi []Koperasi `gorm:"foreignKey:StatusKoperasiID" json:"koperasi,omitempty"`
}

type KBLI struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama      string         `gorm:"type:text;not null" json:"nama"`
	Kategori  string         `gorm:"size:100" json:"kategori"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	KoperasiAktivitasUsaha []KoperasiAktivitasUsaha `gorm:"foreignKey:KBLIID" json:"koperasi_aktivitas_usaha,omitempty"`
}

type Koperasi struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID             uint64    `gorm:"not null;index" json:"tenant_id"`
	NomorSK              string    `gorm:"uniqueIndex;size:100;not null" json:"nomor_sk"`
	NIK                  uint64    `gorm:"uniqueIndex;not null" json:"nik"`
	NamaKoperasi         string    `gorm:"size:255;not null" json:"nama_koperasi"`
	NamaSK               string    `gorm:"size:255;not null" json:"nama_sk"`
	JenisKoperasiID      uint64    `json:"jenis_koperasi_id"`
	BentukKoperasiID     uint64    `json:"bentuk_koperasi_id"`
	StatusKoperasiID     uint64    `gorm:"index" json:"status_koperasi_id"`
	ProvinsiID           uint64    `json:"provinsi_id"`
	KabupatenID          uint64    `json:"kabupaten_id"`
	KecamatanID          uint64    `json:"kecamatan_id"`
	KelurahanID          uint64    `json:"kelurahan_id"`
	Alamat               string    `gorm:"type:text" json:"alamat"`
	RT                   string    `gorm:"size:5" json:"rt"`
	RW                   string    `gorm:"size:5" json:"rw"`
	KodePos              string    `gorm:"size:10" json:"kode_pos"`
	Email                string    `gorm:"size:255" json:"email"`
	Telepon              string    `gorm:"size:20" json:"telepon"`
	Website              string    `gorm:"size:255" json:"website"`
	TanggalBerdiri       *time.Time `json:"tanggal_berdiri"`
	TanggalSK            *time.Time `json:"tanggal_sk"`
	TanggalPengesahan    *time.Time `json:"tanggal_pengesahan"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy            uint64    `json:"created_by"`
	UpdatedBy            uint64    `json:"updated_by"`

	Tenant             Tenant                     `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	JenisKoperasi      JenisKoperasi              `gorm:"foreignKey:JenisKoperasiID" json:"jenis_koperasi,omitempty"`
	BentukKoperasi     BentukKoperasi             `gorm:"foreignKey:BentukKoperasiID" json:"bentuk_koperasi,omitempty"`
	StatusKoperasi     StatusKoperasi             `gorm:"foreignKey:StatusKoperasiID" json:"status_koperasi,omitempty"`
	Provinsi           WilayahProvinsi            `gorm:"foreignKey:ProvinsiID" json:"provinsi,omitempty"`
	Kabupaten          WilayahKabupaten           `gorm:"foreignKey:KabupatenID" json:"kabupaten,omitempty"`
	Kecamatan          WilayahKecamatan           `gorm:"foreignKey:KecamatanID" json:"kecamatan,omitempty"`
	Kelurahan          WilayahKelurahan           `gorm:"foreignKey:KelurahanID" json:"kelurahan,omitempty"`
	AktivitasUsaha     []KoperasiAktivitasUsaha   `gorm:"foreignKey:KoperasiID" json:"aktivitas_usaha,omitempty"`
	AnggotaKoperasi    []AnggotaKoperasi          `gorm:"foreignKey:KoperasiID" json:"anggota_koperasi,omitempty"`
	ModalKoperasi      []ModalKoperasi            `gorm:"foreignKey:KoperasiID" json:"modal_koperasi,omitempty"`
	COAAkun            []COAAkun                  `gorm:"foreignKey:KoperasiID" json:"coa_akun,omitempty"`
	JurnalUmum         []JurnalUmum               `gorm:"foreignKey:KoperasiID" json:"jurnal_umum,omitempty"`
	ProdukSimpanPinjam []ProdukSimpanPinjam       `gorm:"foreignKey:KoperasiID" json:"produk_simpan_pinjam,omitempty"`
	PPOBTransaksi      []PPOBTransaksi            `gorm:"foreignKey:KoperasiID" json:"ppob_transaksi,omitempty"`
	PaymentTransactions []PaymentTransaction      `gorm:"foreignKey:KoperasiID" json:"payment_transactions,omitempty"`
	SimpananPokokConfig SimpananPokokConfig       `gorm:"foreignKey:KoperasiID" json:"simpanan_pokok_config,omitempty"`
}

type KoperasiAktivitasUsaha struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID   uint64    `gorm:"not null" json:"koperasi_id"`
	KBLIID       uint64    `gorm:"not null" json:"kbli_id"`
	JenisUsaha   string    `gorm:"type:enum('utama','sampingan');default:'utama'" json:"jenis_usaha"`
	Keterangan   string    `gorm:"type:text" json:"keterangan"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	KBLI     KBLI     `gorm:"foreignKey:KBLIID" json:"kbli,omitempty"`
}