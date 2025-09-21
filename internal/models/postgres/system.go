package postgres

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint64    `gorm:"not null" json:"tenant_id"`
	KoperasiID  uint64    `json:"koperasi_id"`
	UserID      uint64    `gorm:"index" json:"user_id"`
	TableName   string    `gorm:"size:100;not null;index" json:"table_name"`
	RecordID    uint64    `gorm:"not null;index" json:"record_id"`
	Action      string    `gorm:"type:varchar(10);not null" json:"action"`
	OldValues   string    `gorm:"type:json" json:"old_values"`
	NewValues   string    `gorm:"type:json" json:"new_values"`
	IPAddress   string    `gorm:"size:45" json:"ip_address"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	Tenant   Tenant   `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type SystemSetting struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID    uint64         `json:"tenant_id"`
	KoperasiID  uint64         `json:"koperasi_id"`
	KeyName     string         `gorm:"size:100;not null" json:"key_name"`
	KeyValue    string         `gorm:"type:text" json:"key_value"`
	DataType    string         `gorm:"type:varchar(10);default:'string'" json:"data_type"`
	Category    string         `gorm:"size:50;default:'general';index" json:"category"`
	Description string         `gorm:"type:text" json:"description"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Tenant   Tenant   `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
}

type SequenceNumber struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID        uint64    `gorm:"not null" json:"tenant_id"`
	KoperasiID      uint64    `json:"koperasi_id"`
	SequenceName    string    `gorm:"size:50;not null" json:"sequence_name"`
	Prefix          string    `gorm:"size:10" json:"prefix"`
	CurrentNumber   uint64    `gorm:"default:0" json:"current_number"`
	IncrementBy     int       `gorm:"default:1" json:"increment_by"`
	ResetPeriod     string    `gorm:"type:varchar(10);default:'never'" json:"reset_period"`
	LastResetDate   *time.Time `json:"last_reset_date"`

	Tenant   Tenant   `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi Koperasi `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
}