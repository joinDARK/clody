package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "02.01.2006 15:04:05",
		NoColor:    false,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			"message",
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("%s", i)
		},
	}
	
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	Log = zerolog.New(output).With().Timestamp().Logger()
}
