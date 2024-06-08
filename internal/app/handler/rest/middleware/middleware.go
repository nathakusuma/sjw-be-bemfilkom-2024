package middleware

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IMiddleware interface {
	Authenticate() gin.HandlerFunc
	RequireRole(role string) gin.HandlerFunc
	CORS() gin.HandlerFunc
	IpRateLimiter(key string, limit rate.Limit, burst int, message string) gin.HandlerFunc
}

type middleware struct {
	jwtAuth jwt.JWT
}

func NewMiddleware(jwtAuth jwt.JWT) IMiddleware {
	return middleware{jwtAuth: jwtAuth}
}
