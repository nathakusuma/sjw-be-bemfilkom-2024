package main

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest"
	"github.com/bem-filkom/sjw-be-2024/internal/app/middleware"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/database/postgresql"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
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

	hopeRepo := repository.NewHopeCornerRepository(db)
	hopeService := service.NewHopeCornerService(hopeRepo)
	hopeHandler := rest.NewHopeCornerHandler(hopeService)

	gin.SetMode(os.Getenv("GIN_MODE"))

	router := gin.Default()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	router.Use(func(ctx *gin.Context) {
		corsMiddleware.HandlerFunc(ctx.Writer, ctx.Request)
	})

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
	admin.GET("/hopes/:id", hopeHandler.FindByID(true))
	admin.GET("/hopes", hopeHandler.FindByLazyLoad(true))
	admin.PATCH("/hopes/:id", hopeHandler.Update)
	admin.DELETE("/hopes/:id", hopeHandler.Delete)

	hopes := v1.Group("/hopes")
	hopes.POST("/", hopeHandler.Create)
	hopes.GET("/:id", hopeHandler.FindByID(false))
	hopes.GET("/", hopeHandler.FindByLazyLoad(false))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalln(err)
	}
}
