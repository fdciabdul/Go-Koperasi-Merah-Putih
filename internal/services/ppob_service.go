package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type PPOBService struct {
	ppobRepo        *postgresRepo.PPOBRepository
	paymentService  *PaymentService
	sequenceService *SequenceService
}

func NewPPOBService(
	ppobRepo *postgresRepo.PPOBRepository,
	paymentService *PaymentService,
	sequenceService *SequenceService,
) *PPOBService {
	return &PPOBService{
		ppobRepo:        ppobRepo,
		paymentService:  paymentService,
		sequenceService: sequenceService,
	}
}

func (s *PPOBService) GetKategoriList() ([]postgres.PPOBKategori, error) {
	return s.ppobRepo.GetKategoriList()
}

func (s *PPOBService) GetProdukByKategori(kategoriID uint64) ([]postgres.PPOBProduk, error) {
	return s.ppobRepo.GetProdukByKategori(kategoriID)
}

func (s *PPOBService) CreateTransaction(req *PPOBTransactionRequest) (*postgres.PPOBTransaksi, error) {
	produk, err := s.ppobRepo.GetProdukByID(req.ProdukID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	if !produk.IsAktif {
		return nil, fmt.Errorf("product is not active")
	}

	config, err := s.ppobRepo.GetPaymentConfig(req.KoperasiID)
	if err != nil {
		config = &postgres.PPOBPaymentConfig{
			PPOBAdminFee:     5000,
			PPOBAdminFeeType: "fixed",
		}
	}

	adminFee := s.calculatePPOBAdminFee(config, produk.HargaJual)
	totalAmount := produk.HargaJual + adminFee

	nomorTransaksi, err := s.generateNomorTransaksi(req.KoperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate transaction number: %v", err)
	}

	transaksi := &postgres.PPOBTransaksi{
		KoperasiID:     req.KoperasiID,
		AnggotaID:      req.AnggotaID,
		ProdukID:       req.ProdukID,
		NomorTransaksi: nomorTransaksi,
		NomorTujuan:    req.NomorTujuan,
		NamaPelanggan:  req.NamaPelanggan,
		HargaBeli:      produk.HargaBeli,
		HargaJual:      produk.HargaJual,
		FeeAgen:        produk.FeeAgen,
		Status:         "pending",
		CustomerName:   req.CustomerName,
		CustomerEmail:  req.CustomerEmail,
		CustomerPhone:  req.CustomerPhone,
		AdminFee:       adminFee,
		PaymentStatus:  "pending",
	}

	err = s.ppobRepo.CreateTransaksi(transaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to create PPOB transaction: %v", err)
	}

	paymentReq := &CreatePaymentRequest{
		TenantID:        1,
		KoperasiID:      req.KoperasiID,
		ProviderID:      req.PaymentProviderID,
		MethodID:        req.PaymentMethodID,
		Amount:          totalAmount,
		CustomerName:    req.CustomerName,
		CustomerEmail:   req.CustomerEmail,
		CustomerPhone:   req.CustomerPhone,
		Description:     fmt.Sprintf("PPOB %s - %s", produk.NamaProduk, req.NomorTujuan),
		TransactionType: "ppob",
		ReferenceID:     transaksi.ID,
		ReferenceTable:  "ppob_transaksi",
	}

	payment, err := s.paymentService.CreatePayment(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	transaksi.PaymentID = payment.ID
	err = s.ppobRepo.UpdateTransaksiStatus(transaksi.ID, "pending", "Waiting for payment")
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction with payment ID: %v", err)
	}

	return transaksi, nil
}

func (s *PPOBService) ProcessPayment(paymentID uint64) error {
	transaksi, err := s.ppobRepo.GetTransaksiByPaymentID(paymentID)
	if err != nil {
		return fmt.Errorf("PPOB transaction not found for payment ID %d: %v", paymentID, err)
	}

	if transaksi.PaymentStatus != "pending" {
		return fmt.Errorf("payment already processed")
	}

	err = s.ppobRepo.UpdatePaymentStatus(transaksi.ID, "paid")
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	err = s.processToProvider(transaksi)
	if err != nil {
		s.ppobRepo.UpdateTransaksiStatus(transaksi.ID, "failed", err.Error())
		return fmt.Errorf("failed to process to provider: %v", err)
	}

	return s.ppobRepo.UpdateTransaksiStatus(transaksi.ID, "success", "Transaction successful")
}

func (s *PPOBService) processToProvider(transaksi *postgres.PPOBTransaksi) error {
	return nil
}

func (s *PPOBService) CreateSettlement(koperasiID uint64, dari, sampai time.Time, processedBy uint64) (*postgres.PPOBSettlement, error) {
	transaksis, err := s.ppobRepo.GetTransaksiForSettlement(koperasiID, dari, sampai)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions for settlement: %v", err)
	}

	if len(transaksis) == 0 {
		return nil, fmt.Errorf("no transactions found for settlement")
	}

	nomorSettlement, err := s.generateNomorSettlement(koperasiID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate settlement number: %v", err)
	}

	var totalOmzet, totalFeeAgen, totalAdminFee float64
	var settlementDetails []postgres.PPOBSettlementDetail

	for _, transaksi := range transaksis {
		totalOmzet += transaksi.HargaJual
		totalFeeAgen += transaksi.FeeAgen
		totalAdminFee += transaksi.AdminFee

		detail := postgres.PPOBSettlementDetail{
			PPOBTransaksiID: transaksi.ID,
			Omzet:           transaksi.HargaJual,
			FeeAgen:         transaksi.FeeAgen,
			AdminFee:        transaksi.AdminFee,
		}
		settlementDetails = append(settlementDetails, detail)
	}

	settlement := &postgres.PPOBSettlement{
		KoperasiID:        koperasiID,
		NomorSettlement:   nomorSettlement,
		TanggalSettlement: time.Now(),
		PeriodeDari:       dari,
		PeriodeSampai:     sampai,
		JumlahTransaksi:   len(transaksis),
		TotalOmzet:        totalOmzet,
		TotalFeeAgen:      totalFeeAgen,
		TotalAdminFee:     totalAdminFee,
		TotalSettlement:   totalFeeAgen + totalAdminFee,
		Status:            "draft",
		ProcessedBy:       processedBy,
	}

	err = s.ppobRepo.CreateSettlement(settlement)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement: %v", err)
	}

	for i := range settlementDetails {
		settlementDetails[i].SettlementID = settlement.ID
	}

	err = s.ppobRepo.CreateSettlementDetails(settlementDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement details: %v", err)
	}

	var transaksiIDs []uint64
	for _, transaksi := range transaksis {
		transaksiIDs = append(transaksiIDs, transaksi.ID)
	}

	err = s.ppobRepo.MarkTransaksiSettled(transaksiIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to mark transactions as settled: %v", err)
	}

	return settlement, nil
}

