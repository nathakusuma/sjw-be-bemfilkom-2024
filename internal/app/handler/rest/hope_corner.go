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
	FindByLazyLoad(isAdmin bool) gin.HandlerFunc
	FindByID(isAdmin bool) gin.HandlerFunc
	Update(ctx *gin.Context)
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

func (h *hopeCornerHandler) FindByLazyLoad(isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		afterCreatedAt := ctx.Query("after_created_at")
		afterId := ctx.Query("after_id")
		limit := ctx.Query("limit")

		res := h.s.FindByLazyLoad(afterCreatedAt, afterId, limit, isAdmin)
		res.Send(ctx)
	}
}

func (h *hopeCornerHandler) FindByID(isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		res := h.s.FindByID(id, isAdmin)
		res.Send(ctx)
	}
}

func (h *hopeCornerHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req model.UpdateHopeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewApiResponse(400, "invalid request body", err).Send(ctx)
		return
	}

	res := h.s.Update(id, req)
	res.Send(ctx)
}
