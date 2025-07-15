package main

import (
	"dater/backend/internal/config"
	
	"dater/backend/internal/logger"
	"dater/backend/internal/router"
)

func main() {
	logger.Init()
	// Инициализация конфигурации
	config, err := config.NewConfig("../configs/config.toml")
	if err != nil {
		logger.Log.Error().Msg("Не удалось загрузить конфигурацию: " + err.Error())
	}
	
	r := router.NewRouter(&config.Server)
	logger.Log.Info().Msg("Starting server: http://localhost:8080")
	r.Run(":8080")
}