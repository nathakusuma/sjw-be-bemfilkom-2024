package main

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest"
	"github.com/bem-filkom/sjw-be-2024/internal/app/middleware"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/database/postgresql"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	db := postgresql.Connect()
	jwtAuth := jwt.NewJWT(os.Getenv("JWT_SECRET"), os.Getenv("JWT_TTL"))
	middle := middleware.NewMiddleware(jwtAuth)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtAuth)
	authHandler := rest.NewAuthHandler(userService)

	hopeWhisperRepo := repository.NewHopeWhisperRepository(db)
	hopeWhisperService := service.NewHopeWhisperService(hopeWhisperRepo)
	hopeWhisperHandler := rest.NewHopeWhisperHandler(hopeWhisperService)

	gin.SetMode(os.Getenv("GIN_MODE"))

	router := gin.Default()

	router.Use(middleware.CORS)

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/docs/v1/")
	})

	docs := router.Group("/docs")
	docs.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/docs/v1/")
	})
	docs.GET("/v1.yaml", func(ctx *gin.Context) {
		ctx.File("./docs/api/v1.yaml")
	})
	docs.Static("/v1", "./web/swagger-ui")

	v1 := router.Group("/v1")
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.GET("/check/admin", middle.Authenticate, middle.RequireRole("admin"), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "you are admin"})
	})

	admin := v1.Group("/admin")
	admin.Use(middle.Authenticate, middle.RequireRole("admin"))
	// Admin Hope Corner
	admin.GET("/hopes/:id", hopeWhisperHandler.FindByID(model.HopeCorner, true))
	admin.GET("/hopes", hopeWhisperHandler.FindByLazyLoad(model.HopeCorner, true))
	admin.PATCH("/hopes/:id", hopeWhisperHandler.Update(model.HopeCorner))
	admin.DELETE("/hopes/:id", hopeWhisperHandler.Delete(model.HopeCorner))
	// Admin Whisper Wall
	admin.GET("/whispers/:id", hopeWhisperHandler.FindByID(model.WhisperWall, true))
	admin.GET("/whispers", hopeWhisperHandler.FindByLazyLoad(model.WhisperWall, true))
	admin.PATCH("/whispers/:id", hopeWhisperHandler.Update(model.WhisperWall))
	admin.DELETE("/whispers/:id", hopeWhisperHandler.Delete(model.WhisperWall))

	hopes := v1.Group("/hopes")
	hopes.POST("/", hopeWhisperHandler.Create(model.HopeCorner))
	hopes.GET("/:id", hopeWhisperHandler.FindByID(model.HopeCorner, false))
	hopes.GET("/", hopeWhisperHandler.FindByLazyLoad(model.HopeCorner, false))

	whispers := v1.Group("/whispers")
	whispers.POST("/", hopeWhisperHandler.Create(model.WhisperWall))
	whispers.GET("/:id", hopeWhisperHandler.FindByID(model.WhisperWall, false))
	whispers.GET("/", hopeWhisperHandler.FindByLazyLoad(model.WhisperWall, false))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalln(err)
	}
}
