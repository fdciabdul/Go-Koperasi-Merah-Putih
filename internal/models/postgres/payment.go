package postgres

import (
	"time"
)

type PaymentProvider struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Kode          string    `gorm:"uniqueIndex;size:50;not null" json:"kode"`
	Nama          string    `gorm:"size:255;not null" json:"nama"`
	Jenis         string    `gorm:"type:varchar(20);not null;index" json:"jenis"`
	BaseURL       string    `gorm:"size:500" json:"base_url"`
	MerchantID    string    `gorm:"size:255" json:"merchant_id"`
	APIKey        string    `gorm:"size:500" json:"api_key"`
	SecretKey     string    `gorm:"size:500" json:"secret_key"`
	CallbackURL   string    `gorm:"size:500" json:"callback_url"`
	FeeType       string    `gorm:"type:varchar(20);default:'percentage'" json:"fee_type"`
	FeeAmount     float64   `gorm:"type:decimal(15,2);default:0" json:"fee_amount"`
	FeePercentage float64   `gorm:"type:decimal(5,2);default:0" json:"fee_percentage"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	PaymentMethods      []PaymentMethod      `gorm:"foreignKey:ProviderID" json:"payment_methods,omitempty"`
	PaymentTransactions []PaymentTransaction `gorm:"foreignKey:ProviderID" json:"payment_transactions,omitempty"`
}

type PaymentMethod struct {
	ID             uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ProviderID     uint64  `gorm:"not null" json:"provider_id"`
	Kode           string  `gorm:"size:50;not null" json:"kode"`
	Nama           string  `gorm:"size:255;not null" json:"nama"`
	Jenis          string  `gorm:"type:varchar(20);not null;index" json:"jenis"`
	BankCode       string  `gorm:"size:10" json:"bank_code"`
	WalletCode     string  `gorm:"size:20" json:"wallet_code"`
	LogoURL        string  `gorm:"size:500" json:"logo_url"`
	MinimalAmount  float64 `gorm:"type:decimal(15,2);default:0" json:"minimal_amount"`
	MaksimalAmount float64 `gorm:"type:decimal(15,2);default:0" json:"maksimal_amount"`
	IsActive       bool    `gorm:"default:true" json:"is_active"`

	Provider                   PaymentProvider        `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	PaymentTransactions        []PaymentTransaction   `gorm:"foreignKey:MethodID" json:"payment_transactions,omitempty"`
	PPOBPaymentConfigDefault   []PPOBPaymentConfig    `gorm:"foreignKey:DefaultPaymentMethodID" json:"ppob_payment_config_default,omitempty"`
}

type PaymentTransaction struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID         uint64     `gorm:"not null" json:"tenant_id"`
	KoperasiID       uint64     `gorm:"index" json:"koperasi_id"`
	NomorTransaksi   string     `gorm:"uniqueIndex;size:50;not null" json:"nomor_transaksi"`
	ExternalID       string     `gorm:"size:100;index" json:"external_id"`
	InvoiceID        string     `gorm:"size:100" json:"invoice_id"`
	ProviderID       uint64     `gorm:"not null" json:"provider_id"`
	MethodID         uint64     `gorm:"not null" json:"method_id"`
	Amount           float64    `gorm:"type:decimal(15,2);not null" json:"amount"`
	AdminFee         float64    `gorm:"type:decimal(15,2);default:0" json:"admin_fee"`
	TotalAmount      float64    `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	CustomerName     string     `gorm:"size:255" json:"customer_name"`
	CustomerEmail    string     `gorm:"size:255" json:"customer_email"`
	CustomerPhone    string     `gorm:"size:20" json:"customer_phone"`
	Description      string     `gorm:"type:text" json:"description"`
	Status           string     `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	PaymentDate      *time.Time `json:"payment_date"`
	ExpiredDate      *time.Time `json:"expired_date"`
	GatewayResponse  string     `gorm:"type:json" json:"gateway_response"`
	CallbackData     string     `gorm:"type:json" json:"callback_data"`
	TransactionType  string     `gorm:"type:varchar(20);not null" json:"transaction_type"`
	ReferenceID      uint64     `json:"reference_id"`
	ReferenceTable   string     `gorm:"size:100" json:"reference_table"`
	CreatedAt        time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Tenant                     Tenant                       `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Koperasi                   Koperasi                     `gorm:"foreignKey:KoperasiID" json:"koperasi,omitempty"`
	Provider                   PaymentProvider              `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`
	Method                     PaymentMethod                `gorm:"foreignKey:MethodID" json:"method,omitempty"`
	PaymentCallbacks           []PaymentCallback            `gorm:"foreignKey:PaymentID" json:"payment_callbacks,omitempty"`
	PPOBTransaksi              []PPOBTransaksi              `gorm:"foreignKey:PaymentID" json:"ppob_transaksi,omitempty"`
	SimpananPokokTransaksi     []SimpananPokokTransaksi     `gorm:"foreignKey:PaymentID" json:"simpanan_pokok_transaksi,omitempty"`
	UserRegistrations          []UserRegistration           `gorm:"foreignKey:PaymentID" json:"user_registrations,omitempty"`
}

type PaymentCallback struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentID     uint64     `gorm:"not null;index" json:"payment_id"`
	CallbackType  string     `gorm:"type:varchar(20);not null" json:"callback_type"`
	RawData       string     `gorm:"type:text" json:"raw_data"`
	ProcessedData string     `gorm:"type:json" json:"processed_data"`
	Signature     string     `gorm:"size:500" json:"signature"`
	IsValid       bool       `gorm:"default:false" json:"is_valid"`
	ProcessedAt   *time.Time `json:"processed_at"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`

	Payment PaymentTransaction `gorm:"foreignKey:PaymentID" json:"payment,omitempty"`
}