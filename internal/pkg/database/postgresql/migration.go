package postgresql

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.Table("hopes").AutoMigrate(&entity.HopeWhisper{}); err != nil {
		return err
	}

	if err := db.Table("whispers").AutoMigrate(&entity.HopeWhisper{}); err != nil {
		return err
	}

	return nil
}
