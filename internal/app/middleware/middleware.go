package middleware

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	Authenticate(ctx *gin.Context)
	RequireRole(role string) gin.HandlerFunc
}

type middleware struct {
	jwtAuth jwt.IJWT
}

func NewMiddleware(jwtAuth jwt.IJWT) IMiddleware {
	return middleware{jwtAuth: jwtAuth}
}
