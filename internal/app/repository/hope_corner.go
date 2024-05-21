package repository

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IHopeCornerRepository interface {
	Create(content string) (uuid.UUID, error)
	FindByLazyLoad(afterCreatedAt time.Time, afterId uuid.UUID, limit int, isAdmin bool) ([]entity.Hope, error)
	FindByID(id uuid.UUID) (entity.Hope, error)
	Update(hope entity.Hope) error
	Delete(id uuid.UUID) error
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

func (r *hopeCornerRepository) FindByLazyLoad(afterCreatedAt time.Time, afterId uuid.UUID, limit int, isAdmin bool) ([]entity.Hope, error) {
	var hopes []entity.Hope

	tx := r.db
	if !isAdmin {
		tx = tx.Where("is_approved = ?", true)
	}

	if (afterCreatedAt != time.Time{} && afterId != uuid.Nil) {
		tx = tx.Where("created_at < ? OR (created_at = ? AND id > ?)", afterCreatedAt, afterCreatedAt, afterId)
	}

	tx = tx.
		Order("created_at DESC, id ASC").
		Limit(limit).
		Find(&hopes)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return hopes, nil
}

func (r *hopeCornerRepository) FindByID(id uuid.UUID) (entity.Hope, error) {
	var hope entity.Hope
	if err := r.db.Where("id = ?", id).First(&hope).Error; err != nil {
		return hope, err
	}
	return hope, nil
}

func (r *hopeCornerRepository) Update(hope entity.Hope) error {
	return r.db.Model(&entity.Hope{}).Where("id = ?", hope.ID).Updates(&hope).Error
}

func (r *hopeCornerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Hope{}, id).Error
}
