package main

import (
	"dater/backend/internal/config"
	"dater/backend/internal/database"
	"dater/backend/internal/server"
	"fmt"

	"dater/backend/internal/logger"
)

func main() {	
	// Инициализация конфигурации
	config, err := config.NewConfig("../configs/config.toml")
	if err != nil {
		fmt.Println("\x1b[31mError to initialize config:\x1b[0m " + err.Error())
		return
	}

	// Инициализация логгера
	logger.Init(&config.Logger)

	// Инициализация базы данных
	db, err := database.NewDB(&config.Database)
	if err != nil {
		logger.Log.Error().Msg("Error to initialize database: " + err.Error())
	} else {
		logger.Log.Info().Msg("Database connected")
	}

	// Инициализация роутера
	server := server.NewServer(db, logger.Log)
	
	// Запуск сервера
	server.Start(&config.Server)
}
