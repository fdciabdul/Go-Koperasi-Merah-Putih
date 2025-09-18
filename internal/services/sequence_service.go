package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type SequenceService struct {
	db *gorm.DB
}

func NewSequenceService(db *gorm.DB) *SequenceService {
	return &SequenceService{db: db}
}

func (s *SequenceService) GetNextNumber(tenantID, koperasiID uint64, sequenceName string) (uint64, error) {
	var sequence postgres.SequenceNumber

	err := s.db.Where("tenant_id = ? AND koperasi_id = ? AND sequence_name = ?",
		tenantID, koperasiID, sequenceName).First(&sequence).Error

	if err == gorm.ErrRecordNotFound {
		sequence = postgres.SequenceNumber{
			TenantID:      tenantID,
			KoperasiID:    koperasiID,
			SequenceName:  sequenceName,
			CurrentNumber: 1,
			IncrementBy:   1,
			ResetPeriod:   "never",
		}
		err = s.db.Create(&sequence).Error
		if err != nil {
			return 0, fmt.Errorf("failed to create sequence: %v", err)
		}
		return 1, nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to get sequence: %v", err)
	}

	shouldReset, err := s.shouldResetSequence(&sequence)
	if err != nil {
		return 0, fmt.Errorf("failed to check reset condition: %v", err)
	}

	if shouldReset {
		sequence.CurrentNumber = 1
		now := time.Now()
		sequence.LastResetDate = &now
	} else {
		sequence.CurrentNumber += uint64(sequence.IncrementBy)
	}

	err = s.db.Save(&sequence).Error
	if err != nil {
		return 0, fmt.Errorf("failed to update sequence: %v", err)
	}

	return sequence.CurrentNumber, nil
}

func (s *SequenceService) shouldResetSequence(sequence *postgres.SequenceNumber) (bool, error) {
	if sequence.ResetPeriod == "never" {
		return false, nil
	}

	if sequence.LastResetDate == nil {
		return true, nil
	}

	now := time.Now()
	lastReset := *sequence.LastResetDate

	switch sequence.ResetPeriod {
	case "daily":
		return !isSameDay(lastReset, now), nil
	case "monthly":
		return !isSameMonth(lastReset, now), nil
	case "yearly":
		return lastReset.Year() != now.Year(), nil
	default:
		return false, nil
	}
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}