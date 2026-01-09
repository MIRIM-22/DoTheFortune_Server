package service

import (
	"errors"
	"dothefortune_server/internal/models"
	"dothefortune_server/internal/repository"
	"dothefortune_server/internal/utils"
)

type SimilarUserResult struct {
	User  models.User
	Score float64
	Type  string // "similar", "best_match", "worst_match"
}

type TodayFortuneResult struct {
	TotalFortune    string   `json:"total_fortune"`     // 총운
	WealthFortune   string   `json:"wealth_fortune"`   // 재물운
	LoveFortune      string   `json:"love_fortune"`      // 애정운
	HealthFortune    string   `json:"health_fortune"`    // 건강운
	LuckyColor       string   `json:"lucky_color"`       // 행운의 컬러
	LuckyColorHex    string   `json:"lucky_color_hex"`   // 행운의 컬러 HEX
	LuckyNumbers     []int    `json:"lucky_numbers"`     // 행운의 숫자
}

type FortuneService interface {
	CreateOrUpdateFortuneInfo(userID uint, birthYear, birthMonth, birthDay, birthHour, birthMinute int, unknownTime bool, birthPlace string) (*models.FortuneInfo, error)
	GetFortuneInfo(userID uint) (*models.FortuneInfo, error)
	GetTodayFortune(userID uint) (*TodayFortuneResult, error)
	GetSimilarUsers(userID uint, limit int) ([]models.User, []float64, error)
	GetSimilarUserMatches(userID uint) (*SimilarUserResult, *SimilarUserResult, *SimilarUserResult, error) // 가장 비슷한, 잘 맞는, 잘 안 맞는
}

type fortuneService struct {
	fortuneRepo repository.FortuneRepository
	recordRepo  repository.RecordRepository
	aiService   AIService
}

func NewFortuneService(fortuneRepo repository.FortuneRepository, recordRepo repository.RecordRepository, aiService AIService) FortuneService {
	return &fortuneService{
		fortuneRepo: fortuneRepo,
		recordRepo:  recordRepo,
		aiService:   aiService,
	}
}

