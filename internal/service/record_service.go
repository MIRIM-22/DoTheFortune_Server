package service

import (
	"errors"
	"dothefortune_server/internal/models"
	"dothefortune_server/internal/repository"
)

type RecordService interface {
	GetRecentRecords(userID uint, limit int) ([]models.FortuneRecord, error)
	GetRecordsByType(userID uint, recordType string, limit int) ([]models.FortuneRecord, error)
	GetSpouseImage(userID uint) (string, error)
}

type recordService struct {
	recordRepo  repository.RecordRepository
	fortuneRepo repository.FortuneRepository
}

func NewRecordService(recordRepo repository.RecordRepository, fortuneRepo repository.FortuneRepository) RecordService {
	return &recordService{
		recordRepo:  recordRepo,
		fortuneRepo: fortuneRepo,
	}
}

func (s *recordService) GetRecentRecords(userID uint, limit int) ([]models.FortuneRecord, error) {
	return s.recordRepo.FindByUserID(userID, limit)
}

func (s *recordService) GetRecordsByType(userID uint, recordType string, limit int) ([]models.FortuneRecord, error) {
	return s.recordRepo.FindByUserIDAndType(userID, recordType, limit)
}

func (s *recordService) GetSpouseImage(userID uint) (string, error) {
	fortuneInfo, err := s.fortuneRepo.FindByUserID(userID)
	if err != nil {
		return "", errors.New("fortune info not found")
	}

	if fortuneInfo.SpouseImageURL == "" {
		return "", errors.New("spouse image not found")
	}

	return fortuneInfo.SpouseImageURL, nil
}

