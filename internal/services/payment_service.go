package services

import (
	"encoding/json"
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type PaymentService struct {
	paymentRepo         *postgresRepo.PaymentRepository
	paymentProviderRepo *postgresRepo.PaymentProviderRepository
	sequenceService     *SequenceService
}

func NewPaymentService(
	paymentRepo *postgresRepo.PaymentRepository,
	paymentProviderRepo *postgresRepo.PaymentProviderRepository,
	sequenceService *SequenceService,
) *PaymentService {
	return &PaymentService{
		paymentRepo:         paymentRepo,
		paymentProviderRepo: paymentProviderRepo,
		sequenceService:     sequenceService,
	}
}

func (s *PaymentService) CreatePayment(req *CreatePaymentRequest) (*postgres.PaymentTransaction, error) {
	provider, err := s.paymentProviderRepo.GetByCode(req.ProviderCode)
	if err != nil && req.ProviderID == 0 {
		return nil, fmt.Errorf("payment provider not found: %v", err)
	}

	method, err := s.paymentProviderRepo.GetMethodByID(req.MethodID)
	if err != nil {
		return nil, fmt.Errorf("payment method not found: %v", err)
	}

	if req.Amount < method.MinimalAmount || (method.MaksimalAmount > 0 && req.Amount > method.MaksimalAmount) {
		return nil, fmt.Errorf("amount out of range: min %.2f, max %.2f", method.MinimalAmount, method.MaksimalAmount)
	}

	adminFee := s.calculateAdminFee(provider, req.Amount)
	totalAmount := req.Amount + adminFee

	nomorTransaksi, err := s.generateNomorTransaksi(req.TenantID, req.KoperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction number: %v", err)
	}

	expiredDate := time.Now().Add(24 * time.Hour)

	payment := &postgres.PaymentTransaction{
		TenantID:        req.TenantID,
		KoperasiID:      req.KoperasiID,
		NomorTransaksi:  nomorTransaksi,
		ProviderID:      provider.ID,
		MethodID:        req.MethodID,
		Amount:          req.Amount,
		AdminFee:        adminFee,
		TotalAmount:     totalAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		Description:     req.Description,
		Status:          "pending",
		ExpiredDate:     &expiredDate,
		TransactionType: req.TransactionType,
		ReferenceID:     req.ReferenceID,
		ReferenceTable:  req.ReferenceTable,
	}

	err = s.paymentRepo.CreateTransaction(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment transaction: %v", err)
	}

	switch provider.Kode {
	case "midtrans":
		err = s.createMidtransPayment(payment, method)
	case "xendit":
		err = s.createXenditPayment(payment, method)
	default:
		err = fmt.Errorf("unsupported payment provider: %s", provider.Kode)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create payment with provider: %v", err)
	}

	return payment, nil
}

func (s *PaymentService) HandleCallback(providerCode string, callbackData map[string]interface{}) error {
	var transactionID string
	var status string
	var paymentDate *time.Time

	switch providerCode {
	case "midtrans":
		transactionID, _ = callbackData["order_id"].(string)
		transactionStatus, _ := callbackData["transaction_status"].(string)
		status = s.mapMidtransStatus(transactionStatus)
		if status == "paid" {
			now := time.Now()
			paymentDate = &now
		}
	case "xendit":
		transactionID, _ = callbackData["external_id"].(string)
		callbackStatus, _ := callbackData["status"].(string)
		status = s.mapXenditStatus(callbackStatus)
		if status == "paid" {
			now := time.Now()
			paymentDate = &now
		}
	}

	payment, err := s.paymentRepo.GetTransactionByNomor(transactionID)
	if err != nil {
		return fmt.Errorf("payment transaction not found: %v", err)
	}

	callback := &postgres.PaymentCallback{
		PaymentID:    payment.ID,
		CallbackType: "webhook",
		IsValid:      true,
		ProcessedAt:  &time.Time{},
	}

	rawData, _ := json.Marshal(callbackData)
	callback.RawData = string(rawData)
	callback.ProcessedData = string(rawData)

	err = s.paymentRepo.CreateCallback(callback)
	if err != nil {
		return fmt.Errorf("failed to create callback: %v", err)
	}

	gatewayResponse, _ := json.Marshal(callbackData)
	err = s.paymentRepo.UpdateTransactionResponse(payment.ID, string(gatewayResponse), string(gatewayResponse))
	if err != nil {
		return fmt.Errorf("failed to update payment response: %v", err)
	}

	err = s.paymentRepo.UpdateTransactionStatus(payment.ID, status, paymentDate)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	return nil
}

func (s *PaymentService) ProcessExpiredPayments() error {
	expiredPayments, err := s.paymentRepo.GetExpiredTransactions()
	if err != nil {
		return fmt.Errorf("failed to get expired payments: %v", err)
	}

	for _, payment := range expiredPayments {
		err = s.paymentRepo.UpdateTransactionStatus(payment.ID, "expired", nil)
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *PaymentService) calculateAdminFee(provider *postgres.PaymentProvider, amount float64) float64 {
	switch provider.FeeType {
	case "fixed":
		return provider.FeeAmount
	case "percentage":
		return amount * provider.FeePercentage / 100
	case "both":
		return provider.FeeAmount + (amount * provider.FeePercentage / 100)
	default:
		return 0
	}
}

func (s *PaymentService) generateNomorTransaksi(tenantID, koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(tenantID, koperasiID, "payment_transaction")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PAY%04d%08d", koperasiID, number), nil
}

func (s *PaymentService) createMidtransPayment(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) error {
	return nil
}

func (s *PaymentService) createXenditPayment(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) error {
	return nil
}

func (s *PaymentService) mapMidtransStatus(status string) string {
	switch status {
	case "capture", "settlement":
		return "paid"
	case "pending":
		return "pending"
	case "deny", "cancel", "expire":
		return "failed"
	default:
		return "pending"
	}
}

func (s *PaymentService) mapXenditStatus(status string) string {
	switch status {
	case "PAID":
		return "paid"
	case "PENDING":
		return "pending"
	case "EXPIRED", "FAILED":
		return "failed"
	default:
		return "pending"
	}
}

type CreatePaymentRequest struct {
	TenantID        uint64  `json:"tenant_id"`
	KoperasiID      uint64  `json:"koperasi_id"`
	ProviderID      uint64  `json:"provider_id"`
	ProviderCode    string  `json:"provider_code"`
	MethodID        uint64  `json:"method_id"`
	Amount          float64 `json:"amount"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerPhone   string  `json:"customer_phone"`
	Description     string  `json:"description"`
	TransactionType string  `json:"transaction_type"`
	ReferenceID     uint64  `json:"reference_id"`
	ReferenceTable  string  `json:"reference_table"`
}