func (s *fortuneService) CreateOrUpdateFortuneInfo(userID uint, birthYear, birthMonth, birthDay, birthHour, birthMinute int, unknownTime bool, birthPlace string) (*models.FortuneInfo, error) {
	if unknownTime {
		birthHour = 12
		birthMinute = 0
	}

	yearStem, yearBranch, monthStem, monthBranch, dayStem, dayBranch, hourStem, hourBranch :=
		utils.CalculateFortunePillars(birthYear, birthMonth, birthDay, birthHour)

	existing, err := s.fortuneRepo.FindByUserID(userID)
	if err == nil && existing != nil {
		existing.BirthYear = birthYear
		existing.BirthMonth = birthMonth
		existing.BirthDay = birthDay
		existing.BirthHour = birthHour
		existing.BirthMinute = birthMinute
		existing.UnknownTime = unknownTime
		existing.BirthPlace = birthPlace
		existing.YearHeavenlyStem = yearStem
		existing.YearEarthlyBranch = yearBranch
		existing.MonthHeavenlyStem = monthStem
		existing.MonthEarthlyBranch = monthBranch
		existing.DayHeavenlyStem = dayStem
		existing.DayEarthlyBranch = dayBranch
		existing.HourHeavenlyStem = hourStem
		existing.HourEarthlyBranch = hourBranch

		if err := s.fortuneRepo.Update(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	fortuneInfo := &models.FortuneInfo{
		UserID:            userID,
		BirthYear:         birthYear,
		BirthMonth:        birthMonth,
		BirthDay:          birthDay,
		BirthHour:         birthHour,
		BirthMinute:       birthMinute,
		UnknownTime:       unknownTime,
		BirthPlace:        birthPlace,
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

	return fortuneInfo, nil
}

func (s *fortuneService) GetFortuneInfo(userID uint) (*models.FortuneInfo, error) {
	return s.fortuneRepo.FindByUserID(userID)
}

func (s *fortuneService) GetTodayFortune(userID uint) (*TodayFortuneResult, error) {
	fortuneInfo, err := s.fortuneRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("fortune info not found")
	}

	fortuneMap := map[string]string{
		"year_stem":   fortuneInfo.YearHeavenlyStem,
		"year_branch": fortuneInfo.YearEarthlyBranch,
		"month_stem":  fortuneInfo.MonthHeavenlyStem,
		"month_branch": fortuneInfo.MonthEarthlyBranch,
		"day_stem":    fortuneInfo.DayHeavenlyStem,
		"day_branch":  fortuneInfo.DayEarthlyBranch,
		"hour_stem":   fortuneInfo.HourHeavenlyStem,
		"hour_branch": fortuneInfo.HourEarthlyBranch,
	}

	todayStem, todayBranch := utils.CalculateTodayPillar()
	
	totalFortune, _ := s.aiService.GenerateFortuneText(fortuneMap, todayStem, todayBranch, "총운")
	wealthFortune, _ := s.aiService.GenerateFortuneText(fortuneMap, todayStem, todayBranch, "재물운")
	loveFortune, _ := s.aiService.GenerateFortuneText(fortuneMap, todayStem, todayBranch, "애정운")
	healthFortune, _ := s.aiService.GenerateFortuneText(fortuneMap, todayStem, todayBranch, "건강운")
	
	if totalFortune == "" {
		totalFortune = utils.GetTodayFortune(fortuneMap)
	}
	if wealthFortune == "" {
		wealthFortune = "오늘은 재물운이 안정적인 하루입니다. 계획적인 소비가 도움이 될 거예요."
	}
	if loveFortune == "" {
		loveFortune = "인연운이 평범한 날입니다. 자연스러운 만남을 기대해보세요."
	}
	if healthFortune == "" {
		healthFortune = "건강운이 양호한 하루입니다. 무리하지 말고 적당한 휴식이 필요해요."
	}
	
	luckyElement := utils.CalculateLuckyElement(fortuneMap, todayStem, todayBranch)
	luckyColor, luckyColorHex := utils.GetLuckyColor(luckyElement)
	luckyNumbers := utils.GetLuckyNumbers(luckyElement)

	result := &TodayFortuneResult{
		TotalFortune:  totalFortune,
		WealthFortune: wealthFortune,
		LoveFortune:   loveFortune,
		HealthFortune: healthFortune,
		LuckyColor:    luckyColor,
		LuckyColorHex: luckyColorHex,
		LuckyNumbers:  luckyNumbers,
	}

	record := &models.FortuneRecord{
		UserID:  userID,
		Type:    "today_fortune",
		Content: fortune,
		Metadata: `{"lucky_color": "` + luckyColor + `", "lucky_numbers": ` + utils.IntSliceToJSON(luckyNumbers) + `}`,
	}

	s.recordRepo.Create(record)

	return result, nil
}

func (s *fortuneService) GetSimilarUsers(userID uint, limit int) ([]models.User, []float64, error) {
	users, err := s.fortuneRepo.FindSimilarUsers(userID, limit)
	if err != nil {
		return nil, nil, err
	}

	currentFortune, err := s.fortuneRepo.FindByUserID(userID)
	if err != nil {
		return nil, nil, errors.New("current user fortune info not found")
	}

	currentMap := map[string]string{
		"year_stem":   currentFortune.YearHeavenlyStem,
		"year_branch": currentFortune.YearEarthlyBranch,
		"month_stem":  currentFortune.MonthHeavenlyStem,
		"month_branch": currentFortune.MonthEarthlyBranch,
		"day_stem":    currentFortune.DayHeavenlyStem,
		"day_branch":  currentFortune.DayEarthlyBranch,
		"hour_stem":   currentFortune.HourHeavenlyStem,
		"hour_branch": currentFortune.HourEarthlyBranch,
	}

	scores := make([]float64, len(users))
	for i, user := range users {
		if user.FortuneInfo == nil {
			continue
		}
		userMap := map[string]string{
			"year_stem":   user.FortuneInfo.YearHeavenlyStem,
			"year_branch": user.FortuneInfo.YearEarthlyBranch,
			"month_stem":  user.FortuneInfo.MonthHeavenlyStem,
			"month_branch": user.FortuneInfo.MonthEarthlyBranch,
			"day_stem":    user.FortuneInfo.DayHeavenlyStem,
			"day_branch":  user.FortuneInfo.DayEarthlyBranch,
			"hour_stem":   user.FortuneInfo.HourHeavenlyStem,
			"hour_branch": user.FortuneInfo.HourEarthlyBranch,
		}
		scores[i] = utils.CalculateSimilarityScore(currentMap, userMap)
	}

	return users, scores, nil
}

func (s *fortuneService) GetSimilarUserMatches(userID uint) (*SimilarUserResult, *SimilarUserResult, *SimilarUserResult, error) {
	allUsers, err := s.fortuneRepo.FindSimilarUsers(userID, 100)
	if err != nil {
		return nil, nil, nil, err
	}

	currentFortune, err := s.fortuneRepo.FindByUserID(userID)
	if err != nil {
		return nil, nil, nil, errors.New("current user fortune info not found")
	}

	currentMap := map[string]string{
		"year_stem":   currentFortune.YearHeavenlyStem,
		"year_branch": currentFortune.YearEarthlyBranch,
		"month_stem":  currentFortune.MonthHeavenlyStem,
		"month_branch": currentFortune.MonthEarthlyBranch,
		"day_stem":    currentFortune.DayHeavenlyStem,
		"day_branch":  currentFortune.DayEarthlyBranch,
		"hour_stem":   currentFortune.HourHeavenlyStem,
		"hour_branch": currentFortune.HourEarthlyBranch,
	}

	var similarUser *SimilarUserResult
	var bestMatchUser *SimilarUserResult
	var worstMatchUser *SimilarUserResult

	maxSimilarity := -1.0
	maxCompatibility := -1.0
	minConflict := 101.0

	for _, user := range allUsers {
		if user.FortuneInfo == nil || user.ID == userID {
			continue
		}

		userMap := map[string]string{
			"year_stem":   user.FortuneInfo.YearHeavenlyStem,
			"year_branch": user.FortuneInfo.YearEarthlyBranch,
			"month_stem":  user.FortuneInfo.MonthHeavenlyStem,
			"month_branch": user.FortuneInfo.MonthEarthlyBranch,
			"day_stem":    user.FortuneInfo.DayHeavenlyStem,
			"day_branch":  user.FortuneInfo.DayEarthlyBranch,
			"hour_stem":   user.FortuneInfo.HourHeavenlyStem,
			"hour_branch": user.FortuneInfo.HourEarthlyBranch,
		}

		similarityScore := utils.CalculateSimilarityScore(currentMap, userMap)
		if similarityScore > maxSimilarity {
			maxSimilarity = similarityScore
			similarUser = &SimilarUserResult{
				User:  user,
				Score: similarityScore,
				Type:  "similar",
			}
		}

		compatibilityScore := utils.CalculateCompatibilityScore(currentMap, userMap)
		if compatibilityScore > maxCompatibility {
			maxCompatibility = compatibilityScore
			bestMatchUser = &SimilarUserResult{
				User:  user,
				Score: compatibilityScore,
				Type:  "best_match",
			}
		}

		conflictScore := utils.CalculateConflictScore(currentMap, userMap)
		if conflictScore < minConflict {
			minConflict = conflictScore
			worstMatchUser = &SimilarUserResult{
				User:  user,
				Score: conflictScore,
				Type:  "worst_match",
			}
		}
	}

	return similarUser, bestMatchUser, worstMatchUser, nil
}
