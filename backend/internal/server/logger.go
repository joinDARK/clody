package server

import (
	"dater/backend/internal/config"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func InitLogger(cfg *config.Logger) zerolog.Logger {
	// Формируем вывод логера в консоль
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02.01.2006 15:04:05",
		NoColor:    false,
		PartsOrder: []string{ // Структура вывода логера в консоль
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			"message",
		},
		FormatLevel: func(i any) string { // Форматирование уровня лога
			switch i.(string) {
			case "debug":
				return "\x1b[34mDBG\x1b[0m" // Синий
			case "info":
				return "\x1b[32mINF\x1b[0m" // Зелёный
			case "warn":
				return "\x1b[33mWRN\x1b[0m" // Жёлтый
			case "error":
				return "\x1b[31mERR\x1b[0m" // Красный
			default:
				return fmt.Sprintf("%s", i)
			}
		},
		FormatMessage: func(i any) string { // Форматирование сообщения лога
			return fmt.Sprintf("%s", i)
		},
	}

	switch cfg.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Инициализируем объект логера
	return zerolog.New(output).With().Timestamp().Logger()
}
