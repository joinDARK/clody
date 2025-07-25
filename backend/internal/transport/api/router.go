package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"dater/backend/internal/config"
)

func NewRouter(cfg *config.Server, logger zerolog.Logger) *gin.Engine {
	logger.Debug().Msg("Creating router...")
	
	if cfg == nil {
		logger.Error().Msg("Config is nil")
	}
	if cfg.Host == "" {
		logger.Error().Msg("Config variable host is empty")
	}
	if cfg.Port == 0 {
		logger.Error().Msg("Config variable port == 0")
	}
	if cfg.Mode == "" {
		logger.Error().Msg("Config variable mode is empty")
	}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(GinLogger(logger))

	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Debug().
		Str("mode", cfg.Mode).
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Int("routes_count", len(r.Routes())).
		Msg("Router is created")
	logger.Info().Msg("Router is created")
	return r
}