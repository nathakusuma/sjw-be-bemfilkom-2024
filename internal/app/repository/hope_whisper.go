package repository

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IHopeWhisperRepository interface {
	Create(hwType model.HopeWhisperType, content string) (uuid.UUID, error)
	FindByLazyLoad(hwType model.HopeWhisperType, afterCreatedAt time.Time, afterId uuid.UUID, limit int, isAdmin bool) ([]entity.HopeWhisper, error)
	FindByID(hwType model.HopeWhisperType, id uuid.UUID) (entity.HopeWhisper, error)
	Update(hwType model.HopeWhisperType, hopeWhisper entity.HopeWhisper) error
	Delete(hwType model.HopeWhisperType, id uuid.UUID) error
}

type hopeWhisperRepository struct {
	db *gorm.DB
}

func NewHopeWhisperRepository(db *gorm.DB) IHopeWhisperRepository {
	return &hopeWhisperRepository{db: db}
}

func (r *hopeWhisperRepository) Create(hwType model.HopeWhisperType, content string) (uuid.UUID, error) {
	hopeWhisper := entity.HopeWhisper{
		ID:      uuid.New(),
		Content: content,
	}
	if err := r.db.Table(hwType.String()).Create(&hopeWhisper).Error; err != nil {
		return uuid.Nil, err
	}

	return hopeWhisper.ID, nil
}

func (r *hopeWhisperRepository) FindByLazyLoad(hwType model.HopeWhisperType, afterCreatedAt time.Time, afterId uuid.UUID, limit int, isAdmin bool) ([]entity.HopeWhisper, error) {
	var hopes []entity.HopeWhisper

	tx := r.db.Table(hwType.String())
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

func (r *hopeWhisperRepository) FindByID(hwType model.HopeWhisperType, id uuid.UUID) (entity.HopeWhisper, error) {
	var hopeWhisper entity.HopeWhisper
	if err := r.db.Table(hwType.String()).Where("id = ?", id).First(&hopeWhisper).Error; err != nil {
		return hopeWhisper, err
	}
	return hopeWhisper, nil
}

func (r *hopeWhisperRepository) Update(hwType model.HopeWhisperType, hopeWhisper entity.HopeWhisper) error {
	return r.db.Table(hwType.String()).Where("id = ?", hopeWhisper.ID).Updates(&hopeWhisper).Error
}

func (r *hopeWhisperRepository) Delete(hwType model.HopeWhisperType, id uuid.UUID) error {
	return r.db.Table(hwType.String()).Delete(&entity.HopeWhisper{}, id).Error
}
