package middleware

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func (m middleware) RequireRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claimsRaw, ok := ctx.Get("claims")
		if !ok {
			response.NewApiResponse(401, "empty token", nil).Send(ctx)
			ctx.Abort()
			return
		}
		claims := claimsRaw.(jwt.Claims)
		if claims.Role != role {
			response.NewApiResponse(403, "forbidden", nil).Send(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
