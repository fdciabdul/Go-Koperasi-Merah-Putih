package postgres

import (
	"time"
)

type User struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID     uint64     `gorm:"not null;index" json:"tenant_id"`
	KoperasiID   uint64     `gorm:"index" json:"koperasi_id"`
	Username     string     `gorm:"uniqueIndex;size:100;not null" json:"username"`
	Email        string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"password_hash"`
	NamaLengkap  string     `gorm:"size:255;not null" json:"nama_lengkap"`
	Telepon      string     `gorm:"size:20" json:"telepon"`
	Role         string     `gorm:"type:varchar(20);not null" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	AnggotaID    uint64     `json:"anggota_id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Tenant  Tenant          `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi Koperasi       `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota AnggotaKoperasi `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	PPOBSettlementProcessed []PPOBSettlement `gorm:"foreignKey:ProcessedBy" json:"ppob_settlement_processed,omitempty"`
	UserRegistrationApproved []UserRegistration `gorm:"foreignKey:ApprovedBy" json:"user_registration_approved,omitempty"`
	UserRegistrationLogs []UserRegistrationLog `gorm:"foreignKey:CreatedBy" json:"user_registration_logs,omitempty"`
	SimpananPokokTransaksi []SimpananPokokTransaksi `gorm:"foreignKey:CreatedBy" json:"simpanan_pokok_transaksi,omitempty"`
}

type Permission struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Module      string    `gorm:"size:50;not null" json:"module"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
}

type RolePermission struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Role         string `gorm:"type:varchar(20);not null" json:"role"`
	PermissionID uint64 `gorm:"not null" json:"permission_id"`

	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

type UserRegistration struct {
	ID                   uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID           uint64     `gorm:"not null" json:"koperasi_id"`
	NIK                  string     `gorm:"size:20;not null" json:"nik"`
	NamaLengkap          string     `gorm:"size:255;not null" json:"nama_lengkap"`
	JenisKelamin         string     `gorm:"type:varchar(1);not null" json:"jenis_kelamin"`
	TempatLahir          string     `gorm:"size:100" json:"tempat_lahir"`
	TanggalLahir         *time.Time `json:"tanggal_lahir"`
	Alamat               string     `gorm:"type:text" json:"alamat"`
	RT                   string     `gorm:"size:5" json:"rt"`
	RW                   string     `gorm:"size:5" json:"rw"`
	KelurahanID          uint64     `json:"kelurahan_id"`
	Telepon              string     `gorm:"size:20;not null" json:"telepon"`
	Email                string     `gorm:"size:255;not null" json:"email"`
	Username             string     `gorm:"size:100;not null;uniqueIndex" json:"username"`
	PasswordHash         string     `gorm:"size:255;not null" json:"password_hash"`
	SimpananPokokAmount  float64    `gorm:"type:decimal(15,2);not null" json:"simpanan_pokok_amount"`
	PaymentID            uint64     `json:"payment_id"`
	Status               string     `gorm:"type:varchar(20);default:'pending_payment';index" json:"status"`
	VerificationToken    string     `gorm:"size:100" json:"verification_token"`
	ApprovedBy           uint64     `json:"approved_by"`
	ApprovedAt           *time.Time `json:"approved_at"`
	RejectionReason      string     `gorm:"type:text" json:"rejection_reason"`
	ExpiresAt            *time.Time `gorm:"index" json:"expires_at"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Koperasi                   Koperasi                     `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Kelurahan                  WilayahKelurahan             `gorm:"foreignKey:KelurahanID" json:"kelurahan,omitempty"`
	Payment                    PaymentTransaction           `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	ApprovedByUser             User                         `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	UserRegistrationLogs       []UserRegistrationLog        `gorm:"foreignKey:RegistrationID" json:"user_registration_logs,omitempty"`
	SimpananPokokTransaksi     []SimpananPokokTransaksi     `gorm:"foreignKey:RegistrationID" json:"simpanan_pokok_transaksi,omitempty"`
}

type UserRegistrationLog struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	RegistrationID uint64    `gorm:"not null;index" json:"registration_id"`
	Action         string    `gorm:"type:varchar(20);not null" json:"action"`
	Description    string    `gorm:"type:text" json:"description"`
	Metadata       string    `gorm:"type:json" json:"metadata"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy      uint64    `json:"created_by"`

	Registration UserRegistration `gorm:"foreignKey:RegistrationID" json:"registration,omitempty"`
	CreatedByUser User            `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

type SimpananPokokConfig struct {
	ID                     uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID             uint64    `gorm:"not null;uniqueIndex" json:"koperasi_id"`
	JumlahSimpananPokok    float64   `gorm:"type:decimal(15,2);not null" json:"jumlah_simpanan_pokok"`
	IsWajib                bool      `gorm:"default:true" json:"is_wajib"`
	AllowedPaymentMethods  string    `gorm:"type:json" json:"allowed_payment_methods"`
	PaymentDeadlineDays    int       `gorm:"default:7" json:"payment_deadline_days"`
	AkunSimpananPokokID    uint64    `json:"akun_simpanan_pokok_id"`
	IsActive               bool      `gorm:"default:true" json:"is_active"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Koperasi            Koperasi                 `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	AkunSimpananPokok   COAAkun                  `gorm:"foreignKey:AkunSimpananPokokID" json:"akun_simpanan_pokok,omitempty"`
	SimpananPokokTransaksi []SimpananPokokTransaksi `gorm:"foreignKey:KoperasiID" json:"simpanan_pokok_transaksi,omitempty"`
}

type SimpananPokokTransaksi struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	KoperasiID       uint64     `gorm:"not null" json:"koperasi_id"`
	AnggotaID        uint64     `json:"anggota_id"`
	RegistrationID   uint64     `gorm:"index" json:"registration_id"`
	NomorTransaksi   string     `gorm:"size:50;not null" json:"nomor_transaksi"`
	Jumlah           float64    `gorm:"type:decimal(15,2);not null" json:"jumlah"`
	PaymentID        uint64     `json:"payment_id"`
	Status           string     `gorm:"type:varchar(10);default:'pending';index" json:"status"`
	TanggalTransaksi time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_transaksi"`
	TanggalLunas     *time.Time `json:"tanggal_lunas"`
	JurnalID         uint64     `json:"jurnal_id"`
	Keterangan       string     `gorm:"type:text" json:"keterangan"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy        uint64     `json:"created_by"`

	Koperasi     Koperasi           `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Anggota      AnggotaKoperasi    `gorm:"foreignKey:AnggotaID" json:"anggota,omitempty"`
	Registration UserRegistration   `gorm:"foreignKey:RegistrationID" json:"registration,omitempty"`
	Payment      PaymentTransaction `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
	Jurnal       JurnalUmum         `gorm:"foreignKey:JurnalID" json:"jurnal,omitempty"`
	CreatedByUser User              `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}