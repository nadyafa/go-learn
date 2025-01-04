package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nadyafa/go-learn/config/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Unable to fetch .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort)
	// dsn := "host=localhost user=developer password=dev123 dbname=go-learn port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprintf("Cannot connect to database : %s", err.Error()), err)

		// panic(err)
	}

	return db, nil
}
