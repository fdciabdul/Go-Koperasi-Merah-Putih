package postgres

import (
	"gorm.io/gorm"
	"koperasi-merah-putih/internal/models/postgres"
)

type SequenceRepository struct {
	db *gorm.DB
}

func NewSequenceRepository(db *gorm.DB) *SequenceRepository {
	return &SequenceRepository{db: db}
}

func (r *SequenceRepository) GetNextSequenceNumber(tenantID, koperasiID uint64, sequenceType string) (uint64, error) {
	var sequence postgres.Sequence

	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("tenant_id = ? AND koperasi_id = ? AND sequence_type = ?",
			tenantID, koperasiID, sequenceType).First(&sequence)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				sequence = postgres.Sequence{
					TenantID:     tenantID,
					KoperasiID:   koperasiID,
					SequenceType: sequenceType,
					CurrentValue: 1,
				}
				return tx.Create(&sequence).Error
			}
			return result.Error
		}

		sequence.CurrentValue++
		return tx.Save(&sequence).Error
	})

	if err != nil {
		return 0, err
	}

	return sequence.CurrentValue, nil
}

func (r *SequenceRepository) GetSequenceList(tenantID uint64, koperasiID *uint64) ([]postgres.Sequence, error) {
	var sequences []postgres.Sequence
	query := r.db.Where("tenant_id = ?", tenantID)

	if koperasiID != nil {
		query = query.Where("koperasi_id = ?", *koperasiID)
	}

	err := query.Order("sequence_type ASC").Find(&sequences).Error
	return sequences, err
}

func (r *SequenceRepository) UpdateSequenceValue(tenantID, koperasiID uint64, sequenceType string, value uint64) error {
	return r.db.Where("tenant_id = ? AND koperasi_id = ? AND sequence_type = ?",
		tenantID, koperasiID, sequenceType).
		Updates(map[string]interface{}{"current_value": value}).Error
}

func (r *SequenceRepository) ResetSequence(tenantID, koperasiID uint64, sequenceType string) error {
	return r.db.Where("tenant_id = ? AND koperasi_id = ? AND sequence_type = ?",
		tenantID, koperasiID, sequenceType).
		Updates(map[string]interface{}{"current_value": 0}).Error
}