package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"dothefortune_server/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email" example:"user@example.com" swaggertype:"string" format:"email"`
	Password   string `json:"password" binding:"required,min=6" example:"password123" swaggertype:"string" minLength:"6"`
	Name       string `json:"name" binding:"required" example:"홍길동" swaggertype:"string"`
	Gender     string `json:"gender" binding:"required,oneof=M F" example:"M" swaggertype:"string" description:"성별 (M: 남성, F: 여성)"`
	BirthYear  int    `json:"birth_year" binding:"required" example:"2000" swaggertype:"integer" minimum:"1900" maximum:"2100"`
	BirthMonth int    `json:"birth_month" binding:"required" example:"1" swaggertype:"integer" minimum:"1" maximum:"12"`
	BirthDay   int    `json:"birth_day" binding:"required" example:"1" swaggertype:"integer" minimum:"1" maximum:"31"`
	BirthHour  int    `json:"birth_hour" example:"12" swaggertype:"integer" minimum:"0" maximum:"23" description:"태어난 시간 (0-23), 모를 경우 생략 가능"`
	BirthMinute int    `json:"birth_minute" example:"0" swaggertype:"integer" minimum:"0" maximum:"59" description:"태어난 분 (0-59), 모를 경우 생략 가능"`
	IsLunar    bool   `json:"is_lunar" example:"false" swaggertype:"boolean" description:"양력(false) 또는 음력(true)"`
	BirthPlace string `json:"birth_place" binding:"required" example:"서울" swaggertype:"string" description:"태어난 도시명"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com" swaggertype:"string" format:"email"`
	Password string `json:"password" binding:"required" example:"password123" swaggertype:"string"`
}

type AuthResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  interface{} `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type MessageResponse struct {
	Message string `json:"message" example:"success message"`
}

// Register godoc
// @Summary      회원가입
// @Description  새로운 사용자를 등록합니다. 이메일, 비밀번호(최소 6자), 이름, 성별, 생년월일(양력/음력), 태어난 시간, 태어난 도시명을 입력받아 계정과 사주 정보를 생성하고 JWT 토큰을 발급합니다. 토큰은 쿠키에도 저장됩니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  RegisterRequest  true  "회원가입 요청 정보"
// @Success      201      {object}  AuthResponse  "회원가입 성공"
// @Failure      400      {object}  ErrorResponse  "잘못된 요청 (이메일 형식 오류, 비밀번호 길이 부족, 이미 존재하는 이메일 등)"
// @Failure      500      {object}  ErrorResponse  "서버 내부 오류"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(
		req.Email,
		req.Password,
		req.Name,
		req.Gender,
		req.BirthYear,
		req.BirthMonth,
		req.BirthDay,
		req.BirthHour,
		req.BirthMinute,
		req.IsLunar,
		req.BirthPlace,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, int(24*time.Hour.Seconds()), "/", "", false, true)

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login godoc
// @Summary      로그인
// @Description  이메일과 비밀번호로 로그인합니다. 인증 성공 시 JWT 토큰을 발급하며, 토큰은 응답 본문과 쿠키에 포함됩니다. 토큰은 24시간 동안 유효합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  LoginRequest  true  "로그인 요청 정보"
// @Success      200      {object}  AuthResponse  "로그인 성공"
// @Failure      400      {object}  ErrorResponse  "잘못된 요청 (이메일 형식 오류 등)"
// @Failure      401      {object}  ErrorResponse  "인증 실패 (이메일 또는 비밀번호 불일치)"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, int(24*time.Hour.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout godoc
// @Summary      로그아웃
// @Description  로그아웃하고 인증 쿠키를 삭제합니다. 클라이언트 측에서도 토큰을 제거해야 합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  MessageResponse  "로그아웃 성공"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetMe godoc
// @Summary      현재 사용자 정보 조회
// @Description  현재 인증된 사용자의 ID를 반환합니다. Authorization 헤더의 Bearer 토큰 또는 쿠키의 토큰을 사용하여 인증합니다.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "사용자 ID 반환"  example:"{\"user_id\": 1}"
// @Failure      401  {object}  ErrorResponse  "인증 실패 (토큰 없음, 만료, 또는 유효하지 않음)"
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

