package postgres

import (
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateTransaction(payment *postgres.PaymentTransaction) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetTransactionByID(id uint64) (*postgres.PaymentTransaction, error) {
	var payment postgres.PaymentTransaction
	err := r.db.Preload("Provider").Preload("Method").Preload("Koperasi").First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetTransactionByNomor(nomor string) (*postgres.PaymentTransaction, error) {
	var payment postgres.PaymentTransaction
	err := r.db.Where("nomor_transaksi = ?", nomor).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetTransactionByExternalID(externalID string) (*postgres.PaymentTransaction, error) {
	var payment postgres.PaymentTransaction
	err := r.db.Where("external_id = ?", externalID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdateTransactionStatus(id uint64, status string, paymentDate *time.Time) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if paymentDate != nil {
		updates["payment_date"] = paymentDate
	}
	return r.db.Model(&postgres.PaymentTransaction{}).Where("id = ?", id).Updates(updates).Error
}

func (r *PaymentRepository) UpdateTransactionResponse(id uint64, gatewayResponse, callbackData string) error {
	return r.db.Model(&postgres.PaymentTransaction{}).Where("id = ?", id).Updates(map[string]interface{}{
		"gateway_response": gatewayResponse,
		"callback_data":    callbackData,
		"updated_at":       time.Now(),
	}).Error
}

func (r *PaymentRepository) GetExpiredTransactions() ([]postgres.PaymentTransaction, error) {
	var payments []postgres.PaymentTransaction
	now := time.Now()
	err := r.db.Where("status = ? AND expired_date < ?", "pending", now).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) CreateCallback(callback *postgres.PaymentCallback) error {
	return r.db.Create(callback).Error
}

func (r *PaymentRepository) GetCallbacksByPaymentID(paymentID uint64) ([]postgres.PaymentCallback, error) {
	var callbacks []postgres.PaymentCallback
	err := r.db.Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&callbacks).Error
	return callbacks, err
}

type PaymentProviderRepository struct {
	db *gorm.DB
}

func NewPaymentProviderRepository(db *gorm.DB) *PaymentProviderRepository {
	return &PaymentProviderRepository{db: db}
}

func (r *PaymentProviderRepository) GetActiveProviders() ([]postgres.PaymentProvider, error) {
	var providers []postgres.PaymentProvider
	err := r.db.Where("is_active = ?", true).Preload("PaymentMethods", "is_active = ?", true).Find(&providers).Error
	return providers, err
}

func (r *PaymentProviderRepository) GetByCode(code string) (*postgres.PaymentProvider, error) {
	var provider postgres.PaymentProvider
	err := r.db.Where("kode = ? AND is_active = ?", code, true).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *PaymentProviderRepository) GetMethodByID(id uint64) (*postgres.PaymentMethod, error) {
	var method postgres.PaymentMethod
	err := r.db.Preload("Provider").First(&method, id).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *PaymentProviderRepository) GetMethodsByProviderID(providerID uint64) ([]postgres.PaymentMethod, error) {
	var methods []postgres.PaymentMethod
	err := r.db.Where("provider_id = ? AND is_active = ?", providerID, true).Find(&methods).Error
	return methods, err
}