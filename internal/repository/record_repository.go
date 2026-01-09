package repository

import (
	"dothefortune_server/internal/database"
	"dothefortune_server/internal/models"
)

type RecordRepository interface {
	Create(record *models.FortuneRecord) error
	FindByUserID(userID uint, limit int) ([]models.FortuneRecord, error)
	FindByUserIDAndType(userID uint, recordType string, limit int) ([]models.FortuneRecord, error)
}

type recordRepository struct{}

func NewRecordRepository() RecordRepository {
	return &recordRepository{}
}

func (r *recordRepository) Create(record *models.FortuneRecord) error {
	return database.DB.Create(record).Error
}

func (r *recordRepository) FindByUserID(userID uint, limit int) ([]models.FortuneRecord, error) {
	var records []models.FortuneRecord
	err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

func (r *recordRepository) FindByUserIDAndType(userID uint, recordType string, limit int) ([]models.FortuneRecord, error) {
	var records []models.FortuneRecord
	err := database.DB.
		Where("user_id = ? AND type = ?", userID, recordType).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

