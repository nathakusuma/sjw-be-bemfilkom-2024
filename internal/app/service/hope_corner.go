package service

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type hopeCornerService struct {
	r repository.IHopeCornerRepository
}

type IHopeCornerService interface {
	Create(content string) response.ApiResponse
	FindByLazyLoad(afterCreatedAt, afterId, limitStr string, isAdmin bool) response.ApiResponse
	Update(idStr string, req model.UpdateHopeRequest) response.ApiResponse
}

func NewHopeCornerService(r repository.IHopeCornerRepository) IHopeCornerService {
	return &hopeCornerService{r: r}
}

func (s *hopeCornerService) Create(content string) response.ApiResponse {
	id, err := s.r.Create(content)
	if err != nil {
		return response.NewApiResponse(500, "fail to create hope", err)
	}

	return response.NewApiResponse(201, "hope created", gin.H{"id": id})
}

func (s *hopeCornerService) FindByLazyLoad(afterCreatedAt, afterId, limitStr string, isAdmin bool) response.ApiResponse {
	var MAX_FETCH = 10
	if isAdmin {
		MAX_FETCH = 20
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		return response.NewApiResponse(400, "invalid limit format", err)
	}

	if limit > MAX_FETCH {
		return response.NewApiResponse(400, "limit exceeds maximum fetch", nil)
	}

	isAfterCreateAtExist := afterCreatedAt != ""
	isAfterIdExist := afterId != ""

	if isAfterCreateAtExist != isAfterIdExist {
		return response.NewApiResponse(400, "after_created_at and after_id must be provided together", nil)
	}

	var hopesRaw []entity.Hope

	if !isAfterCreateAtExist && !isAfterIdExist {
		hopesRaw, err = s.r.FindByLazyLoad(time.Time{}, uuid.Nil, limit, isAdmin)
		if err != nil {
			return response.NewApiResponse(500, "fail to get hopes", err)
		}
	} else {
		afterTime, err := time.Parse(time.RFC3339, afterCreatedAt)
		if err != nil {
			return response.NewApiResponse(400, "invalid time format", err)
		}

		afterUuid, err := uuid.Parse(afterId)
		if err != nil {
			return response.NewApiResponse(400, "invalid uuid format", err)
		}

		hopesRaw, err = s.r.FindByLazyLoad(afterTime, afterUuid, limit, isAdmin)
		if err != nil {
			return response.NewApiResponse(500, "fail to get hopes", err)
		}
	}

	hopes := make([]any, len(hopesRaw))
	for i, hope := range hopesRaw {
		res := model.GetHopeResponse{
			ID:      hope.ID,
			Content: hope.Content,
		}
		if isAdmin {
			hopes[i] = model.GetHopeAsAdminResponse{
				GetHopeResponse: res,
				IsApproved:      hope.IsApproved,
			}
		}
		hopes[i] = res
	}

	return response.NewApiResponse(200, "hopes retrieved", gin.H{"hopes": hopes})
}

func (s *hopeCornerService) Update(idStr string, req model.UpdateHopeRequest) response.ApiResponse {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.NewApiResponse(400, "invalid id", err)
	}

	update := entity.Hope{
		ID:         id,
		Content:    req.Content,
		IsApproved: req.IsApproved,
	}

	if err := s.r.Update(update); err != nil {
		return response.NewApiResponse(500, "fail to update hope", err)
	}

	return response.NewApiResponse(201, "hope updated", nil)
}
