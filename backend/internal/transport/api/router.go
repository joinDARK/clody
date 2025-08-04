package api

import (
	"dater/backend/internal/config"
	"dater/backend/internal/repo/postgres"
	"dater/backend/internal/service"
	"dater/backend/internal/transport/api/handlers"

	// "gorm.io/gorm"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Server, logger zerolog.Logger, db *gorm.DB) *gin.Engine {
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
	apiGroup := r.Group("/api")

	{
		baseRepo := postgres.NewPostgresBaseRepo(db)
		baseService := service.NewBaseService(baseRepo)

		base := apiGroup.Group("/bases")
		base.GET("/:id", handlers.GetInfoBaseHandler(baseService, logger))
		base.POST("/", handlers.CreateBaseHandler(baseService, logger))
		base.PATCH("/:id", handlers.UpdateBaseHandler(baseService, logger))
		base.DELETE("/:id", handlers.DeleteBaseHandler(baseService, logger))
	}

	{
		tableRepo := postgres.NewPostgresTableRepo(db)
		tableService := service.NewTableService(tableRepo)

		table := apiGroup.Group("/tables")
		table.GET("/:id", handlers.GetInfoTableHandler(tableService, logger))
		table.GET("/:id/records", handlers.GetAllRecordsByTableHandler(tableService, logger))
		table.GET("/:id/fields", handlers.GetAllFieldsByTableHandler(tableService, logger))
		table.POST("/", handlers.CreateTableHandler(tableService, logger))
		table.PATCH("/:id", handlers.UpdateTableHandler(tableService, logger))
		table.DELETE("/:id", handlers.DeleteTableHandler(tableService, logger))
	}

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
