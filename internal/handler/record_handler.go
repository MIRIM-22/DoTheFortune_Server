package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"dothefortune_server/internal/service"
)

type RecordHandler struct {
	recordService service.RecordService
}

func NewRecordHandler(recordService service.RecordService) *RecordHandler {
	return &RecordHandler{
		recordService: recordService,
	}
}

type RecordsResponse struct {
	Records []interface{} `json:"records" description:"운세 기록 목록"`
}

// GetRecentRecords godoc
// @Summary      최근 기록 조회
// @Description  사용자의 최근 운세 기록을 조회합니다. 오늘의 운세, 궁합 결과, AI 이미지 등 모든 타입의 기록이 최신순으로 반환됩니다.
// @Tags         records
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query  int  false  "반환할 최대 기록 수"  default(20)  minimum(1)  maximum(100)
// @Success      200    {object}  RecordsResponse  "기록 목록 조회 성공"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /records/ [get]
func (h *RecordHandler) GetRecentRecords(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	records, err := h.recordService.GetRecentRecords(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

// GetRecordsByType godoc
// @Summary      타입별 기록 조회
// @Description  특정 타입의 운세 기록만 필터링하여 조회합니다. 타입은 'today_fortune', 'compatibility', 'ai_spouse' 등이 있습니다.
// @Tags         records
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        type   path  string  true   "기록 타입 (today_fortune, compatibility, ai_spouse 등)"  example:"today_fortune"
// @Param        limit  query  int     false  "반환할 최대 기록 수"  default(20)  minimum(1)  maximum(100)
// @Success      200    {object}  RecordsResponse  "타입별 기록 목록"
// @Failure      401    {object}  ErrorResponse  "인증 실패"
// @Failure      500    {object}  ErrorResponse  "서버 내부 오류"
// @Router       /records/{type} [get]
func (h *RecordHandler) GetRecordsByType(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	recordType := c.Param("type")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	records, err := h.recordService.GetRecordsByType(userID, recordType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records})
}

type GetSpouseImageResponse struct {
	ImageURL string `json:"image_url" example:"https://example.com/spouse-image.jpg"`
}

// GetSpouseImage godoc
// @Summary      배우자 이미지 조회
// @Description  사용자의 사주 정보에 저장된 미리 생성된 배우자 이미지 URL을 조회합니다.
// @Tags         records
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  GetSpouseImageResponse  "이미지 조회 성공"
// @Failure      400      {object}  ErrorResponse  "사주 정보가 등록되지 않음"
// @Failure      404      {object}  ErrorResponse  "배우자 이미지가 없음"
// @Failure      401      {object}  ErrorResponse  "인증 실패"
// @Failure      500      {object}  ErrorResponse  "서버 내부 오류"
// @Router       /records/spouse-image [get]
func (h *RecordHandler) GetSpouseImage(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	imageURL, err := h.recordService.GetSpouseImage(userID)
	if err != nil {
		if err.Error() == "spouse image not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, GetSpouseImageResponse{
		ImageURL: imageURL,
	})
}

