package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"dater/backend/internal/config"
	"dater/backend/internal/handlers"
	"dater/backend/internal/logger"
)

func NewRouter(cfg *config.Server) *gin.Engine {
	if cfg == nil {
		logger.Log.Error().Msg("Config is nil")
	}
	if cfg.Host == "" {
		logger.Log.Error().Msg("Config variable host is empty")
	}
	if cfg.Port == 0 {
		logger.Log.Error().Msg("Config variable port == 0")
	}
	if cfg.Mode == "" {
		logger.Log.Error().Msg("Config variable mode is empty")
	}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(logger.GinLogger())

	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.GET("/ping", handlers.Ping)

	logger.Log.Debug().
		Str("mode", cfg.Mode).
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Int("routes_count", len(r.Routes())).
		Msg("Router is initialized")
	return r
}
