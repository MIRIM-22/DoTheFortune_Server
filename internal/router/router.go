package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "dothefortune_server/docs"
	"dothefortune_server/internal/config"
	"dothefortune_server/internal/handler"
	"dothefortune_server/internal/middleware"
	"dothefortune_server/internal/repository"
	"dothefortune_server/internal/service"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))

	userRepo := repository.NewUserRepository()
	fortuneRepo := repository.NewFortuneRepository()
	recordRepo := repository.NewRecordRepository()
	compatibilityRepo := repository.NewCompatibilityRepository()

	authService := service.NewAuthService(userRepo, fortuneRepo)
	aiService := service.NewAIService(fortuneRepo, userRepo, cfg)
	fortuneService := service.NewFortuneService(fortuneRepo, recordRepo, aiService)
	compatibilityService := service.NewCompatibilityService(compatibilityRepo, fortuneRepo, recordRepo)
	recordService := service.NewRecordService(recordRepo, fortuneRepo)

	authHandler := handler.NewAuthHandler(authService)
	fortuneHandler := handler.NewFortuneHandler(fortuneService)
	compatibilityHandler := handler.NewCompatibilityHandler(compatibilityService)
	recordHandler := handler.NewRecordHandler(recordService)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetMe)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			fortune := protected.Group("/fortune")
			{
				fortune.POST("/info", fortuneHandler.CreateOrUpdateFortuneInfo)
				fortune.GET("/info", fortuneHandler.GetFortuneInfo)
				fortune.GET("/today", fortuneHandler.GetTodayFortune)
				fortune.GET("/similar", fortuneHandler.GetSimilarUsers)
				fortune.GET("/similar-matches", fortuneHandler.GetSimilarUserMatches)
			}

			compatibility := protected.Group("/compatibility")
			{
				compatibility.GET("/calculate", compatibilityHandler.CalculateCompatibility)
				compatibility.GET("/", compatibilityHandler.GetCompatibility)
				compatibility.GET("/best", compatibilityHandler.GetBestMatches)
				compatibility.GET("/worst", compatibilityHandler.GetWorstMatches)
			}

			records := protected.Group("/records")
			{
				records.GET("/", recordHandler.GetRecentRecords)
				records.GET("/:type", recordHandler.GetRecordsByType)
				records.GET("/spouse-image", recordHandler.GetSpouseImage)
			}
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

