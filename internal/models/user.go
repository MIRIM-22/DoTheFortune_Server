package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email    string `gorm:"uniqueIndex;not null" json:"email" example:"user@example.com"`
	Name     string `gorm:"not null" json:"name" example:"홍길동"`
	Password string `gorm:"not null" json:"-"`
	Gender   string `gorm:"not null" json:"gender" example:"M" description:"성별 (M: 남성, F: 여성)"`

	FortuneInfo *FortuneInfo `gorm:"foreignKey:UserID" json:"fortune_info,omitempty"`
}

type FortuneInfo struct {
	ID        uint           `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint   `gorm:"uniqueIndex;not null" json:"user_id" example:"1"`
	BirthYear   int    `gorm:"not null" json:"birth_year" example:"2000"`
	BirthMonth  int    `gorm:"not null" json:"birth_month" example:"1"`
	BirthDay    int    `gorm:"not null" json:"birth_day" example:"1"`
	BirthHour   int    `json:"birth_hour" example:"12"`
	BirthMinute int    `json:"birth_minute" example:"0"`
	UnknownTime bool   `gorm:"default:false" json:"unknown_time" example:"false"`
	BirthPlace  string `gorm:"not null" json:"birth_place" example:"서울"`
	IsLunar     bool   `gorm:"default:false" json:"is_lunar" example:"false" description:"양력(false) 또는 음력(true)"`

	YearHeavenlyStem  string `json:"year_heavenly_stem" example:"庚"`
	YearEarthlyBranch string `json:"year_earthly_branch" example:"子"`
	MonthHeavenlyStem string `json:"month_heavenly_stem" example:"戊"`
	MonthEarthlyBranch string `json:"month_earthly_branch" example:"寅"`
	DayHeavenlyStem   string `json:"day_heavenly_stem" example:"甲"`
	DayEarthlyBranch  string `json:"day_earthly_branch" example:"子"`
	HourHeavenlyStem  string `json:"hour_heavenly_stem" example:"甲"`
	HourEarthlyBranch string `json:"hour_earthly_branch" example:"子"`
	
	SpouseImageURL string `json:"spouse_image_url" example:"https://example.com/spouse-image.jpg" description:"미리 생성된 배우자 이미지 URL"`
}

type FortuneRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID    uint   `gorm:"not null;index" json:"user_id"`
	Type      string `gorm:"not null" json:"type"`
	Content   string `gorm:"type:text" json:"content"`
	ImageURL  string `json:"image_url"`
	Metadata  string `gorm:"type:jsonb" json:"metadata"`
}

type Compatibility struct {
	ID        uint           `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User1ID    uint    `gorm:"not null;index" json:"user1_id" example:"1"`
	User2ID    uint    `gorm:"not null;index" json:"user2_id" example:"2"`
	Score      float64 `gorm:"not null" json:"score" example:"85.5" description:"궁합 점수 (0-100)"`
	Analysis   string  `gorm:"type:text" json:"analysis" example:"두 사람은 매우 좋은 궁합을 가지고 있습니다."`
	CompatibilityType string `gorm:"not null" json:"compatibility_type" example:"excellent" description:"궁합 타입 (excellent, good, normal, poor)"`
	
	// 카테고리별 분석 (기획서 기준)
	CommunicationAnalysis string `gorm:"type:text" json:"communication_analysis" example:"말하지 않아도 통하는 텔레파시가 있어요." description:"🗣️ 대화/가치관"`
	EmotionAnalysis       string `gorm:"type:text" json:"emotion_analysis" example:"서로의 부족한 점을 감싸주는 안정감을 느껴요." description:"💖 감정/성격"`
	LifestyleAnalysis     string `gorm:"type:text" json:"lifestyle_analysis" example:"함께 무언가를 도모하면 손발이 척척 맞아요." description:"🏠 목표/생활 방식"`
	CautionAnalysis       string `gorm:"type:text" json:"caution_analysis" example:"특별히 주의할 점은 없으나, 서로 예의를 지키는 게 중요해요." description:"⚡ 주의할 점"`
}

