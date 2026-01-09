package repository

import (
	"dothefortune_server/internal/database"
	"dothefortune_server/internal/models"
)

type FortuneRepository interface {
	Create(fortune *models.FortuneInfo) error
	FindByUserID(userID uint) (*models.FortuneInfo, error)
	Update(fortune *models.FortuneInfo) error
	FindSimilarUsers(userID uint, limit int) ([]models.User, error)
}

type fortuneRepository struct{}

func NewFortuneRepository() FortuneRepository {
	return &fortuneRepository{}
}

func (r *fortuneRepository) Create(fortune *models.FortuneInfo) error {
	return database.DB.Create(fortune).Error
}

func (r *fortuneRepository) FindByUserID(userID uint) (*models.FortuneInfo, error) {
	var fortune models.FortuneInfo
	err := database.DB.Where("user_id = ?", userID).First(&fortune).Error
	if err != nil {
		return nil, err
	}
	return &fortune, nil
}

func (r *fortuneRepository) Update(fortune *models.FortuneInfo) error {
	return database.DB.Save(fortune).Error
}

func (r *fortuneRepository) FindSimilarUsers(userID uint, limit int) ([]models.User, error) {
	var user models.User
	if err := database.DB.Preload("FortuneInfo").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	if user.FortuneInfo == nil {
		return []models.User{}, nil
	}

	var users []models.User
	err := database.DB.
		Preload("FortuneInfo").
		Joins("INNER JOIN fortune_infos ON fortune_infos.user_id = users.id").
		Where("users.id != ?", userID).
		Where("fortune_infos.day_heavenly_stem = ? OR fortune_infos.day_earthly_branch = ?",
			user.FortuneInfo.DayHeavenlyStem, user.FortuneInfo.DayEarthlyBranch).
		Limit(limit).
		Find(&users).Error

	return users, err
}

