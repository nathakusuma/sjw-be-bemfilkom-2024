package repository

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type IHopeWhisperRepository interface {
	Create(hwType model.HopeWhisperType, content string, isPublic bool) (uuid.UUID, error)
	FindByLazyLoad(hwType model.HopeWhisperType, createdAtPivot time.Time, idPivot uuid.UUID, isPrev bool, limit int, isAdmin bool) ([]entity.HopeWhisper, error)
	FindByID(hwType model.HopeWhisperType, id uuid.UUID) (entity.HopeWhisper, error)
	FindAllApproved(hwType model.HopeWhisperType) ([]entity.HopeWhisper, error)
	Update(hwType model.HopeWhisperType, hopeWhisper entity.HopeWhisper) error
	Delete(hwType model.HopeWhisperType, id uuid.UUID) error
}

type hopeWhisperRepository struct {
	db *gorm.DB
}

func NewHopeWhisperRepository(db *gorm.DB) IHopeWhisperRepository {
	return &hopeWhisperRepository{db: db}
}

func (r *hopeWhisperRepository) Create(hwType model.HopeWhisperType, content string, isPublic bool) (uuid.UUID, error) {
	hopeWhisper := entity.HopeWhisper{
		ID:       uuid.New(),
		Content:  content,
		IsPublic: isPublic,
	}
	if err := r.db.Table(hwType.String()).Create(&hopeWhisper).Error; err != nil {
		return uuid.Nil, err
	}

	return hopeWhisper.ID, nil
}

func (r *hopeWhisperRepository) FindByLazyLoad(hwType model.HopeWhisperType, createdAtPivot time.Time, idPivot uuid.UUID, isPrev bool, limit int, isAdmin bool) ([]entity.HopeWhisper, error) {
	var hopes []entity.HopeWhisper

	tx := r.db.Table(hwType.String())
	if !isAdmin {
		tx = tx.Where("is_approved = ? AND is_public = ?", true, true)
	}

	if (createdAtPivot != time.Time{} && idPivot != uuid.Nil) {
		createdAtOperator := "<"
		idOperator := ">"
		if isPrev {
			createdAtOperator = ">"
			idOperator = "<"
		}
		tx = tx.Where("created_at "+createdAtOperator+" ? OR (created_at = ? AND id "+idOperator+" ?)", createdAtPivot, createdAtPivot, idPivot)
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

// FindAllApproved Not a good practice to return all data, so I limit it to 100
func (r *hopeWhisperRepository) FindAllApproved(hwType model.HopeWhisperType) ([]entity.HopeWhisper, error) {
	var hopeWhispers []entity.HopeWhisper
	tx := r.db.Table(hwType.String()).
		Where("is_approved = ? AND is_public = ?", true, true).
		Order("created_at DESC, id ASC").
		Limit(100).
		Find(&hopeWhispers)
	if err := tx.Error; err != nil {
		return nil, err
	}
	return hopeWhispers, nil
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
