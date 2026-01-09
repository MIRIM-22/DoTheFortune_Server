package service

import (
	"errors"
	"fmt"
	"dothefortune_server/internal/models"
	"dothefortune_server/internal/repository"
	"dothefortune_server/internal/utils"
)

type CompatibilityService interface {
	CalculateCompatibility(user1ID, user2ID uint) (*models.Compatibility, error)
	GetCompatibility(user1ID, user2ID uint) (*models.Compatibility, error)
	GetBestMatches(userID uint, limit int) ([]models.Compatibility, error)
	GetWorstMatches(userID uint, limit int) ([]models.Compatibility, error)
}

type compatibilityService struct {
	compatibilityRepo repository.CompatibilityRepository
	fortuneRepo       repository.FortuneRepository
	recordRepo        repository.RecordRepository
}

func NewCompatibilityService(compatibilityRepo repository.CompatibilityRepository, fortuneRepo repository.FortuneRepository, recordRepo repository.RecordRepository) CompatibilityService {
	return &compatibilityService{
		compatibilityRepo: compatibilityRepo,
		fortuneRepo:       fortuneRepo,
		recordRepo:        recordRepo,
	}
}

func (s *compatibilityService) CalculateCompatibility(user1ID, user2ID uint) (*models.Compatibility, error) {
	if user1ID == user2ID {
		return nil, errors.New("cannot calculate compatibility with yourself")
	}

	existing, err := s.compatibilityRepo.FindByUserPair(user1ID, user2ID)
	if err == nil && existing != nil {
		return existing, nil
	}

	fortune1, err := s.fortuneRepo.FindByUserID(user1ID)
	if err != nil {
		return nil, errors.New("user1 fortune info not found")
	}

	fortune2, err := s.fortuneRepo.FindByUserID(user2ID)
	if err != nil {
		return nil, errors.New("user2 fortune info not found")
	}

	fortune1Map := map[string]string{
		"year_stem":   fortune1.YearHeavenlyStem,
		"year_branch": fortune1.YearEarthlyBranch,
		"month_stem":  fortune1.MonthHeavenlyStem,
		"month_branch": fortune1.MonthEarthlyBranch,
		"day_stem":    fortune1.DayHeavenlyStem,
		"day_branch":  fortune1.DayEarthlyBranch,
		"hour_stem":   fortune1.HourHeavenlyStem,
		"hour_branch": fortune1.HourEarthlyBranch,
	}

	fortune2Map := map[string]string{
		"year_stem":   fortune2.YearHeavenlyStem,
		"year_branch": fortune2.YearEarthlyBranch,
		"month_stem":  fortune2.MonthHeavenlyStem,
		"month_branch": fortune2.MonthEarthlyBranch,
		"day_stem":    fortune2.DayHeavenlyStem,
		"day_branch":  fortune2.DayEarthlyBranch,
		"hour_stem":   fortune2.HourHeavenlyStem,
		"hour_branch": fortune2.HourEarthlyBranch,
	}

	score := utils.CalculateCompatibilityScore(fortune1Map, fortune2Map)

	compatibilityType := "normal"
	if score >= 80 {
		compatibilityType = "excellent"
	} else if score >= 60 {
		compatibilityType = "good"
	} else if score < 40 {
		compatibilityType = "poor"
	}

	analysis := generateCompatibilityAnalysis(score, compatibilityType)
	
	commAnalysis, emotionAnalysis, lifestyleAnalysis, cautionAnalysis := 
		generateCategoryAnalysis(fortune1Map, fortune2Map, score)

	compatibility := &models.Compatibility{
		User1ID:              user1ID,
		User2ID:              user2ID,
		Score:                score,
		Analysis:             analysis,
		CompatibilityType:    compatibilityType,
		CommunicationAnalysis: commAnalysis,
		EmotionAnalysis:      emotionAnalysis,
		LifestyleAnalysis:    lifestyleAnalysis,
		CautionAnalysis:      cautionAnalysis,
	}

	if err := s.compatibilityRepo.Create(compatibility); err != nil {
		return nil, err
	}

	record := &models.FortuneRecord{
		UserID:  user1ID,
		Type:    "compatibility",
		Content: fmt.Sprintf("Compatibility with user %d: %.1f%%", user2ID, score),
		Metadata: fmt.Sprintf(`{"user2_id": %d, "score": %.1f, "type": "%s"}`, user2ID, score, compatibilityType),
	}

	s.recordRepo.Create(record)

	return compatibility, nil
}

