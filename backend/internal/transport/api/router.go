package api

import (
	"dater/backend/internal/config"
	"dater/backend/internal/repo/postgres"
	"dater/backend/internal/service"
	"fmt"

	"strconv"

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
		base.GET("/:id", GetInfoBaseHandler(baseService, logger))
		base.POST("/", CreateBaseHandler(baseService, logger))
		base.PUT("/:id", UpdateBaseHandler(baseService, logger))
		base.DELETE("/:id", DeleteBaseHandler(baseService, logger))
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

func GetInfoBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService
	
	return func(c *gin.Context) {
		logger.Debug().Msg("Getting base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}
		
		base, err := service.GetBaseInfo(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get base")
			c.JSON(500, gin.H{"error": "Failed to get base"})
			return
		}

		logger.Info().Msg("Base retrieved")
		c.JSON(200, gin.H{
			"data": base,
			"_links": map[string]map[string]string{
				"self": {
					"href": fmt.Sprintf("/api/bases/%d", base.ID),
				},
				"create": {
					"href": "/api/bases",
					"method": "POST",
				},
				"update": {
					"href": fmt.Sprintf("/api/bases/%d", base.ID),
					"method": "PUT",
				},
				"delete": {
					"href": fmt.Sprintf("/api/bases/%d", base.ID),
					"method": "DELETE",
				},
			},
		})
	}
}

type CreateBaseInput struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

func CreateBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService
	
	return func(c *gin.Context) {
		logger.Debug().Msg("Creating base...")

		var input CreateBaseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		newBase, err := service.CreateBase(input.Name, input.Description, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create base")
			c.JSON(500, gin.H{"error": "Failed to create base"})
			return
		}

		logger.Info().Msg("Base created")
		c.JSON(201, gin.H{"data": newBase})
	}
}

type UpdateBaseInput struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description *string `json:"description" binding:"omitempty"`
}

func UpdateBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService
	
	return func(c *gin.Context) {
		logger.Debug().Msg("Updating base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		var input UpdateBaseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			logger.Error().Err(err).Msg("Failed to bind JSON")
			c.JSON(400, gin.H{"error": "Failed to bind JSON"})
			return
		}

		updBase, err := service.UpdateBase(i, input.Name, input.Description, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to update base")
			c.JSON(500, gin.H{"error": "Failed to update base"})
			return
		}

		logger.Info().Msg("Base updated")
		c.JSON(200, gin.H{"data": updBase})
	}
}

func DeleteBaseHandler(baseService *service.BaseService, logger zerolog.Logger) gin.HandlerFunc {
	service := *baseService
	
	return func(c *gin.Context) {
		logger.Debug().Msg("Deleting base...")

		id := c.Param("id")
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse ID")
			c.JSON(400, gin.H{"error": "Failed to parse ID"})
			return
		}

		deleteID, err := service.DeleteBase(i, logger)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to delete base")
			c.JSON(500, gin.H{"error": "Failed to delete base"})
			return
		}

		logger.Info().Msg("Base deleted")
		c.JSON(200, gin.H{"data": deleteID})
	}
}

// type CreateTableInput struct {
// 	BaseID      int64   `json:"base_id" binding:"required"`
// 	Name        string  `json:"name" binding:"required,min=3,max=100"`
// 	Description *string `json:"description" binding:"omitempty"`
// }

// type UpdateTableInput struct {
// 	BaseID      int64  `json:"base_id" binding:"omitempty"`
// 	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
// 	Description string `json:"description" binding:"omitempty"`
// }