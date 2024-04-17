package main

import (
	"github.com/bem-filkom/sjw-be-2024/internal/pkg/database/postgresql"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	_ = postgresql.Connect()
}
