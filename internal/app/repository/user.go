package repository

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	FindByNim(nim string) (entity.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByNim(nim string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("nim = ?", nim).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
