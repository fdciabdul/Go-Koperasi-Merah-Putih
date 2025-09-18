package postgres

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *postgres.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint64) (*postgres.User, error) {
	var user postgres.User
	err := r.db.Preload("Tenant").Preload("Koperasi").Preload("Anggota").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*postgres.User, error) {
	var user postgres.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*postgres.User, error) {
	var user postgres.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *postgres.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) UpdateLastLogin(id uint64) error {
	now := time.Now()
	return r.db.Model(&postgres.User{}).Where("id = ?", id).Update("last_login", &now).Error
}

func (r *UserRepository) Delete(id uint64) error {
	return r.db.Delete(&postgres.User{}, id).Error
}

func (r *UserRepository) GetByKoperasiID(koperasiID uint64, limit, offset int) ([]postgres.User, error) {
	var users []postgres.User
	err := r.db.Where("koperasi_id = ?", koperasiID).
		Limit(limit).Offset(offset).
		Find(&users).Error
	return users, err
}

func (r *UserRepository) CountByKoperasiID(koperasiID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&postgres.User{}).Where("koperasi_id = ?", koperasiID).Count(&count).Error
	return count, err
}

type UserRegistrationRepository struct {
	db *gorm.DB
}

func NewUserRegistrationRepository(db *gorm.DB) *UserRegistrationRepository {
	return &UserRegistrationRepository{db: db}
}

func (r *UserRegistrationRepository) Create(registration *postgres.UserRegistration) error {
	return r.db.Create(registration).Error
}

func (r *UserRegistrationRepository) GetByID(id uint64) (*postgres.UserRegistration, error) {
	var registration postgres.UserRegistration
	err := r.db.Preload("Koperasi").Preload("Payment").First(&registration, id).Error
	if err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *UserRegistrationRepository) GetByPaymentID(paymentID uint64) (*postgres.UserRegistration, error) {
	var registration postgres.UserRegistration
	err := r.db.Where("payment_id = ?", paymentID).First(&registration).Error
	if err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *UserRegistrationRepository) Update(registration *postgres.UserRegistration) error {
	return r.db.Save(registration).Error
}

func (r *UserRegistrationRepository) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&postgres.UserRegistration{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRegistrationRepository) GetExpiredRegistrations() ([]postgres.UserRegistration, error) {
	var registrations []postgres.UserRegistration
	now := time.Now()
	err := r.db.Where("status = ? AND expires_at < ?", "pending_payment", now).Find(&registrations).Error
	return registrations, err
}

func (r *UserRegistrationRepository) GetPendingApproval(koperasiID uint64) ([]postgres.UserRegistration, error) {
	var registrations []postgres.UserRegistration
	err := r.db.Where("koperasi_id = ? AND status = ?", koperasiID, "payment_verified").Find(&registrations).Error
	return registrations, err
}