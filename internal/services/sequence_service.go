package services

import (
	"fmt"
	"time"

	"koperasi-merah-putih/internal/models/postgres"
	postgresRepo "koperasi-merah-putih/internal/repository/postgres"
)

type SequenceService struct {
	sequenceRepo *postgresRepo.SequenceRepository
}

func NewSequenceService(sequenceRepo *postgresRepo.SequenceRepository) *SequenceService {
	return &SequenceService{sequenceRepo: sequenceRepo}
}

func (s *SequenceService) GetNextNumber(tenantID, koperasiID uint64, sequenceType string) (uint64, error) {
	return s.sequenceRepo.GetNextSequenceNumber(tenantID, koperasiID, sequenceType)
}

func (s *SequenceService) GetSequenceList(tenantID uint64, koperasiID *uint64) ([]postgres.Sequence, error) {
	return s.sequenceRepo.GetSequenceList(tenantID, koperasiID)
}

func (s *SequenceService) UpdateSequenceValue(tenantID, koperasiID uint64, sequenceType string, value uint64) error {
	if value < 0 {
		return fmt.Errorf("sequence value cannot be negative")
	}
	return s.sequenceRepo.UpdateSequenceValue(tenantID, koperasiID, sequenceType, value)
}

func (s *SequenceService) ResetSequence(tenantID, koperasiID uint64, sequenceType string) error {
	return s.sequenceRepo.ResetSequence(tenantID, koperasiID, sequenceType)
}

