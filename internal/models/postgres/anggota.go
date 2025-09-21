package postgres

import (
	"time"

	"gorm.io/gorm"
)

type JabatanKoperasi struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Tingkat   string         `gorm:"type:varchar(20);not null" json:"tingkat"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	AnggotaKoperasi []AnggotaKoperasi `gorm:"foreignKey:JabatanID" json:"anggota_koperasi,omitempty"`
}

type AnggotaKoperasi struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID     uint64     `gorm:"not null;index" json:"koperasi_id"`
	NIAK           string     `gorm:"uniqueIndex;size:20;not null" json:"niak"`
	NIK            string     `gorm:"size:20;index" json:"nik"`
	Nama           string     `gorm:"size:255;not null" json:"nama"`
	JenisKelamin   string     `gorm:"type:varchar(1);not null" json:"jenis_kelamin"`
	TempatLahir    string     `gorm:"size:100" json:"tempat_lahir"`
	TanggalLahir   *time.Time `json:"tanggal_lahir"`
	Alamat         string     `gorm:"type:text" json:"alamat"`
	RT             string     `gorm:"size:5" json:"rt"`
	RW             string     `gorm:"size:5" json:"rw"`
	KelurahanID    uint64     `json:"kelurahan_id"`
	Telepon        string     `gorm:"size:20" json:"telepon"`
	Email          string     `gorm:"size:255" json:"email"`
	Posisi         string     `gorm:"type:varchar(20);default:'anggota'" json:"posisi"`
	JabatanID      uint64     `json:"jabatan_id"`
	TanggalMasuk   *time.Time `json:"tanggal_masuk"`
	TanggalKeluar  *time.Time `json:"tanggal_keluar"`
	StatusAnggota  string     `gorm:"type:varchar(20);default:'aktif';index" json:"status_anggota"`
	NPWP           string     `gorm:"size:20" json:"npwp"`
	Pekerjaan      string     `gorm:"size:100" json:"pekerjaan"`
	Pendidikan     string     `gorm:"size:50" json:"pendidikan"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Koperasi               Koperasi                  `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Jabatan                JabatanKoperasi           `gorm:"foreignKey:JabatanID" json:"jabatan,omitempty"`
	Kelurahan              WilayahKelurahan          `gorm:"foreignKey:KelurahanID" json:"kelurahan,omitempty"`
	RekeningSimPan         []RekeningSimpanPinjam    `gorm:"foreignKey:AnggotaID" json:"rekening_simpan_pinjam,omitempty"`
	PPOBTransaksi          []PPOBTransaksi           `gorm:"foreignKey:AnggotaID" json:"ppob_transaksi,omitempty"`
	SimpananPokokTransaksi []SimpananPokokTransaksi  `gorm:"foreignKey:AnggotaID" json:"simpanan_pokok_transaksi,omitempty"`
	Users                  []User                    `gorm:"foreignKey:AnggotaID" json:"users,omitempty"`
}