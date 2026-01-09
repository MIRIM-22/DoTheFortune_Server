package service

import (
	"errors"
	"dothefortune_server/internal/models"
	"dothefortune_server/internal/repository"
	"dothefortune_server/internal/utils"
)

type AuthService interface {
	Register(email, password, name, gender string, birthYear, birthMonth, birthDay, birthHour, birthMinute int, isLunar bool, birthPlace string) (*models.User, error)
	Login(email, password string) (*models.User, string, error)
	GenerateToken(userID uint, email string) (string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	fortuneRepo repository.FortuneRepository
}

func NewAuthService(userRepo repository.UserRepository, fortuneRepo repository.FortuneRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		fortuneRepo: fortuneRepo,
	}
}

func (s *authService) Register(email, password, name, gender string, birthYear, birthMonth, birthDay, birthHour, birthMinute int, isLunar bool, birthPlace string) (*models.User, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	unknownTime := birthHour == 0 && birthMinute == 0
	if unknownTime {
		birthHour = 12
		birthMinute = 0
	}

	newUser := &models.User{
		Email:    email,
		Name:     name,
		Password: hashedPassword,
		Gender:   gender,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	yearStem, yearBranch, monthStem, monthBranch, dayStem, dayBranch, hourStem, hourBranch :=
		utils.CalculateFortunePillars(birthYear, birthMonth, birthDay, birthHour)

	fortuneInfo := &models.FortuneInfo{
		UserID:            newUser.ID,
		BirthYear:         birthYear,
		BirthMonth:        birthMonth,
		BirthDay:          birthDay,
		BirthHour:         birthHour,
		BirthMinute:       birthMinute,
		UnknownTime:       unknownTime,
		BirthPlace:        birthPlace,
		IsLunar:           isLunar,
		YearHeavenlyStem:  yearStem,
		YearEarthlyBranch: yearBranch,
		MonthHeavenlyStem: monthStem,
		MonthEarthlyBranch: monthBranch,
		DayHeavenlyStem:   dayStem,
		DayEarthlyBranch:  dayBranch,
		HourHeavenlyStem:  hourStem,
		HourEarthlyBranch: hourBranch,
	}

	if err := s.fortuneRepo.Create(fortuneInfo); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := s.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) GenerateToken(userID uint, email string) (string, error) {
	return utils.GenerateToken(userID, email)
}

