package main

import (
	"log"

	"dothefortune_server/internal/config"
	"dothefortune_server/internal/database"
	"dothefortune_server/internal/router"
	"dothefortune_server/internal/utils"
)

// @title           DoTheFortune API
// @version         1.0
// @description     2025 파이썬 반 대항 프로젝트 빌려온 사주 서버

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT 토큰을 사용한 인증. "Bearer " 뒤에 공백을 두고 JWT 토큰을 입력하세요. 예: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	cfg := config.Load()

	utils.InitJWT(cfg.JWTSecret)

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := router.SetupRouter(cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

