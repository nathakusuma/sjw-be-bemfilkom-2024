package config

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if err := migrateTables(db); err != nil {
		log.Fatalln(err)
	}

	return db
}

func migrateTables(db *gorm.DB) error {
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
