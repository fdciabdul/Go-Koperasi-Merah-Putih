package postgres

import (
	"time"
)

type KlinikPasien struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID     uint64     `gorm:"not null" json:"koperasi_id"`
	NomorRM        string     `gorm:"size:20;not null" json:"nomor_rm"`
	NIK            string     `gorm:"size:20;index" json:"nik"`
	NamaLengkap    string     `gorm:"size:255;not null;index" json:"nama_lengkap"`
	JenisKelamin   string     `gorm:"type:varchar(1);not null" json:"jenis_kelamin"`
	TempatLahir    string     `gorm:"size:100" json:"tempat_lahir"`
	TanggalLahir   *time.Time `json:"tanggal_lahir"`
	Alamat         string     `gorm:"type:text" json:"alamat"`
	Telepon        string     `gorm:"size:20" json:"telepon"`
	Email          string     `gorm:"size:100" json:"email"`
	GolonganDarah  string     `gorm:"type:varchar(3)" json:"golongan_darah"`
	Alergi         string     `gorm:"type:text" json:"alergi"`
	RiwayatPenyakit string    `gorm:"type:text" json:"riwayat_penyakit"`
	AnggotaID      uint64     `json:"anggota_id"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Koperasi        Koperasi          `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota         AnggotaKoperasi   `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	KlinikKunjungan []KlinikKunjungan `gorm:"foreignKey:PasienID" json:"klinik_kunjungan,omitempty"`
}

type KlinikTenagaMedis struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64    `gorm:"not null" json:"koperasi_id"`
	NIP               string    `gorm:"size:30" json:"nip"`
	NamaLengkap       string    `gorm:"size:255;not null;index" json:"nama_lengkap"`
	JenisKelamin      string    `gorm:"type:varchar(1);not null" json:"jenis_kelamin"`
	Spesialisasi      string    `gorm:"size:100;index" json:"spesialisasi"`
	NoSTR             string    `gorm:"size:50" json:"no_str"`
	NoSIP             string    `gorm:"size:50" json:"no_sip"`
	Telepon           string    `gorm:"size:20" json:"telepon"`
	Email             string    `gorm:"size:100" json:"email"`
	JadwalPraktik     string    `gorm:"type:json" json:"jadwal_praktik"`
	TarifKonsultasi   float64   `gorm:"type:decimal(15,2);default:0" json:"tarif_konsultasi"`
	Status            string    `gorm:"type:varchar(20);default:'aktif'" json:"status"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`

	Koperasi        Koperasi          `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	KlinikKunjungan []KlinikKunjungan `gorm:"foreignKey:DokterID" json:"klinik_kunjungan,omitempty"`
}

type KlinikKunjungan struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID        uint64    `gorm:"not null" json:"koperasi_id"`
	PasienID          uint64    `gorm:"not null;index" json:"pasien_id"`
	DokterID          uint64    `gorm:"not null;index" json:"dokter_id"`
	NomorKunjungan    string    `gorm:"size:30;not null" json:"nomor_kunjungan"`
	TanggalKunjungan  time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"tanggal_kunjungan"`
	KeluhanUtama      string    `gorm:"type:text" json:"keluhan_utama"`
	Anamnesis         string    `gorm:"type:text" json:"anamnesis"`
	PemeriksaanFisik  string    `gorm:"type:text" json:"pemeriksaan_fisik"`
	Diagnosis         string    `gorm:"type:text" json:"diagnosis"`
	TerapiPengobatan  string    `gorm:"type:text" json:"terapi_pengobatan"`
	BiayaKonsultasi   float64   `gorm:"type:decimal(15,2);default:0" json:"biaya_konsultasi"`
	BiayaTindakan     float64   `gorm:"type:decimal(15,2);default:0" json:"biaya_tindakan"`
	BiayaObat         float64   `gorm:"type:decimal(15,2);default:0" json:"biaya_obat"`
	TotalBiaya        float64   `gorm:"type:decimal(15,2);default:0" json:"total_biaya"`
	StatusPembayaran  string    `gorm:"type:varchar(20);default:'belum_bayar'" json:"status_pembayaran"`
	JurnalID          uint64    `json:"jurnal_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`

	Koperasi     Koperasi            `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Pasien       KlinikPasien        `gorm:"foreignKey:PasienID" json:"pasien,omitempty"`
	Dokter       KlinikTenagaMedis   `gorm:"foreignKey:DokterID" json:"dokter,omitempty"`
	Jurnal       JurnalUmum          `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
	KlinikResep  []KlinikResep       `gorm:"foreignKey:KunjunganID" json:"klinik_resep,omitempty"`
}

type KlinikObat struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID      uint64    `gorm:"not null" json:"koperasi_id"`
	KodeObat        string    `gorm:"size:20;not null" json:"kode_obat"`
	NamaObat        string    `gorm:"size:255;not null;index" json:"nama_obat"`
	Kategori        string    `gorm:"size:100" json:"kategori"`
	BentukSediaan   string    `gorm:"size:100" json:"bentuk_sediaan"`
	Kekuatan        string    `gorm:"size:50" json:"kekuatan"`
	Satuan          string    `gorm:"size:20" json:"satuan"`
	StokMinimal     int       `gorm:"default:0" json:"stok_minimal"`
	StokCurrent     int       `gorm:"default:0" json:"stok_current"`
	HargaBeli       float64   `gorm:"type:decimal(15,2);default:0" json:"harga_beli"`
	HargaJual       float64   `gorm:"type:decimal(15,2);default:0" json:"harga_jual"`
	IsAktif         bool      `gorm:"default:true" json:"is_aktif"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`

	Koperasi    Koperasi      `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	KlinikResep []KlinikResep `gorm:"foreignKey:ObatID" json:"klinik_resep,omitempty"`
}

type KlinikResep struct {
	ID           uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	KunjunganID  uint64  `gorm:"not null;index" json:"kunjungan_id"`
	ObatID       uint64  `gorm:"not null" json:"obat_id"`
	Jumlah       int     `gorm:"not null" json:"jumlah"`
	AturanPakai  string  `gorm:"size:255" json:"aturan_pakai"`
	Keterangan   string  `gorm:"size:255" json:"keterangan"`
	HargaSatuan  float64 `gorm:"type:decimal(15,2);default:0" json:"harga_satuan"`
	TotalHarga   float64 `gorm:"type:decimal(15,2);default:0" json:"total_harga"`

	Kunjungan KlinikKunjungan `gorm:"foreignKey:KunjunganID" json:"kunjungan,omitempty"`
	Obat      KlinikObat      `gorm:"foreignKey:ObatID" json:"obat,omitempty"`
}