func (s *PPOBService) calculatePPOBAdminFee(config *postgres.PPOBPaymentConfig, amount float64) float64 {
	if config.PPOBAdminFeeType == "percentage" {
		return amount * config.PPOBAdminFee / 100
	}
	return config.PPOBAdminFee
}

func (s *PPOBService) generateNomorTransaksi(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "ppob_transaksi")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("PPOB%04d%08d", koperasiID, number), nil
}

func (s *PPOBService) generateNomorSettlement(koperasiID uint64) (string, error) {
	number, err := s.sequenceService.GetNextNumber(1, koperasiID, "ppob_settlement")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("SET%04d%06d", koperasiID, number), nil
}

type PPOBTransactionRequest struct {
	KoperasiID        uint64 `json:"koperasi_id"`
	AnggotaID         uint64 `json:"anggota_id"`
	ProdukID          uint64 `json:"produk_id"`
	NomorTujuan       string `json:"nomor_tujuan"`
	NamaPelanggan     string `json:"nama_pelanggan"`
	CustomerName      string `json:"customer_name"`
	CustomerEmail     string `json:"customer_email"`
	CustomerPhone     string `json:"customer_phone"`
	PaymentProviderID uint64 `json:"payment_provider_id"`
	PaymentMethodID   uint64 `json:"payment_method_id"`
}