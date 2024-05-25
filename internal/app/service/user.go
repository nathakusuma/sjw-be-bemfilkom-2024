package service

import (
	"errors"
	ubauth "github.com/ahmdyaasiin/ub-auth-without-notification/v2"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userService struct {
	r   repository.IUserRepository
	jwt jwt.IJWT
}

type IUserService interface {
	Login(nimOrEmail, password string) response.ApiResponse
}

func NewUserService(repository repository.IUserRepository, jwt jwt.IJWT) IUserService {
	return &userService{r: repository, jwt: jwt}
}

func (s *userService) Login(username, password string) response.ApiResponse {
	studentDetails, err := ubauth.AuthUB(username, password)
	if err != nil {
		var respErr *ubauth.ResponseDetails
		ok := errors.As(err, &respErr)
		if !ok {
			return response.NewApiResponse(500, "fail to authenticate to ub auth", err)
		}
		return response.NewApiResponse(respErr.Code, respErr.Message, nil)
	}

	user, err := s.r.FindByNim(studentDetails.NIM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewApiResponse(404, "user not found", nil)
		}
		return response.NewApiResponse(500, "fail to get user data", err)
	}

	token, err := s.jwt.Create(user)
	if err != nil {
		return response.NewApiResponse(500, "fail to generate token", err)
	}

	return response.NewApiResponse(200, "successfully logged in", gin.H{"token": token})
}
