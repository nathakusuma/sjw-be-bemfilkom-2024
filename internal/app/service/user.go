package service

import (
	"errors"
	"fmt"
	ubauth "github.com/ahmdyaasiin/ub-auth-without-notification/v2"
	"github.com/bem-filkom/sjw-be-2024/internal/app/repository"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/jwt"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/response"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type userService struct {
	r   repository.IUserRepository
	jwt jwt.JWT
}

type IUserService interface {
	Login(username, password string) response.ApiResponse
}

func NewUserService(repository repository.IUserRepository, jwt jwt.JWT) IUserService {
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
		return response.NewApiResponse(respErr.Code, respErr.Message, gin.H{})
	}

	if studentDetails.Fakultas != "Fakultas Ilmu Komputer" {
		return response.NewApiResponse(403, "only for FILKOM students", gin.H{"fakultas": studentDetails.Fakultas})
	}

	user, err := s.r.FindByNim(studentDetails.NIM)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = entity.User{
				Nim:          studentDetails.NIM,
				Email:        studentDetails.Email,
				FullName:     studentDetails.FullName,
				ProgramStudi: studentDetails.ProgramStudi,
				Role:         "user",
			}
		} else {
			return response.NewApiResponse(500, "fail to get user data", err)
		}
	}

	angkatan := "20" + studentDetails.NIM[0:2]
	profilePictureURL := fmt.Sprintf("https://siakad.ub.ac.id/dirfoto/foto/foto_%s/%s.jpg", angkatan, studentDetails.NIM)

	claims := jwt.Claims{
		Role:           user.Role,
		FullName:       user.FullName,
		Email:          user.Email,
		ProgramStudi:   user.ProgramStudi,
		Angkatan:       angkatan,
		ProfilePicture: profilePictureURL,
		RegisteredClaims: jwt2.RegisteredClaims{
			Subject: studentDetails.NIM,
		},
	}

	token, err := s.jwt.Create(&claims)
	if err != nil {
		return response.NewApiResponse(500, "fail to generate token", err)
	}

	return response.NewApiResponse(200, "successfully logged in", model.LoginResponse{
		Token:          token,
		NIM:            studentDetails.NIM,
		Email:          studentDetails.Email,
		FullName:       studentDetails.FullName,
		Role:           user.Role,
		Angkatan:       angkatan,
		ProgramStudi:   studentDetails.ProgramStudi,
		ProfilePicture: profilePictureURL,
	})
}
