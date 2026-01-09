package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"dothefortune_server/internal/service"
)

type FortuneHandler struct {
	fortuneService service.FortuneService
}

func NewFortuneHandler(fortuneService service.FortuneService) *FortuneHandler {
	return &FortuneHandler{
		fortuneService: fortuneService,
	}
}

type CreateFortuneInfoRequest struct {
	BirthYear   int    `json:"birth_year" binding:"required" example:"2000" swaggertype:"integer" minimum:"1900" maximum:"2100" description:"출생 연도"`
	BirthMonth  int    `json:"birth_month" binding:"required" example:"1" swaggertype:"integer" minimum:"1" maximum:"12" description:"출생 월"`
	BirthDay    int    `json:"birth_day" binding:"required" example:"1" swaggertype:"integer" minimum:"1" maximum:"31" description:"출생 일"`
	BirthHour   int    `json:"birth_hour" example:"12" swaggertype:"integer" minimum:"0" maximum:"23" description:"출생 시각 (0-23), unknown_time이 true면 무시됨"`
	BirthMinute int    `json:"birth_minute" example:"0" swaggertype:"integer" minimum:"0" maximum:"59" description:"출생 분 (0-59), unknown_time이 true면 무시됨"`
	UnknownTime bool   `json:"unknown_time" example:"false" swaggertype:"boolean" description:"출생 시각을 모를 경우 true, 기본값은 12시 0분으로 설정됨"`
	BirthPlace  string `json:"birth_place" binding:"required" example:"서울" swaggertype:"string" description:"출생지"`
}

type TodayFortuneResponse struct {
	TotalFortune    string   `json:"total_fortune" example:"오늘은 새로운 시작에 좋은 날입니다."`
	WealthFortune   string   `json:"wealth_fortune" example:"재물운이 안정적인 하루입니다."`
	LoveFortune     string   `json:"love_fortune" example:"인연운이 평범한 날입니다."`
	HealthFortune   string   `json:"health_fortune" example:"건강운이 양호한 하루입니다."`
	LuckyColor      string   `json:"lucky_color" example:"초록"`
	LuckyColorHex   string   `json:"lucky_color_hex" example:"#4CAF50"`
	LuckyNumbers    []int    `json:"lucky_numbers" example:"[3,8]"`
}

type SimilarUsersResponse struct {
	Users []UserWithScore `json:"users"`
}

type UserWithScore struct {
	User  interface{} `json:"user"`
	Score float64     `json:"similarity_score" example:"85.5" description:"유사도 점수 (0-100)"`
}

// CreateOrUpdateFortuneInfo godoc
// @Summary      사주 정보 등록/수정
// @Description  사용자의 사주 정보를 등록하거나 수정합니다. 생년월일, 출생 시각, 출생지를 입력받아 사주를 계산하고 저장합니다. 출생 시각을 모를 경우 unknown_time을 true로 설정하면 자동으로 12시 0분으로 설정됩니다.
// @Tags         fortune
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  CreateFortuneInfoRequest  true  "사주 정보"
// @Success      200      {object}  models.FortuneInfo  "사주 정보 저장 성공"
// @Failure      400      {object}  ErrorResponse  "잘못된 요청 (필수 필드 누락, 잘못된 날짜 형식 등)"
// @Failure      401      {object}  ErrorResponse  "인증 실패"
// @Failure      500      {object}  ErrorResponse  "서버 내부 오류"
// @Router       /fortune/info [post]
func (h *FortuneHandler) CreateOrUpdateFortuneInfo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req CreateFortuneInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fortuneInfo, err := h.fortuneService.CreateOrUpdateFortuneInfo(
		userID,
		req.BirthYear,
		req.BirthMonth,
		req.BirthDay,
		req.BirthHour,
		req.BirthMinute,
		req.UnknownTime,
		req.BirthPlace,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fortuneInfo)
}

// GetFortuneInfo godoc
// @Summary      사주 정보 조회
// @Description  현재 사용자의 등록된 사주 정보를 조회합니다. 사주 정보가 없으면 404 오류가 반환됩니다.
// @Tags         fortune
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.FortuneInfo  "사주 정보 조회 성공"
// @Failure      401  {object}  ErrorResponse  "인증 실패"
// @Failure      404  {object}  ErrorResponse  "사주 정보를 찾을 수 없음"
// @Router       /fortune/info [get]
func (h *FortuneHandler) GetFortuneInfo(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	fortuneInfo, err := h.fortuneService.GetFortuneInfo(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fortune info not found"})
		return
	}

	c.JSON(http.StatusOK, fortuneInfo)
}

// GetTodayFortune godoc
// @Summary      오늘의 운세 조회
// @Description  사용자의 사주 정보를 기반으로 오늘의 운세를 제공합니다. 총운, 재물운, 애정운, 건강운과 행운의 컬러, 행운의 숫자를 포함합니다. 사주 정보가 등록되어 있어야 하며, 조회 시 기록이 자동으로 저장됩니다.
// @Tags         fortune
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  TodayFortuneResponse  "오늘의 운세 조회 성공"
// @Failure      400  {object}  ErrorResponse  "사주 정보가 등록되지 않음"
// @Failure      401  {object}  ErrorResponse  "인증 실패"
// @Router       /fortune/today [get]
func (h *FortuneHandler) GetTodayFortune(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	fortune, err := h.fortuneService.GetTodayFortune(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fortune)
}

// GetSimilarUsers godoc
// @Summary      유사 사주 친구 찾기
// @Description  현재 사용자와 유사한 사주를 가진 다른 사용자들을 찾아 반환합니다. 일간 천간 또는 지지를 기준으로 유사도를 계산하며, 유사도 점수(0-100)와 함께 반환됩니다.
// @Tags         fortune
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "반환할 최대 사용자 수"  default(10)  minimum(1)  maximum(100)
// @Success      200    {object}  SimilarUsersResponse  "유사 사용자 목록"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /fortune/similar [get]
func (h *FortuneHandler) GetSimilarUsers(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	users, scores, err := h.fortuneService.GetSimilarUsers(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]UserWithScore, len(users))
	for i, user := range users {
		user.FortuneInfo = nil
		result[i] = UserWithScore{
			User:  user,
			Score: scores[i],
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": result})
}

type SimilarUserMatchesResponse struct {
	SimilarUser  *service.SimilarUserResult `json:"similar_user" description:"가장 비슷한 사주 친구"`
	BestMatch    *service.SimilarUserResult `json:"best_match" description:"잘 맞는 사주 친구"`
	WorstMatch   *service.SimilarUserResult `json:"worst_match" description:"잘 안 맞는 사주 친구"`
}

// GetSimilarUserMatches godoc
// @Summary      유사 사주 친구 매칭 (가장 비슷한, 잘 맞는, 잘 안 맞는)
// @Description  현재 사용자와 가장 비슷한 사주, 잘 맞는 사주, 잘 안 맞는 사주를 가진 친구를 각각 1명씩 찾아 반환합니다. 기획서 기준의 상세한 매칭 로직을 사용합니다.
// @Tags         fortune
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200    {object}  SimilarUserMatchesResponse  "유사 사주 친구 매칭 결과"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /fortune/similar-matches [get]
func (h *FortuneHandler) GetSimilarUserMatches(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	similar, best, worst, err := h.fortuneService.GetSimilarUserMatches(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SimilarUserMatchesResponse{
		SimilarUser: similar,
		BestMatch:   best,
		WorstMatch:  worst,
	})
}

