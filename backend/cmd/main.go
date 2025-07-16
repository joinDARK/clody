package main

import (
	"dater/backend/internal/config"
	"dater/backend/internal/database"

	"dater/backend/internal/logger"
	"dater/backend/internal/router"
)

func main() {
	// Инициализация конфигурации
	config, err := config.NewConfig("../configs/config.toml")
	if err != nil {
		logger.Init(&config.Logger)
		logger.Log.Error().Msg("Не удалось загрузить конфигурацию: " + err.Error())
		return
	}

	// Инициализация логгера
	logger.Init(&config.Logger)

	// Инициализация базы данных
	db, err := database.NewDB(&config.Database)
	if err != nil {
		logger.Log.Error().Msg("Errro to initialize database: " + err.Error())
	} else {
		logger.Log.Info().Msg("Database connected")
	}
	_ = db // Временно убираем ошибку линтера UnusedVar

	// Инициализация роутера
	r := router.NewRouter(&config.Server)
	logger.Log.Info().Msg("Router is initialized")
	
	// Запуск сервера
	r.Run(":8080")
	logger.Log.Info().Msg("Starting server: http://localhost:8080")
}
