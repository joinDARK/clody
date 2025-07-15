package database

import (
	"dater/backend/internal/config"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Database) (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// Получаем пароль из переменной окружения
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("В .env пустое поле DB_PASSWORD или его нету")
	}
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.DBUser, dbPassword, cfg.DBName, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
