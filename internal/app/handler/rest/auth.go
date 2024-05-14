package rest

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	s service.IUserService
}

type IAuthHandler interface {
	Login(ctx *gin.Context)
}

func NewAuthHandler(service service.IUserService) IAuthHandler {
	return &authHandler{s: service}
}

func (h *authHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewApiResponse(400, "invalid request body", err).Send(ctx)
		return
	}

	res := h.s.Login(req.Username, req.Password)
	res.Send(ctx)
}
