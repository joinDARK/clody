package api

import (
	"dater/backend/internal/config"
	"dater/backend/internal/repo/postgres"
	"dater/backend/internal/service"

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
	
	{
		baseRepo := postgres.NewPostgresBaseRepo(db)
		baseService := service.NewBaseService(baseRepo)
		
		base := r.Group("/bases")
		base.GET("/:id", GetInfoBaseHandler(baseService, logger))
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
		c.JSON(200, gin.H{"data": base})
	}
}

// type CreateBaseInput struct {
// 	Name        string  `json:"name" binding:"required,min=3,max=100"`
// 	Description *string `json:"description" binding:"omitempty"`
// }

// type UpdateBaseInput struct {
// 	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
// 	Description string `json:"description" binding:"omitempty"`
// }

// func GetBases(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		logger.Debug().Msg("Getting base...")

// 		var base []Base
// 		if err := db.Find(&base).Error; err != nil {
// 			logger.Error().Err(err).Msg("Failed to get base")
// 			c.JSON(500, gin.H{"error": "Failed to get base"})
// 			return
// 		}

// 		logger.Debug().Msg("Bases retrieved")
// 		logger.Info().Msg("Bases retrieved")
// 		c.JSON(200, gin.H{"data": base})
// 	}
// }

// func CreateBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		logger.Debug().Msg("Creating base...")

// 		var input CreateBaseInput
// 		if err := c.ShouldBindJSON(&input); err.Error() == "EOF" || err != nil {
// 			logger.Error().Err(err).Msg("Failed to bind base")
// 			c.JSON(400, gin.H{"error": "Failed to bind base"})
// 			return
// 		}

// 		if input.Name == "" {
// 			input.Name = "Default"
// 		}

// 		defaultBase := &Base{
// 			Name:        input.Name,
// 			Description: input.Description,
// 		}

// 		if err := db.Create(defaultBase).Error; err != nil {
// 			logger.Error().Err(err).Msg("Failed to create base")
// 			c.JSON(500, gin.H{"error": "Failed to create base"})
// 			return
// 		}

// 		logger.Debug().Msg("Base created")
// 		logger.Info().Msg("Base created")
// 		c.JSON(201, gin.H{"data": defaultBase})
// 	}
// }

// func UpdateBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Получаем ID пользователя из пути
// 		id := c.Param("id")
// 		if id == "" {
// 			logger.Error().Msg("Empty ID in URI")
// 			c.JSON(404, gin.H{"error": "Empty ID in URI"})
// 			return
// 		}

// 		logger.Debug().Msg("Updating base...")

// 		var input UpdateBaseInput
// 		if err := c.ShouldBindJSON(&input); err != nil && err.Error() != "EOF" {
// 			logger.Error().Err(err).Msg("Failed to bind base")
// 			c.JSON(400, gin.H{"error": "Failed to bind base"})
// 			return
// 		}

// 		i, err := strconv.ParseInt(id, 10, 64)
// 		if err != nil {
// 			logger.Error().Err(err).Msg("Failed to parse ID")
// 			c.JSON(400, gin.H{"error": "Failed to parse ID"})
// 			return
// 		}

// 		updates := map[string]interface{}{}
// 		if input.Name != "" {
// 			updates["name"] = input.Name
// 		}
// 		if input.Description != "" {
// 			updates["description"] = input.Description
// 		}

// 		if len(updates) == 0 {
// 			logger.Error().Msg("No fields to update")
// 			c.JSON(400, gin.H{"error": "No fields to update"})
// 			return
// 		}

// 		if err := db.Model(&Base{}).Where("id = ?", i).Updates(updates).Error; err != nil {
// 			logger.Error().Err(err).Msg("Failed to update base")
// 			c.JSON(500, gin.H{"error": "Failed to update base"})
// 			return
// 		}
		
// 		var updateBase Base
// 		if err := db.First(&updateBase, i).Error; err != nil {
// 			logger.Error().Err(err).Msg("Failed to get updated base")
// 			c.JSON(500, gin.H{"error": "Failed to get updated base"})
// 			return
// 		}

// 		logger.Debug().Msg("Base updated")
// 		logger.Info().Msg("Base updated")
// 		c.JSON(205, gin.H{"data": updateBase})
// 	}
// }

// func DeleteBase(db *gorm.DB, logger zerolog.Logger) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Получаем ID базы из пути
// 		id := c.Param("id")
// 		if id == "" {
// 			logger.Error().Msg("Empty ID in URI")
// 			c.JSON(404, gin.H{"error": "Empty ID in URI"})
// 			return
// 		}

// 		logger.Debug().Msg("Deleting base...")

// 		i, err := strconv.ParseInt(id, 10, 64)
// 		if err != nil {
// 			logger.Error().Err(err).Msg("Failed to parse ID")
// 			c.JSON(400, gin.H{"error": "Failed to parse ID"})
// 			return
// 		}

// 		if err := db.Delete(&Base{}, i).Error; err != nil {
// 			logger.Error().Err(err).Msg("Failed to delete base")
// 			c.JSON(500, gin.H{"error": "Failed to delete base"})
// 			return
// 		}

// 		logger.Debug().Msg("Base deleted")
// 		logger.Info().Msg("Base deleted")
// 		c.JSON(204, gin.H{})
// 	}
// }