package database

import (
	"dater/backend/internal/config"
	"dater/backend/internal/logger"
	"dater/backend/internal/models"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Database) (*gorm.DB, error) {
	// Загружаем переменные окружения
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}

	// Получаем пароль из переменной окружения
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("В .env пустое поле DB_PASSWORD или его нету")
	}

	// Формируем строку подключения к базе данных
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.DBUser, dbPassword, cfg.DBName, cfg.Port)

	// Открываем соединение с БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Error connecting to database")
		return nil, err
	}
	logger.Log.Debug().
		Str("host", cfg.Host).
		Str("db_name", cfg.DBName).
		Str("db_user", cfg.DBUser).
		Str("db", db.Name()).
		Msg("Database connected")

	// Создаем enum в бд для типов полей
	db.Exec(`
		DO $$
		BEGIN
		  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'field_type') THEN
		    CREATE TYPE field_type AS ENUM (
		      'text',
		      'number',
		      'checkbox',
		      'date',
		      'single_select',
		      'multi_select',
		      'many2one',
		      'one2many',
		      'many2many',
		      'formula'
		    );
		  END IF;
		END
		$$;
	`)

	// Миграция базы данных
	err = db.AutoMigrate(
		&models.Base{},
		&models.Table{},
		&models.Field{},
		&models.SelectOption{},
		&models.Record{},
		&models.Cell{},
		&models.RecordRelation{},
	)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Failed to migrate database")
		return nil, err
	}
	logger.Log.Debug().Msg("Migration completed")
	return db, nil
}
