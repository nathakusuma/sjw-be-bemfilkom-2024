package route

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/handler/rest"
	"github.com/bem-filkom/sjw-be-2024/internal/app/middleware"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type Config struct {
	App                *gin.Engine
	AuthHandler        rest.IAuthHandler
	HopeWhisperHandler rest.IHopeWhisperHandler
	Middleware         middleware.IMiddleware
}

func (c *Config) Setup() {
	c.App.Use(gin.Logger())
	c.App.Use(gin.Recovery())
	c.App.Use(c.Middleware.CORS())

	c.docsRoute()

	v1 := c.App.Group("/v1")
	c.authRoute(v1)
	c.hopeCornerRoute(v1)
	c.whisperWallRoute(v1)
	c.adminRoute(v1)
}

func (c *Config) docsRoute() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/docs/v1/")
	})

	docs := c.App.Group("/docs")
	docs.GET("", func(ctx *gin.Context) {
		ctx.Redirect(301, "/docs/v1/")
	})
	docs.GET("/v1.yaml", func(ctx *gin.Context) {
		ctx.File("./docs/api/v1.yaml")
	})
	docs.Static("/v1", "./web/swagger-ui")
}

func (c *Config) authRoute(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/login", c.AuthHandler.Login())
	auth.GET("/check/admin", c.Middleware.Authenticate(), c.Middleware.RequireRole("admin"), func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "you are admin"})
	})
}

func (c *Config) hopeCornerRoute(r *gin.RouterGroup) {
	hopes := r.Group("/hopes")
	hopes.POST("", c.HopeWhisperHandler.Create(model.HopeCorner))
	hopes.GET("/:id", c.HopeWhisperHandler.FindByID(model.HopeCorner, false))
	hopes.GET("", c.HopeWhisperHandler.FindByLazyLoad(model.HopeCorner, false))
}

func (c *Config) whisperWallRoute(r *gin.RouterGroup) {
	whispers := r.Group("/whispers")
	whispers.POST("", c.HopeWhisperHandler.Create(model.WhisperWall))
	whispers.GET("/:id", c.HopeWhisperHandler.FindByID(model.WhisperWall, false))
	whispers.GET("", c.HopeWhisperHandler.FindByLazyLoad(model.WhisperWall, false))
}

func (c *Config) adminRoute(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	admin.Use(c.Middleware.Authenticate(), c.Middleware.RequireRole("admin"))
	// Admin Hope Corner
	admin.GET("/hopes/:id", c.HopeWhisperHandler.FindByID(model.HopeCorner, true))
	admin.GET("/hopes", c.HopeWhisperHandler.FindByLazyLoad(model.HopeCorner, true))
	admin.PATCH("/hopes/:id", c.HopeWhisperHandler.Update(model.HopeCorner))
	admin.DELETE("/hopes/:id", c.HopeWhisperHandler.Delete(model.HopeCorner))
	// Admin Whisper Wall
	admin.GET("/whispers/:id", c.HopeWhisperHandler.FindByID(model.WhisperWall, true))
	admin.GET("/whispers", c.HopeWhisperHandler.FindByLazyLoad(model.WhisperWall, true))
	admin.PATCH("/whispers/:id", c.HopeWhisperHandler.Update(model.WhisperWall))
	admin.DELETE("/whispers/:id", c.HopeWhisperHandler.Delete(model.WhisperWall))
}
