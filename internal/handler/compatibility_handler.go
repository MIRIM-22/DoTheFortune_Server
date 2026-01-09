package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"dothefortune_server/internal/service"
)

type CompatibilityHandler struct {
	compatibilityService service.CompatibilityService
}

func NewCompatibilityHandler(compatibilityService service.CompatibilityService) *CompatibilityHandler {
	return &CompatibilityHandler{
		compatibilityService: compatibilityService,
	}
}

// CalculateCompatibility godoc
// @Summary      궁합 계산
// @Description  현재 사용자와 다른 사용자 간의 궁합을 계산합니다. 두 사용자의 사주 정보를 비교하여 궁합 점수(0-100)와 분석 결과를 반환합니다. 계산 결과는 데이터베이스에 저장되며, 이후 조회 시 재계산 없이 저장된 결과를 반환합니다.
// @Tags         compatibility
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user2_id  query  int  true  "상대방 사용자 ID"  minimum(1)
// @Success      200       {object}  models.Compatibility  "궁합 계산 성공"
// @Failure      400       {object}  ErrorResponse  "잘못된 요청 (자기 자신과의 궁합 계산 시도 등)"
// @Failure      401       {object}  ErrorResponse  "인증 실패"
// @Failure      500       {object}  ErrorResponse  "서버 내부 오류 또는 사주 정보 없음"
// @Router       /compatibility/calculate [get]
func (h *CompatibilityHandler) CalculateCompatibility(c *gin.Context) {
	user1ID := c.MustGet("user_id").(uint)

	user2IDStr := c.Query("user2_id")
	if user2IDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user2_id is required"})
		return
	}

	user2ID, err := strconv.ParseUint(user2IDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user2_id"})
		return
	}

	compatibility, err := h.compatibilityService.CalculateCompatibility(user1ID, uint(user2ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, compatibility)
}

// GetCompatibility godoc
// @Summary      궁합 조회
// @Description  현재 사용자와 다른 사용자 간의 저장된 궁합 정보를 조회합니다. 저장된 궁합이 없으면 자동으로 계산하여 반환합니다.
// @Tags         compatibility
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user2_id  query  int  true  "상대방 사용자 ID"  minimum(1)
// @Success      200       {object}  models.Compatibility  "궁합 정보 조회 성공"
// @Failure      400       {object}  ErrorResponse  "잘못된 요청"
// @Failure      401       {object}  ErrorResponse  "인증 실패"
// @Failure      500       {object}  ErrorResponse  "서버 내부 오류"
// @Router       /compatibility/ [get]
func (h *CompatibilityHandler) GetCompatibility(c *gin.Context) {
	user1ID := c.MustGet("user_id").(uint)

	user2IDStr := c.Query("user2_id")
	if user2IDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user2_id is required"})
		return
	}

	user2ID, err := strconv.ParseUint(user2IDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user2_id"})
		return
	}

	compatibility, err := h.compatibilityService.GetCompatibility(user1ID, uint(user2ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, compatibility)
}

type CompatibilityMatchesResponse struct {
	Matches []interface{} `json:"matches" description:"궁합 정보 목록"`
}

// GetBestMatches godoc
// @Summary      최고 궁합 목록 조회
// @Description  현재 사용자와 가장 좋은 궁합을 가진 사용자들을 궁합 점수 순으로 반환합니다. 점수가 높은 순서대로 정렬되어 반환됩니다.
// @Tags         compatibility
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "반환할 최대 개수"  default(10)  minimum(1)  maximum(100)
// @Success      200    {object}  CompatibilityMatchesResponse  "최고 궁합 목록"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /compatibility/best [get]
func (h *CompatibilityHandler) GetBestMatches(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	matches, err := h.compatibilityService.GetBestMatches(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

// GetWorstMatches godoc
// @Summary      최악 궁합 목록 조회
// @Description  현재 사용자와 가장 나쁜 궁합을 가진 사용자들을 궁합 점수 순으로 반환합니다. 점수가 낮은 순서대로 정렬되어 반환됩니다.
// @Tags         compatibility
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "반환할 최대 개수"  default(10)  minimum(1)  maximum(100)
// @Success      200    {object}  CompatibilityMatchesResponse  "최악 궁합 목록"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /compatibility/worst [get]
func (h *CompatibilityHandler) GetWorstMatches(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	matches, err := h.compatibilityService.GetWorstMatches(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

