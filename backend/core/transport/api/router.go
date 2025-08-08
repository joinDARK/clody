package api

import (
	"clody/core/config"
	"clody/core/repo/postgres"
	"clody/core/service"
	"clody/core/transport/api/handlers"

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

	{
		fieldRepo := postgres.NewPostgresFieldRepo(db)
		fieldService := service.NewFieldService(fieldRepo)

		field := apiGroup.Group("/fields")
		field.GET("/:id", handlers.GetFieldByIDHandler(fieldService, logger))
		field.POST("/", handlers.CreateFieldHandler(fieldService, logger))
		field.PATCH("/:id", handlers.UpdateFieldHandler(fieldService, logger))
		field.DELETE("/:id", handlers.DeleteFieldHandler(fieldService, logger))
	}

	{
		recordRepo := postgres.NewPostgresRecordRepo(db)
		recordService := service.NewRecordService(recordRepo)

		record := apiGroup.Group("/records")
		record.GET("/:id", handlers.GetRecordByIDHandler(recordService, logger))
		record.POST("/", handlers.CreateRecordHandler(recordService, logger))
		record.PATCH("/:id", handlers.UpdateRecordHandler(recordService, logger))
		record.DELETE("/:id", handlers.DeleteRecordHandler(recordService, logger))
	}

	{
		cellRepo := postgres.NewPostgresCellRepo(db)
		cellService := service.NewCellService(cellRepo)

		cell := apiGroup.Group("/cells")
		cell.GET("/:id", handlers.GetCellByIDHandler(cellService, logger))
		cell.POST("/", handlers.CreateCellHandler(cellService, logger))
		cell.PATCH("/:id", handlers.UpdateCellHandler(cellService, logger))
		cell.DELETE("/:id", handlers.DeleteCellHandler(cellService, logger))
	}

	{
		recordRelationRepo := postgres.NewRecordRelationRepo(db)
		recordRelationService := service.NewRecordRelationService(recordRelationRepo)

		recordRelation := apiGroup.Group("/record_relations")
		recordRelation.GET("/:id", handlers.GetRecordRelationByIDHandler(recordRelationService, logger))
		recordRelation.POST("/", handlers.CreateRecordRelationHandler(recordRelationService, logger))
		recordRelation.PATCH("/:id", handlers.UpdateRecordRelationHandler(recordRelationService, logger))
		recordRelation.DELETE("/:id", handlers.DeleteRecordRelationHandler(recordRelationService, logger))
	}

	{
		selectOptionRepo := postgres.NewPostgresSelectOptionRepo(db)
		selectOptionService := service.NewSelectOptionService(selectOptionRepo)

		selectOption := apiGroup.Group("/select_options")
		selectOption.GET("/:id", handlers.GetSelectOptionByIDHandler(selectOptionService, logger))
		selectOption.POST("/", handlers.CreateSelectOptionHandler(selectOptionService, logger))
		selectOption.PATCH("/:id", handlers.UpdateSelectOptionHandler(selectOptionService, logger))
		selectOption.DELETE("/:id", handlers.DeleteSelectOptionHandler(selectOptionService, logger))
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
