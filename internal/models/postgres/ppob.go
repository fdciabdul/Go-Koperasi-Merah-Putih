package postgres

import (
	"time"

	"gorm.io/gorm"
)

type PPOBKategori struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string         `gorm:"uniqueIndex;size:20;not null" json:"kode"`
	Nama      string         `gorm:"size:100;not null" json:"nama"`
	Icon      string         `gorm:"size:255" json:"icon"`
	Urutan    int            `gorm:"default:0" json:"urutan"`
	IsAktif   bool           `gorm:"default:true" json:"is_aktif"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	PPOBProduk []PPOBProduk `gorm:"foreignKey:KategoriID" json:"ppob_produk,omitempty"`
}

type PPOBProvider struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode      string    `gorm:"uniqueIndex;size:50;not null" json:"kode"`
	Nama      string    `gorm:"size:255;not null" json:"nama"`
	BaseURL   string    `gorm:"size:500" json:"base_url"`
	APIKey    string    `gorm:"size:500" json:"api_key"`
	SecretKey string    `gorm:"size:500" json:"secret_key"`
	IsAktif   bool      `gorm:"default:true" json:"is_aktif"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	PPOBProduk []PPOBProduk `gorm:"foreignKey:ProviderID" json:"ppob_produk,omitempty"`
}

type PPOBProduk struct {
	ID               uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderID       uint64  `gorm:"not null" json:"provider_id"`
	KategoriID       uint64  `gorm:"not null;index" json:"kategori_id"`
	KodeProduk       string  `gorm:"size:50;not null" json:"kode_produk"`
	NamaProduk       string  `gorm:"size:255;not null" json:"nama_produk"`
	Deskripsi        string  `gorm:"type:text" json:"deskripsi"`
	HargaBeli        float64 `gorm:"type:decimal(15,2);default:0" json:"harga_beli"`
	HargaJual        float64 `gorm:"type:decimal(15,2);default:0" json:"harga_jual"`
	FeeAgen          float64 `gorm:"type:decimal(15,2);default:0" json:"fee_agen"`
	IsAktif          bool    `gorm:"default:true" json:"is_aktif"`
	ValidasiFormat   string  `gorm:"size:100" json:"validasi_format"`

	Provider      PPOBProvider    `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Kategori      PPOBKategori    `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	PPOBTransaksi []PPOBTransaksi `gorm:"foreignKey:ProdukID" json:"ppob_transaksi,omitempty"`
}

type PPOBTransaksi struct {
	ID                 uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID         uint64     `gorm:"not null;index" json:"koperasi_id"`
	AnggotaID          uint64     `gorm:"index" json:"anggota_id"`
	ProdukID           uint64     `gorm:"not null" json:"produk_id"`
	NomorTransaksi     string     `gorm:"uniqueIndex;size:50;not null" json:"nomor_transaksi"`
	NomorReferensi     string     `gorm:"size:100" json:"nomor_referensi"`
	NomorTujuan        string     `gorm:"size:50;not null" json:"nomor_tujuan"`
	NamaPelanggan      string     `gorm:"size:255" json:"nama_pelanggan"`
	HargaBeli          float64    `gorm:"type:decimal(15,2);not null" json:"harga_beli"`
	HargaJual          float64    `gorm:"type:decimal(15,2);not null" json:"harga_jual"`
	FeeAgen            float64    `gorm:"type:decimal(15,2);default:0" json:"fee_agen"`
	Status             string     `gorm:"type:enum('pending','success','failed','cancelled');default:'pending';index" json:"status"`
	PesanResponse      string     `gorm:"type:text" json:"pesan_response"`
	TanggalTransaksi   time.Time  `gorm:"default:CURRENT_TIMESTAMP;index" json:"tanggal_transaksi"`
	TanggalSettlement  *time.Time `json:"tanggal_settlement"`
	JurnalID           uint64     `json:"jurnal_id"`
	PaymentID          uint64     `json:"payment_id"`
	PaymentStatus      string     `gorm:"type:enum('pending','paid','failed');default:'pending';index" json:"payment_status"`
	CustomerName       string     `gorm:"size:255" json:"customer_name"`
	CustomerEmail      string     `gorm:"size:255" json:"customer_email"`
	CustomerPhone      string     `gorm:"size:20" json:"customer_phone"`
	AdminFee           float64    `gorm:"type:decimal(15,2);default:0" json:"admin_fee"`

	Koperasi        Koperasi           `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota         AnggotaKoperasi    `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	Produk          PPOBProduk         `gorm:"foreignKey:ProdukID" json:"produk,omitempty"`
	Jurnal          JurnalUmum         `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
	Payment         PaymentTransaction `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	SettlementDetails []PPOBSettlementDetail `gorm:"foreignKey:PPOBTransaksiID" json:"settlement_details,omitempty"`
}

