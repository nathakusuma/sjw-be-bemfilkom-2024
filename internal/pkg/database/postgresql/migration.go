package postgresql

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
