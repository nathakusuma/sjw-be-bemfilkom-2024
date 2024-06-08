package main

import (
	"github.com/bem-filkom/sjw-be-2024/internal/app/config"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	env := os.Getenv("ENV")
	if err != nil && env == "" {
		log.Println("fail to load env")
	}

	db := config.NewDatabase()
	app := config.NewGin()

	config.StartApp(&config.StartAppConfig{
		DB:  db,
		App: app,
	})

	if err := app.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalln(err)
	}
}
