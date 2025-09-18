package postgres

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantCode string         `gorm:"uniqueIndex;size:20;not null" json:"tenant_code"`
	TenantName string         `gorm:"size:255;not null" json:"tenant_name"`
	Domain     string         `gorm:"size:100" json:"domain"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Koperasi []Koperasi `gorm:"foreignKey:TenantID" json:"koperasi,omitempty"`
	Users    []User     `gorm:"foreignKey:TenantID" json:"users,omitempty"`
}

type WilayahProvinsi struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Kode string `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama string `gorm:"size:255;not null" json:"nama"`

	Kabupaten []WilayahKabupaten `gorm:"foreignKey:ProvinsiID" json:"kabupaten,omitempty"`
}

type WilayahKabupaten struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	Kode       string `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama       string `gorm:"size:255;not null" json:"nama"`
	ProvinsiID uint64 `gorm:"not null" json:"provinsi_id"`

	Provinsi  WilayahProvinsi   `gorm:"foreignKey:ProvinsiID" json:"provinsi,omitempty"`
	Kecamatan []WilayahKecamatan `gorm:"foreignKey:KabupatenID" json:"kecamatan,omitempty"`
}

type WilayahKecamatan struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Kode        string `gorm:"uniqueIndex;size:10;not null" json:"kode"`
	Nama        string `gorm:"size:255;not null" json:"nama"`
	KabupatenID uint64 `gorm:"not null" json:"kabupaten_id"`

	Kabupaten WilayahKabupaten   `gorm:"foreignKey:KabupatenID" json:"kabupaten,omitempty"`
	Kelurahan []WilayahKelurahan `gorm:"foreignKey:KecamatanID" json:"kelurahan,omitempty"`
}

type WilayahKelurahan struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	Kode        string `gorm:"uniqueIndex;size:15;not null" json:"kode"`
	Nama        string `gorm:"size:255;not null" json:"nama"`
	KecamatanID uint64 `gorm:"not null" json:"kecamatan_id"`
	Jenis       string `gorm:"type:enum('kelurahan','desa');default:'desa'" json:"jenis"`

	Kecamatan      WilayahKecamatan     `gorm:"foreignKey:KecamatanID" json:"kecamatan,omitempty"`
	AnggotaKoperasi []AnggotaKoperasi   `gorm:"foreignKey:KelurahanID" json:"anggota_koperasi,omitempty"`
	UserRegistrations []UserRegistration `gorm:"foreignKey:KelurahanID" json:"user_registrations,omitempty"`
}