func (s *compatibilityService) GetCompatibility(user1ID, user2ID uint) (*models.Compatibility, error) {
	compatibility, err := s.compatibilityRepo.FindByUserPair(user1ID, user2ID)
	if err != nil {
		return s.CalculateCompatibility(user1ID, user2ID)
	}
	return compatibility, nil
}

func (s *compatibilityService) GetBestMatches(userID uint, limit int) ([]models.Compatibility, error) {
	return s.compatibilityRepo.FindBestMatches(userID, limit)
}

func (s *compatibilityService) GetWorstMatches(userID uint, limit int) ([]models.Compatibility, error) {
	return s.compatibilityRepo.FindWorstMatches(userID, limit)
}

func generateCompatibilityAnalysis(score float64, compatibilityType string) string {
	switch compatibilityType {
	case "excellent":
		return "두 사람은 매우 좋은 궁합을 가지고 있습니다. 서로를 잘 이해하고 보완하는 관계가 될 것입니다."
	case "good":
		return "두 사람은 좋은 궁합을 가지고 있습니다. 서로 협력하며 발전할 수 있는 관계입니다."
	case "poor":
		return "두 사람은 서로 다른 성향을 가지고 있어 이해가 필요합니다. 인내심과 소통이 중요합니다."
	default:
		return "두 사람은 평범한 궁합을 가지고 있습니다. 서로의 차이를 존중하며 관계를 발전시켜 나가세요."
	}
}

func generateCategoryAnalysis(fortune1, fortune2 map[string]string, score float64) (commAnalysis, emotionAnalysis, lifestyleAnalysis, cautionAnalysis string) {
	dayStem1 := fortune1["day_stem"]
	dayStem2 := fortune2["day_stem"]
	dayBranch1 := fortune1["day_branch"]
	dayBranch2 := fortune2["day_branch"]

	commAnalysis = analyzeCommunication(dayStem1, dayStem2)
	emotionAnalysis = analyzeEmotion(fortune1, fortune2)
	lifestyleAnalysis = analyzeLifestyle(dayBranch1, dayBranch2)
	cautionAnalysis = analyzeCaution(dayBranch1, dayBranch2)

	return commAnalysis, emotionAnalysis, lifestyleAnalysis, cautionAnalysis
}

func analyzeCommunication(stem1, stem2 string) string {
	if utils.IsHeavenlyStemPair(stem1, stem2) {
		return "말하지 않아도 통하는 텔레파시가 있어요."
	}
	if utils.IsHeavenlyStemClash(stem1, stem2) {
		return "가치관이 달라 논쟁이 될 수 있지만, 새로운 시각을 줘요."
	}
	element1 := utils.GetElement(stem1)
	element2 := utils.GetElement(stem2)
	if element1 == element2 && element1 != "" {
		return "친구처럼 편안하게 대화가 흘러가요."
	}
	return "서로 다른 관점을 나누며 대화가 이어져요."
}

func analyzeEmotion(fortune1, fortune2 map[string]string) string {
	user1Elements := utils.GetFiveElements(fortune1)
	user2Elements := utils.GetFiveElements(fortune2)
	
	complementCount := utils.CountComplementaryElements(user1Elements, user2Elements)
	if complementCount >= 2 {
		return "서로의 부족한 점을 감싸주는 안정감을 느껴요."
	}
	
	if utils.HasElementBias(user1Elements, user2Elements) {
		return "성격이 너무 비슷해서 오히려 부딪힐 때가 있어요."
	}
	
	return "서로의 감정을 잘 이해하고 공감할 수 있어요."
}

func analyzeLifestyle(branch1, branch2 string) string {
	if utils.IsEarthlyBranchSixPair(branch1, branch2) {
		return "함께 무언가를 도모하면 손발이 척척 맞아요."
	}
	if utils.IsEarthlyBranchThreePair(branch1, branch2) {
		return "목표와 가치관이 잘 맞아 협력이 잘 돼요."
	}
	if utils.IsEarthlyBranchClash(branch1, branch2) {
		return "활동 반경이나 생활 패턴이 달라서 조율이 필요해요."
	}
	return "서로의 생활 방식을 존중하며 조화롭게 지낼 수 있어요."
}

func analyzeCaution(branch1, branch2 string) string {
	if utils.IsEarthlyBranchResentment(branch1, branch2) {
		return "사소한 오해가 감정 싸움으로 번지지 않게 배려가 필요해요."
	}
	if utils.IsEarthlyBranchClash(branch1, branch2) {
		return "의견 차이가 있을 때 바로 해결하지 않으면 오래가요."
	}
	return "특별히 주의할 점은 없으나, 서로 예의를 지키는 게 중요해요."
}

