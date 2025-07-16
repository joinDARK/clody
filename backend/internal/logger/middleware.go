package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обрабатываем запрос
		c.Next()

		// После обработки запроса
		duration := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		if status >= 400 {
			Log.Error().
				Str("PATH", path).
				Str("CLIENT_IP", clientIP).
				Dur("DURATION", duration).
				Msg(fmt.Sprintf("%s %d", method, status))
		} else {
			Log.Info().
				Str("PATH", path).
				Str("CLIENT_IP", clientIP).
				Dur("DURATION", duration).
				Msg(fmt.Sprintf("%s %d", method, status))
		}
	}
}