type PPOBPaymentConfig struct {
	ID                      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID              uint64 `gorm:"not null;uniqueIndex" json:"koperasi_id"`
	AllowedPaymentMethods   string `gorm:"type:json" json:"allowed_payment_methods"`
	DefaultPaymentMethodID  uint64 `json:"default_payment_method_id"`
	AutoSettlement          bool   `gorm:"default:false" json:"auto_settlement"`
	SettlementSchedule      string `gorm:"type:enum('immediate','daily','weekly','monthly');default:'daily'" json:"settlement_schedule"`
	PPOBAdminFee            float64 `gorm:"type:decimal(15,2);default:0" json:"ppob_admin_fee"`
	PPOBAdminFeeType        string `gorm:"type:enum('fixed','percentage');default:'fixed'" json:"ppob_admin_fee_type"`
	IsActive                bool   `gorm:"default:true" json:"is_active"`

	Koperasi              Koperasi       `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	DefaultPaymentMethod  PaymentMethod  `gorm:"foreignKey:DefaultPaymentMethodID" json:"default_payment_method,omitempty"`
}

type PPOBSettlement struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID       uint64     `gorm:"not null" json:"koperasi_id"`
	NomorSettlement  string     `gorm:"size:50;not null" json:"nomor_settlement"`
	TanggalSettlement time.Time  `gorm:"not null;index" json:"tanggal_settlement"`
	PeriodeDari      time.Time  `gorm:"not null" json:"periode_dari"`
	PeriodeSampai    time.Time  `gorm:"not null" json:"periode_sampai"`
	JumlahTransaksi  int        `gorm:"default:0" json:"jumlah_transaksi"`
	TotalOmzet       float64    `gorm:"type:decimal(15,2);default:0" json:"total_omzet"`
	TotalFeeAgen     float64    `gorm:"type:decimal(15,2);default:0" json:"total_fee_agen"`
	TotalAdminFee    float64    `gorm:"type:decimal(15,2);default:0" json:"total_admin_fee"`
	TotalSettlement  float64    `gorm:"type:decimal(15,2);default:0" json:"total_settlement"`
	Status           string     `gorm:"type:enum('draft','processed','paid');default:'draft';index" json:"status"`
	JurnalID         uint64     `json:"jurnal_id"`
	ProcessedAt      *time.Time `json:"processed_at"`
	ProcessedBy      uint64     `json:"processed_by"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Koperasi          Koperasi               `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Jurnal            JurnalUmum             `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
	ProcessedByUser   User                   `gorm:"foreignKey:ProcessedBy" json:"processed_by_user,omitempty"`
	SettlementDetails []PPOBSettlementDetail `gorm:"foreignKey:SettlementID" json:"settlement_details,omitempty"`
}

type PPOBSettlementDetail struct {
	ID               uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SettlementID     uint64  `gorm:"not null" json:"settlement_id"`
	PPOBTransaksiID  uint64  `gorm:"not null" json:"ppob_transaksi_id"`
	Omzet            float64 `gorm:"type:decimal(15,2);not null" json:"omzet"`
	FeeAgen          float64 `gorm:"type:decimal(15,2);not null" json:"fee_agen"`
	AdminFee         float64 `gorm:"type:decimal(15,2);not null" json:"admin_fee"`

	Settlement    PPOBSettlement  `gorm:"foreignKey:SettlementID" json:"settlement,omitempty"`
	PPOBTransaksi PPOBTransaksi   `gorm:"foreignKey:PPOBTransaksiID" json:"ppob_transaksi,omitempty"`
}