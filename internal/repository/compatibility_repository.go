package repository

import (
	"dothefortune_server/internal/database"
	"dothefortune_server/internal/models"
)

type CompatibilityRepository interface {
	Create(compatibility *models.Compatibility) error
	FindByUserPair(user1ID, user2ID uint) (*models.Compatibility, error)
	FindBestMatches(userID uint, limit int) ([]models.Compatibility, error)
	FindWorstMatches(userID uint, limit int) ([]models.Compatibility, error)
}

type compatibilityRepository struct{}

func NewCompatibilityRepository() CompatibilityRepository {
	return &compatibilityRepository{}
}

func (r *compatibilityRepository) Create(compatibility *models.Compatibility) error {
	return database.DB.Create(compatibility).Error
}

func (r *compatibilityRepository) FindByUserPair(user1ID, user2ID uint) (*models.Compatibility, error) {
	var compatibility models.Compatibility
	err := database.DB.
		Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		First(&compatibility).Error
	if err != nil {
		return nil, err
	}
	return &compatibility, nil
}

func (r *compatibilityRepository) FindBestMatches(userID uint, limit int) ([]models.Compatibility, error) {
	var compatibilities []models.Compatibility
	err := database.DB.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("score DESC").
		Limit(limit).
		Find(&compatibilities).Error
	return compatibilities, err
}

func (r *compatibilityRepository) FindWorstMatches(userID uint, limit int) ([]models.Compatibility, error) {
	var compatibilities []models.Compatibility
	err := database.DB.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("score ASC").
		Limit(limit).
		Find(&compatibilities).Error
	return compatibilities, err
}

