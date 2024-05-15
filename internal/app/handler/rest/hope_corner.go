package rest

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type hopeCornerHandler struct {
	s service.IHopeCornerService
}

type IHopeCornerHandler interface {
	Create(ctx *gin.Context)
	GetLazyLoad(isAdmin bool) gin.HandlerFunc
}

func NewHopeCornerHandler(service service.IHopeCornerService) IHopeCornerHandler {
	return &hopeCornerHandler{s: service}
}

func (h *hopeCornerHandler) Create(ctx *gin.Context) {
	var req model.CreateHopeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewApiResponse(400, "invalid request body", err).Send(ctx)
		return
	}

	res := h.s.Create(req.Content)
	res.Send(ctx)
}

func (h *hopeCornerHandler) GetLazyLoad(isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		afterCreatedAt := ctx.Query("after_created_at")
		afterId := ctx.Query("after_id")
		limit := ctx.Query("limit")

		res := h.s.GetLazyLoad(afterCreatedAt, afterId, limit, isAdmin)
		res.Send(ctx)
	}
}
