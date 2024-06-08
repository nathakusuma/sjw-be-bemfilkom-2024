package middleware

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

func (m middleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader("Authorization")
		if bearer == "" {
			response.NewApiResponse(401, "empty token", nil).Send(ctx)
			ctx.Abort()
			return
		}

		token := strings.Split(bearer, " ")[1]
		var claims jwt.Claims
		err := m.jwtAuth.Decode(token, &claims)
		if err != nil {
			response.NewApiResponse(401, "fail to validate token", err).Send(ctx)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
