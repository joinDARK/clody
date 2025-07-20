package main

import (
	"dater/backend/internal/config"
	"dater/backend/internal/server"
	"fmt"
)

func main() {	
	// Инициализация конфигурации
	config, err := config.NewConfig("../configs/config.toml")
	if err != nil {
		fmt.Println("\x1b[31mError to initialize config:\x1b[0m " + err.Error())
		return
	}

	// Инициализация роутера
	server := server.NewServer(config)
	
	// Запуск сервера
	server.Start(&config.Server)
}
