package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type UserService struct {
	userRepo         *postgresRepo.UserRepository
	registrationRepo *postgresRepo.UserRegistrationRepository
	anggotaRepo      *postgresRepo.AnggotaKoperasiRepository
	paymentService   *PaymentService
	sequenceService  *SequenceService
}

func NewUserService(
	userRepo *postgresRepo.UserRepository,
	registrationRepo *postgresRepo.UserRegistrationRepository,
	anggotaRepo *postgresRepo.AnggotaKoperasiRepository,
	paymentService *PaymentService,
	sequenceService *SequenceService,
) *UserService {
	return &UserService{
		userRepo:         userRepo,
		registrationRepo: registrationRepo,
		anggotaRepo:      anggotaRepo,
		paymentService:   paymentService,
		sequenceService:  sequenceService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token     string         `json:"token"`
	User      *postgres.User `json:"user"`
	ExpiresAt time.Time      `json:"expires_at"`
}

func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is not active")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	expiresAt := time.Now().Add(24 * time.Hour)
	token, err := GenerateJWT(user.ID, user.TenantID, user.Role, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &LoginResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *UserService) RegisterUser(req *UserRegistrationRequest) (*postgres.UserRegistration, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	token, err := generateVerificationToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %v", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	registration := &postgres.UserRegistration{
		KoperasiID:          req.KoperasiID,
		NIK:                 req.NIK,
		NamaLengkap:         req.NamaLengkap,
		JenisKelamin:        req.JenisKelamin,
		TempatLahir:         req.TempatLahir,
		TanggalLahir:        req.TanggalLahir,
		Alamat:              req.Alamat,
		RT:                  req.RT,
		RW:                  req.RW,
		KelurahanID:         req.KelurahanID,
		Telepon:             req.Telepon,
		Email:               req.Email,
		Username:            req.Username,
		PasswordHash:        string(hashedPassword),
		SimpananPokokAmount: req.SimpananPokokAmount,
		Status:              "pending_payment",
		VerificationToken:   token,
		ExpiresAt:           &expiresAt,
	}

	err = s.registrationRepo.Create(registration)
	if err != nil {
		return nil, fmt.Errorf("failed to create registration: %v", err)
	}

	paymentReq := &CreatePaymentRequest{
		TenantID:        req.TenantID,
		KoperasiID:      req.KoperasiID,
		ProviderID:      req.PaymentProviderID,
		MethodID:        req.PaymentMethodID,
		Amount:          req.SimpananPokokAmount,
		CustomerName:    req.NamaLengkap,
		CustomerEmail:   req.Email,
		CustomerPhone:   req.Telepon,
		Description:     fmt.Sprintf("Pembayaran Simpanan Pokok - %s", req.NamaLengkap),
		TransactionType: "simpanan_pokok",
		ReferenceID:     registration.ID,
		ReferenceTable:  "user_registrations",
	}

	payment, err := s.paymentService.CreatePayment(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	registration.PaymentID = payment.ID
	err = s.registrationRepo.Update(registration)
	if err != nil {
		return nil, fmt.Errorf("failed to update registration with payment ID: %v", err)
	}

	return registration, nil
}

func (s *UserService) VerifyPayment(paymentID uint64) error {
	registration, err := s.registrationRepo.GetByPaymentID(paymentID)
	if err != nil {
		return fmt.Errorf("registration not found for payment ID %d: %v", paymentID, err)
	}

	if registration.Status != "pending_payment" {
		return errors.New("registration is not in pending_payment status")
	}

	registration.Status = "payment_verified"
	return s.registrationRepo.Update(registration)
}

func (s *UserService) ApproveRegistration(registrationID uint64, approvedBy uint64) error {
	registration, err := s.registrationRepo.GetByID(registrationID)
	if err != nil {
		return fmt.Errorf("registration not found: %v", err)
	}

	if registration.Status != "payment_verified" {
		return errors.New("registration is not in payment_verified status")
	}

	now := time.Now()
	registration.Status = "approved"
	registration.ApprovedBy = approvedBy
	registration.ApprovedAt = &now

	err = s.registrationRepo.Update(registration)
	if err != nil {
		return fmt.Errorf("failed to update registration: %v", err)
	}

	anggota := &postgres.AnggotaKoperasi{
		KoperasiID:    registration.KoperasiID,
		NIAK:          s.generateNIAK(registration.KoperasiID),
		NIK:           registration.NIK,
		Nama:          registration.NamaLengkap,
		JenisKelamin:  registration.JenisKelamin,
		TempatLahir:   registration.TempatLahir,
		TanggalLahir:  registration.TanggalLahir,
		Alamat:        registration.Alamat,
		RT:            registration.RT,
		RW:            registration.RW,
		KelurahanID:   registration.KelurahanID,
		Telepon:       registration.Telepon,
		Email:         registration.Email,
		Posisi:        "anggota",
		TanggalMasuk:  &now,
		StatusAnggota: "aktif",
	}

	// Save anggota to database
	if err := s.anggotaRepo.Create(anggota); err != nil {
		return fmt.Errorf("failed to create anggota: %v", err)
	}

	user := &postgres.User{
		TenantID:     1,
		KoperasiID:   registration.KoperasiID,
		Username:     registration.Username,
		Email:        registration.Email,
		PasswordHash: registration.PasswordHash,
		NamaLengkap:  registration.NamaLengkap,
		Telepon:      registration.Telepon,
		Role:         "anggota",
		IsActive:     true,
	}

	return s.userRepo.Create(user)
}

func (s *UserService) RejectRegistration(registrationID uint64, rejectedBy uint64, reason string) error {
	registration, err := s.registrationRepo.GetByID(registrationID)
	if err != nil {
		return fmt.Errorf("registration not found: %v", err)
	}

	registration.Status = "rejected"
	registration.ApprovedBy = rejectedBy
	registration.RejectionReason = reason

	return s.registrationRepo.Update(registration)
}

func (s *UserService) ProcessExpiredRegistrations() error {
	expiredRegistrations, err := s.registrationRepo.GetExpiredRegistrations()
	if err != nil {
		return fmt.Errorf("failed to get expired registrations: %v", err)
	}

	for _, registration := range expiredRegistrations {
		registration.Status = "expired"
		err = s.registrationRepo.Update(&registration)
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *UserService) generateNIAK(koperasiID uint64) string {
	number, _ := s.sequenceService.GetNextNumber(1, koperasiID, "anggota")
	return fmt.Sprintf("ANG%04d%06d", koperasiID, number)
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type UserRegistrationRequest struct {
	TenantID            uint64     `json:"tenant_id"`
	KoperasiID          uint64     `json:"koperasi_id"`
	NIK                 string     `json:"nik"`
	NamaLengkap         string     `json:"nama_lengkap"`
	JenisKelamin        string     `json:"jenis_kelamin"`
	TempatLahir         string     `json:"tempat_lahir"`
	TanggalLahir        *time.Time `json:"tanggal_lahir"`
	Alamat              string     `json:"alamat"`
	RT                  string     `json:"rt"`
	RW                  string     `json:"rw"`
	KelurahanID         uint64     `json:"kelurahan_id"`
	Telepon             string     `json:"telepon"`
	Email               string     `json:"email"`
	Username            string     `json:"username"`
	Password            string     `json:"password"`
	SimpananPokokAmount float64    `json:"simpanan_pokok_amount"`
	PaymentProviderID   uint64     `json:"payment_provider_id"`
	PaymentMethodID     uint64     `json:"payment_method_id"`
}