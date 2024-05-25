package rest

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/service"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type hopeWhisperHandler struct {
	s service.IHopeWhisperService
}

type IHopeWhisperHandler interface {
	Create(hwType model.HopeWhisperType) gin.HandlerFunc
	FindByLazyLoad(hwType model.HopeWhisperType, isAdmin bool) gin.HandlerFunc
	FindByID(hwType model.HopeWhisperType, isAdmin bool) gin.HandlerFunc
	Update(hwType model.HopeWhisperType) gin.HandlerFunc
	Delete(hwType model.HopeWhisperType) gin.HandlerFunc
}

func NewHopeWhisperHandler(service service.IHopeWhisperService) IHopeWhisperHandler {
	return &hopeWhisperHandler{s: service}
}

func (h *hopeWhisperHandler) Create(hwType model.HopeWhisperType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.CreateHopeWhisperRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.NewApiResponse(400, "invalid request body", err).Send(ctx)
			return
		}

		res := h.s.Create(hwType, req.Content)
		res.Send(ctx)
	}
}

func (h *hopeWhisperHandler) FindByLazyLoad(hwType model.HopeWhisperType, isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		afterCreatedAt := ctx.Query("after_created_at")
		afterId := ctx.Query("after_id")
		limit := ctx.Query("limit")

		res := h.s.FindByLazyLoad(hwType, afterCreatedAt, afterId, limit, isAdmin)
		res.Send(ctx)
	}
}

func (h *hopeWhisperHandler) FindByID(hwType model.HopeWhisperType, isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		res := h.s.FindByID(hwType, id, isAdmin)
		res.Send(ctx)
	}
}

func (h *hopeWhisperHandler) Update(hwType model.HopeWhisperType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		var req model.UpdateHopeWhisperRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			response.NewApiResponse(400, "invalid request body", err).Send(ctx)
			return
		}

		res := h.s.Update(hwType, id, req)
		res.Send(ctx)
	}
}

func (h *hopeWhisperHandler) Delete(hwType model.HopeWhisperType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		res := h.s.Delete(hwType, id)
		res.Send(ctx)
	}
}
