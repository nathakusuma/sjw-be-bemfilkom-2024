package service

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type hopeCornerService struct {
	r repository.IHopeCornerRepository
}

type IHopeCornerService interface {
	Create(content string) response.ApiResponse
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
