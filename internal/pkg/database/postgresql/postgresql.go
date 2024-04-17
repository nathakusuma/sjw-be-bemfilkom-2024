package postgresql

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if err := AutoMigrate(db); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalln(err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Fatalln(err)
			return
		}
	}()

	return db
}
