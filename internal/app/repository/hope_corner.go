package repository

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IHopeCornerRepository interface {
	Create(content string) (uuid.UUID, error)
}

type hopeCornerRepository struct {
	db *gorm.DB
}

func NewHopeCornerRepository(db *gorm.DB) IHopeCornerRepository {
	return &hopeCornerRepository{db: db}
}

func (r *hopeCornerRepository) Create(content string) (uuid.UUID, error) {
	hope := entity.Hope{
		ID:      uuid.New(),
		Content: content,
	}
	if err := r.db.Create(&hope).Error; err != nil {
		return uuid.Nil, err
	}

	return hope.ID, nil
}
