package service

import (
	"errors"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type hopeWhisperService struct {
	r repository.IHopeWhisperRepository
}

type IHopeWhisperService interface {
	Create(hwType model.HopeWhisperType, content string) response.ApiResponse
	FindByLazyLoad(hwType model.HopeWhisperType, createdAtPivot, idPivot, direction, limitStr string, isAdmin bool) response.ApiResponse
	FindByID(hwType model.HopeWhisperType, idStr string, isAdmin bool) response.ApiResponse
	Update(hwType model.HopeWhisperType, idStr string, req model.UpdateHopeWhisperRequest) response.ApiResponse
	Delete(hwType model.HopeWhisperType, idStr string) response.ApiResponse
}

func NewHopeWhisperService(r repository.IHopeWhisperRepository) IHopeWhisperService {
	return &hopeWhisperService{r: r}
}

func (s *hopeWhisperService) Create(hwType model.HopeWhisperType, content string) response.ApiResponse {
	id, err := s.r.Create(hwType, content)
	if err != nil {
		return response.NewApiResponse(500, "fail to create "+hwType.Singular(), err)
	}

	return response.NewApiResponse(201, hwType.Singular()+" created", gin.H{"id": id})
}

func (s *hopeWhisperService) FindByLazyLoad(hwType model.HopeWhisperType, createdAtPivot, idPivot, direction, limitStr string, isAdmin bool) response.ApiResponse {
	var MAX_FETCH = 10
	if isAdmin {
		MAX_FETCH = 20
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		return response.NewApiResponse(400, "invalid limit format", err)
	}

	if limit > MAX_FETCH {
		return response.NewApiResponse(400, "limit exceeds maximum fetch", gin.H{})
	}

	isCreateAtPivotExist := createdAtPivot != ""
	isIdPivotExist := idPivot != ""

	if isCreateAtPivotExist != isIdPivotExist {
		return response.NewApiResponse(400, "created_at_pivot and id_pivot must be provided together", gin.H{})
	}

	isPrev := direction == "prev"

	var hopesWhispersRaw []entity.HopeWhisper

	if !isCreateAtPivotExist && !isIdPivotExist {
		hopesWhispersRaw, err = s.r.FindByLazyLoad(hwType, time.Time{}, uuid.Nil, isPrev, limit, isAdmin)
		if err != nil {
			return response.NewApiResponse(500, "fail to get "+hwType.String(), err)
		}
	} else {
		timePivot, err := time.Parse(time.RFC3339, createdAtPivot)
		if err != nil {
			return response.NewApiResponse(400, "invalid time format", err)
		}

		uuidPivot, err := uuid.Parse(idPivot)
		if err != nil {
			return response.NewApiResponse(400, "invalid uuid format", err)
		}

		hopesWhispersRaw, err = s.r.FindByLazyLoad(hwType, timePivot, uuidPivot, isPrev, limit, isAdmin)
		if err != nil {
			return response.NewApiResponse(500, "fail to get "+hwType.String(), err)
		}
	}

	hopesWhispers := make([]any, len(hopesWhispersRaw))
	for i, hopeWhisper := range hopesWhispersRaw {
		res := model.FindHopeWhisperResponse{
			ID:        hopeWhisper.ID,
			Content:   hopeWhisper.Content,
			CreatedAt: hopeWhisper.CreatedAt.Format(time.RFC3339),
		}
		if isAdmin {
			hopesWhispers[i] = model.FindHopeWhisperAsAdminResponse{
				FindHopeWhisperResponse: res,
				IsApproved:              hopeWhisper.IsApproved,
				UpdatedAt:               hopeWhisper.UpdatedAt.Format(time.RFC3339),
			}
		}
		hopesWhispers[i] = res
	}

	return response.NewApiResponse(200, hwType.String()+" retrieved", hopesWhispers)
}

func (s *hopeWhisperService) FindByID(hwType model.HopeWhisperType, idStr string, isAdmin bool) response.ApiResponse {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.NewApiResponse(400, "invalid id", err)
	}

	hopeWhisperRaw, err := s.r.FindByID(hwType, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewApiResponse(404, hwType.Singular()+" not found", err)
		}
		return response.NewApiResponse(500, "fail to get "+hwType.Singular(), err)
	}

	if !isAdmin && (hopeWhisperRaw.IsApproved == nil || !*hopeWhisperRaw.IsApproved) {
		return response.NewApiResponse(403, hwType.Singular()+" not approved yet", gin.H{})
	}

	var hopeWhisper any
	res := model.FindHopeWhisperResponse{
		ID:      hopeWhisperRaw.ID,
		Content: hopeWhisperRaw.Content,
	}
	if isAdmin {
		hopeWhisper = model.FindHopeWhisperAsAdminResponse{
			FindHopeWhisperResponse: res,
			IsApproved:              hopeWhisperRaw.IsApproved,
		}
	} else {
		hopeWhisper = res
	}

	return response.NewApiResponse(200, hwType.Singular()+" retrieved", hopeWhisper)
}

func (s *hopeWhisperService) Update(hwType model.HopeWhisperType, idStr string, req model.UpdateHopeWhisperRequest) response.ApiResponse {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.NewApiResponse(400, "invalid id", err)
	}

	update := entity.HopeWhisper{
		Content:    req.Content,
		ID:         id,
		IsApproved: req.IsApproved,
	}

	if err := s.r.Update(hwType, update); err != nil {
		return response.NewApiResponse(500, "fail to update "+hwType.Singular(), err)
	}

	return response.NewApiResponse(201, hwType.Singular()+" updated", gin.H{})
}

func (s *hopeWhisperService) Delete(hwType model.HopeWhisperType, idStr string) response.ApiResponse {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return response.NewApiResponse(400, "invalid id", err)
	}

	if err := s.r.Delete(hwType, id); err != nil {
		return response.NewApiResponse(500, "fail to delete "+hwType.Singular(), err)
	}

	return response.NewApiResponse(200, hwType.Singular()+" deleted", gin.H{})
}
