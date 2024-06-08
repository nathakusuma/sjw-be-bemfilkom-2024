package config

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest"
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest/middleware"
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest/route"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

type StartAppConfig struct {
	DB  *gorm.DB
	App *gin.Engine
}

func StartApp(config *StartAppConfig) {
	jwtAuth := jwt.NewJWT(os.Getenv("JWT_SECRET"), os.Getenv("JWT_TTL"))

	// Repository
	userRepo := repository.NewUserRepository(config.DB)
	hopeWhisperRepo := repository.NewHopeWhisperRepository(config.DB)

	// Service
	userService := service.NewUserService(userRepo, jwtAuth)
	hopeWhisperService := service.NewHopeWhisperService(hopeWhisperRepo)

	// Middleware
	middle := middleware.NewMiddleware(jwtAuth)

	// Handler
	authHandler := rest.NewAuthHandler(userService)
	hopeWhisperHandler := rest.NewHopeWhisperHandler(hopeWhisperService)

	routeConfig := route.Config{
		App:                config.App,
		AuthHandler:        authHandler,
		HopeWhisperHandler: hopeWhisperHandler,
		Middleware:         middle,
	}
	routeConfig.Setup()
}
